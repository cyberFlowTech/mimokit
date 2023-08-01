package response

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// JSON
// @Description 统一项目返回结构
// 用法: 在handler中调用业务logic后封装resp返回
func JSON(w http.ResponseWriter, resp interface{}, err error) {
	if err != nil {
		httpx.OkJson(w, err)
	} else {
		var body BizError
		body.SMsg = "OK"
		body.Data = resp
		body.IRet = 1
		httpx.OkJson(w, body)
	}

}

type BizError struct {
	IRet int         `json:"iRet"`
	SMsg string      `json:"sMsg"`
	Data interface{} `json:"data"`
}

func New(code int, msg string) *BizError {
	return &BizError{IRet: code, SMsg: msg}
}

func (e *BizError) Error() string {
	return e.SMsg
}

func WrapError(bizError *BizError, err error) error {
	if err == nil {
		return nil
	}
	return &BizError{
		IRet: bizError.IRet,
		SMsg: bizError.SMsg + ":" + err.Error(),
	}
}

type InternalError struct {
	IRet int         `json:"iRet"`
	SMsg string      `json:"sMsg"`
	Data interface{} `json:"data"`
}

func InternalErrorRes() *InternalError {
	return &InternalError{
		IRet: ServiceInternalError.IRet,
		SMsg: ServiceInternalError.SMsg,
		Data: ServiceInternalError.Data,
	}
}
