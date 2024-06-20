package response

import (
	"fmt"
	"regexp"
	"strconv"
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

func IsCodeError(errMsg string) (bool, CodeError) {
	var codeError CodeError
	reCode := regexp.MustCompile(`ErrCode:[-]?(\d+)`)
	if reCode.MatchString(errMsg) {
		s := reCode.FindString(errMsg)
		splits := strings.Split(s, ":")
		if len(splits) == 2 {
			codeError.IRet, _ = strconv.Atoi(splits[1])
			// 自定义消息优先级更高
			reMsg := regexp.MustCompile(`ErrMsg:(.*)`)
			msg := reMsg.FindString(errMsg)
			msg = strings.Replace(msg, "ErrMsg:", "", 1)
			if msg != "" {
				codeError.SMsg = msg
			}
		}
		return true, codeError
	}
	return false, codeError
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
