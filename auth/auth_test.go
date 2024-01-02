package auth

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAuth(t *testing.T) {
	// 创建一个测试用的请求
	reqBody := []byte("user_id=123&sessid=abc&uuid=xyz&api=a_1681997958&sign=0d00d67d476eea33723363f2ff2f0a87&sign_time=1641183445")
	req, err := http.NewRequest("POST", "/auth", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// 创建一个 ResponseRecorder 用于记录响应
	rr := httptest.NewRecorder()

	// 测试 Auth 函数
	Auth(req, rr, nil)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := ""
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
