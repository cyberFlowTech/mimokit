package auth

import (
	"github.com/cyberFlowTech/mimokit/response/v3"
	"github.com/cyberFlowTech/mimokit/utils"
	"github.com/cyberFlowTech/rpclient/authen"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

// Auth
//
//	@Description: 校验签名和登陆态
//	@param r
//	@param w
//	@param client
//	@param checkLogin 是否校验登陆态
//	@return error
func AuthV2(r *http.Request, w http.ResponseWriter, client authen.AuthenClient, checkLogin bool) error {
	args, err := utils.GetPostFormData(r)
	if err != nil {
		httpx.WriteJson(w, http.StatusOK, response.NewErrCodeMsg(-1, "auth server error"))
		return err
	}
	var res *authen.CheckAuthResp
	res, err = client.CheckAuth(r.Context(), &authen.CheckAuthReq{
		Args:       args,
		CheckLogin: checkLogin, // 检查登陆态
	})
	if err != nil {
		if b, c := response.IsCodeError(err.Error()); b {
			httpx.WriteJson(w, http.StatusOK, response.NewErrCodeMsg(c.GetErrCode(), c.GetErrMsg()))
		}
		return err
	}
	if res.Pass == false {
		httpx.WriteJson(w, http.StatusOK, response.NewErrCodeMsg(-1, "auth fail"))
		return errors.New("auth fail")
	}
	return nil
}
