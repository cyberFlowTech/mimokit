package response

// 定义全局公共错误

var (
	SignError            = New(1000004, "sign参数校验错误")
	ServiceInternalError = New(1000009, "服务内部错误")

	ExpireError = New(5000003, "token已过期")
)
