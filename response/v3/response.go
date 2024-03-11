package response

import (
	"fmt"
	"strings"
)

type CodeError struct {
	IRet int         `json:"iRet"`
	SMsg string      `json:"sMsg"`
	Data interface{} `json:"data"`
}

// GetErrCode 返回给前端的错误码
func (e *CodeError) GetErrCode() int {
	return e.IRet
}

// GetErrMsg 返回给前端显示端错误信息
func (e *CodeError) GetErrMsg() string {
	return e.SMsg
}

func (e *CodeError) Error() string {
	return fmt.Sprintf("ErrCode:%d，ErrMsg:%s", e.IRet, e.SMsg)
}

type CodeSuccess struct {
	IRet int         `json:"iRet"`
	SMsg string      `json:"sMsg"`
	Data interface{} `json:"data"`
}

func Success(data interface{}) *CodeSuccess {
	return &CodeSuccess{OK, "OK", data}
}

func NewErrCodeMsg(errCode int, errMsg ...string) *CodeError {
	msg := ""
	if len(errMsg) == 1 {
		msg = errMsg[0]
	} else if len(errMsg) > 1 {
		msg = strings.Join(errMsg, ",")
	}
	return &CodeError{IRet: errCode, SMsg: msg}
}

func NewErrMsg(errMsg string) *CodeError {
	return &CodeError{IRet: -1, SMsg: errMsg}
}
