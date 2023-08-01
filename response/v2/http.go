package response

import (
	"net/http"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"google.golang.org/grpc/status"
)

/*HTTP请求返回处理*/

type HTTPResponse struct {
	//内部错误码，非自定义错误
	ServerCommonErrorCode int
	//登录失效错误码
	TokenExpireErrorCode int
	//所有的错误信息，key为错误码，value为错误描述
	message map[int]string
}

func NewHTTPResponse(serverCommonErrorCode, tokenExpireErrorCode int, message map[int]string) *HTTPResponse {
	return &HTTPResponse{ServerCommonErrorCode: serverCommonErrorCode, TokenExpireErrorCode: tokenExpireErrorCode, message: message}
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
		errmsg := "服务器开小差啦，稍后再来试一试"

		causeErr := errors.Cause(err)           // err类型
		if e, ok := causeErr.(*CodeError); ok { //自定义错误类型
			//自定义CodeError
			errcode = e.GetErrCode()
			errmsg = e.GetErrMsg()
		} else {
			if gstatus, ok := status.FromError(causeErr); ok { // grpc err错误
				grpcCode := uint32(gstatus.Code())
				if h.IsCodeErr(int(grpcCode)) { //区分自定义错误跟系统底层、db等错误，底层、db错误不能返回给前端
					errcode = int(grpcCode)
					errmsg = gstatus.Message()
				}
			}
		}
		if errcode == h.TokenExpireErrorCode {
			errcode = TokenExpiredErrorCode
		} else {
			errcode = UniformErrorCode
		}

		logx.WithContext(r.Context()).Errorf("【API-ERR】 : %+v ", err)
		httpx.WriteJson(w, http.StatusOK, NewErrCodeMsg(errcode, errmsg))
	}
}

// MapErrMsg 判断是否自定义错误
func (h *HTTPResponse) MapErrMsg(errCode int) string {
	if msg, ok := h.message[errCode]; ok {
		return msg
	} else {
		return "服务器开小差啦,稍后再来试一试"
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
