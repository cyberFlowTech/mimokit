package response

import (
	"github.com/cyberFlowTech/mimokit/lan"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"google.golang.org/grpc/status"
)

type Config struct {
	Trans                 bool // 是否进行多语言翻译
	ServerCommonErrorCode int  //内部错误码，非自定义错误
	TokenExpireErrorCode  int  // 	//登录失效错误码
}

/*HTTP请求返回处理*/
type HTTPResponse struct {
	Config
	//所有的错误信息，key为错误码，value为错误描述
	message map[int]string
}

func NewHTTPResponse(config Config, message map[int]string) *HTTPResponse {
	return &HTTPResponse{Config: config, message: message}
}

// JSON http请求返回JSON数据结果，同时记录日志
// 当error是自定义error，直接返回给前端，其它错误不直接返回
// 调用时可以用errors.Wrapf()，尽可能包裹详细错误信息，Logic层错误直接返回无须再打印日志
func (h *HTTPResponse) JSON(r *http.Request, w http.ResponseWriter, resp interface{}, err error) {
	if err == nil {
		//成功返回
		r := Success(resp)
		httpx.WriteJson(w, http.StatusOK, r)
	} else {
		//错误返回
		errcode := h.ServerCommonErrorCode
		errmsg := "The server has something wrong, please try again later."

		causeErr := errors.Cause(err)           // err类型
		if e, ok := causeErr.(*CodeError); ok { //自定义错误类型
			//自定义CodeError
			errcode = e.GetErrCode()
			errmsg = e.GetErrMsg()
		} else {
			if gstatus, ok := status.FromError(causeErr); ok { // grpc err错误
				// 判断是否为自定义错误。根据CodeError Error格式进行判断
				str := gstatus.String()
				reCode := regexp.MustCompile(`ErrCode:(\d+)`)
				if reCode.MatchString(str) {
					s := reCode.FindString(str)
					splits := strings.Split(s, ":")
					if len(splits) == 2 {
						//fmt.Println(splits[1]) // errcode
						reMsg := regexp.MustCompile(`ErrMsg:(.*)`)
						msg := reMsg.FindString(str)
						if msg != "" {
							errmsg = strings.Replace(msg, "ErrMsg:", "", 1)
						}
					}
				}
			}
		}
		// 多语言转换
		if h.Config.Trans == true && r.FormValue("lan") != "" {
			if msg := lan.Trans(r.FormValue("lan"), strconv.Itoa(errcode)); msg != "" {
				errmsg = msg
			}
		}
		if errcode == h.TokenExpireErrorCode {
			errcode = TokenExpiredErrorCode
		} else {
			errcode = UniformErrorCode
		}

		logx.WithContext(r.Context()).Errorf("【API-ERR】Uri:%v user_id:%v err: %+v ", r.RequestURI, r.FormValue("user_id"), err)
		httpx.WriteJson(w, http.StatusOK, NewErrCodeMsg(errcode, errmsg))
	}
}

// MapErrMsg 判断是否自定义错误
func (h *HTTPResponse) MapErrMsg(errCode int) string {
	if msg, ok := h.message[errCode]; ok {
		return msg
	} else {
		return "The current network is congested, please wait"
	}
}

// IsCodeErr 判断是否自定义错误
func (h *HTTPResponse) IsCodeErr(errCode int) bool {
	if _, ok := h.message[errCode]; ok {
		return true
	} else {
		return false
	}
}

func (h *HTTPResponse) NewErrCode(errCode int) *CodeError {
	return &CodeError{IRet: errCode, SMsg: h.MapErrMsg(errCode)}
}

func (h *HTTPResponse) NewErrMsg(errMsg string) *CodeError {
	return &CodeError{IRet: h.ServerCommonErrorCode, SMsg: errMsg}
}