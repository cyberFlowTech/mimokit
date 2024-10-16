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
	Trans bool // 是否进行多语言翻译
}

/*HTTP请求返回处理*/
type HTTPResponse struct {
	Config
}

func NewHTTPResponse(config Config) *HTTPResponse {
	return &HTTPResponse{Config: config}
}

// JSON http请求返回JSON数据结果，同时记录日志
// 当error是自定义error，直接返回给前端，其它错误不直接返回
// 调用时可以用errors.Wrapf()，尽可能包裹详细错误信息，Logic层错误直接返回无须再打印日志
func (h *HTTPResponse) JSON(r *http.Request, w http.ResponseWriter, resp interface{}, err error) {
	if err == nil {
		//成功返回
		httpx.WriteJson(w, http.StatusOK, Success(resp))
	} else {
		//错误返回
		errorCode := UniformErrorCode
		errMsg := "The server has encountered an issue. Please try again later."
		//文案提示优先级：多语言转换->自定义消息->兜底文案errMsg
		causeErr := errors.Cause(err)           // err类型
		if e, ok := causeErr.(*CodeError); ok { //自定义错误类型
			//自定义CodeError
			errorCode = e.GetErrCode()
			// 自定义消息
			if e.GetErrMsg() != "" {
				errMsg = e.GetErrMsg()
			}

		} else {
			if gstatus, ok := status.FromError(causeErr); ok { // grpc err错误
				// 判断是否为自定义错误。根据CodeError Error格式进行判断
				str := gstatus.String()
				reCode := regexp.MustCompile(`ErrCode:[-]?(\d+)`)
				if reCode.MatchString(str) {
					s := reCode.FindString(str)
					splits := strings.Split(s, ":")
					if len(splits) == 2 {
						errorCode, _ = strconv.Atoi(splits[1])
						// 自定义消息优先级更高
						reMsg := regexp.MustCompile(`ErrMsg:(.*)`)
						msg := reMsg.FindString(str)
						msg = strings.Replace(msg, "ErrMsg:", "", 1)
						if msg != "" {
							errMsg = msg
						}

					}
				}
			}
		}
		// 多语言转换优先级最高
		if h.Config.Trans == true && r.FormValue("lan") != "" && errorCode != -1 && errorCode != 2 {
			if msg := lan.Trans(r.FormValue("lan"), strconv.Itoa(errorCode)); msg != "" {
				if !strings.Contains(msg, "The current network is congested, please wait") {
					//能转译成功则赋值
					errMsg = msg
				}
			}
		}
		if errorCode == -100 {
			errorCode = TokenExpiredErrorCode
		} else if errorCode == 2 {
			errorCode = RecordNotFound //记录不存在,终端做跳转处理
		} else if errorCode == -101 {
			errorCode = UuidChanged // uuid更换了
		} else {
			errorCode = UniformErrorCode
		}

		logx.WithContext(r.Context()).Errorf("【ERR】uri:%v|uid:%s|v:%s|p:%s|err:%s",
			r.RequestURI,
			r.FormValue("user_id"),
			r.FormValue("version"),
			r.FormValue("api"),
			err.Error())
		httpx.WriteJson(w, http.StatusOK, NewErrCodeMsg(errorCode, errMsg))
	}
}
