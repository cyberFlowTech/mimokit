# HTTP返回请求处理

* 我们在平时开发时候，程序在出错时，希望可以通过错误日志能快速定位问题（那么传递进来的参数、包括堆栈信息肯定就要都要打印到日志），
  但同时又想返回给前端用户比较友善、能看得懂的错误提示， 除非在返回前端错误提示的地方同时在记录log，这样的话日志满天飞，代码不简洁，
  日志到时候也会很难看。

* V2版本支持统一的地方记录日志，同时在业务代码中只需要一个return err 就能将返回给前端的错误提示信息、日志记录相信信息分开提示跟记录

* 目前返回给前端的code分为-1、-100、1其中：-1表示前端可直接toast的错误；-100表示登录过期；1表示成功

* 自定义的全局错误码仅内部使用

* 调用实例：

~~~
//import包

var httpResponse *response.HTTPResponse

func init() {
	httpResponse = response.NewHTTPResponse(ServerCommonError, TokenExpireError, message)
}

func JSON(r *http.Request, w http.ResponseWriter, resp interface{}, err error) {
	httpResponse.JSON(r, w, resp, err)
}

func NewErrCode(errCode int) *response.CodeError {
	return httpResponse.NewErrCode(errCode)
}
~~~