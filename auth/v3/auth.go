package auth

import (
	"github.com/cyberFlowTech/mimokit/response/v3"
	"github.com/cyberFlowTech/mimokit/utils"
	"github.com/cyberFlowTech/rpclient/authen"
	"net/http"
)

// Auth V3
//
//	@Description: 校验签名和登陆态v3
//	@param r 请求体
//	@param w 响应体
//	@param client 认证客户端
//	@param checkLogin 是否校验登陆态
//	@return error
func AuthV3(r *http.Request, w http.ResponseWriter, client authen.AuthenClient, checkLogin bool) error {
	args, err := utils.GetPostFormData(r)
	if err != nil {
		return err
	}
	var res *authen.CheckAuthResp
	res, err = client.CheckAuth(r.Context(), &authen.CheckAuthReq{
		Args:       args,
		CheckLogin: checkLogin, // 检查登陆态
	})
	if err != nil {
		return err
	}
	if res.Pass == false {
		return response.NewErrMsg("Failed to authenticate")
	}
	return nil
}
