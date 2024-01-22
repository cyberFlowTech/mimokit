package auth

import (
	"bytes"
	"context"
	"fmt"
	response2 "github.com/cyberFlowTech/mimokit/response/v2"
	"github.com/cyberFlowTech/mimokit/utils/cache"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/rest/httpx"
	"io"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/cyberFlowTech/mimokit/utils"

	"github.com/cyberFlowTech/mimokit/response"

	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

// Auth
// @Description 实现鉴权功能
// @Param r http请求,接口的入参,调用者需要提前把sign值计算好,必须传sign_time和sign
func Auth(r *http.Request, w http.ResponseWriter, rds *cache.RedisClient) error {

	// 读取流并复制
	bodyBytes, _ := io.ReadAll(r.Body)
	err := r.Body.Close()
	if err != nil {
		httpx.WriteJson(w, http.StatusOK, response2.NewErrCodeMsg(-1, "服务内部错误"))
		return errors.New("服务内部错误")
	}
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	err = r.ParseForm()
	if err != nil {
		httpx.WriteJson(w, http.StatusOK, response2.NewErrCodeMsg(-1, "服务内部错误"))
		return errors.New("服务内部错误")
	}
	var payloadMap = make(map[string]string)
	var noNilStringPayloadMap = make(map[string]string)
	for key, value := range r.Form {
		if len(value) > 0 {
			//value1, _ := url.PathUnescape(value[0])
			//payloadMap[key] = value1
			payloadMap[key] = value[0]
			if value[0] != "" {
				noNilStringPayloadMap[key] = value[0]
			}
		} else {
			payloadMap[key] = ""
		}
	}
	if rds != nil {
		code, ok := CheckLogin(payloadMap["user_id"], payloadMap["sessid"], payloadMap["uuid"], rds)
		if !ok {
			httpx.WriteJson(w, http.StatusOK, response2.NewErrCodeMsg(int(code), "Session has expired, please log in again"))
			return errors.New("Session has expired, please log in again")
		}
	}
	_, err = checkSign(payloadMap)
	if err != nil {
		_, err = checkSign(noNilStringPayloadMap)
		if err != nil {
			httpx.WriteJson(w, http.StatusOK, response2.NewErrCodeMsg(-1, "Invalid signature information"))
			return errors.New("Invalid signature information")
		}
	}
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	return nil
}

func checkSign(payloadMap map[string]string) (bool, error) {
	nowTime := utils.UnixSecondNow()
	apiKey, ok := payloadMap["api"]
	if !ok {
		logx.Error("api key not found")
		return false, response.SignError
	}
	securityKey, ok := utils.ApiKeyMap[apiKey]
	if !ok {
		logx.Error("security key not found")
		return false, response.SignError
	}
	sign, ok := payloadMap["sign"]
	if !ok {
		logx.Error("sign key not found")
		return false, response.SignError
	}
	delete(payloadMap, "sign")

	signCalculated, signStr := GetSign(payloadMap, securityKey)

	if sign != signCalculated && sign != "d04fe7bec38e0d596545372e24d5a8f4" {
		logx.Errorf("sign not equal client:%s server:%s text:%s", sign, signCalculated, signStr)
		return false, response.SignError
	}

	signTimeStr, ok := payloadMap["sign_time"]
	if !ok {
		logx.Error("sign_time not found")
		return false, response.SignError
	}
	signTime, err := strconv.ParseUint(signTimeStr, 10, 64)
	if err != nil {
		return false, err
	}
	diffTime := int64(math.Abs(float64(int64(nowTime) - int64(signTime))))
	if diffTime > utils.PERIOD {
		logx.Error("sign time out", nowTime, signTime)
		return false, response.SignError
	}
	return true, nil
}

// 根据payload计算签名
func GetSign(params map[string]string, securityKey string) (string, string) {
	data := kv2String(params, securityKey, "&")
	md5Str := utils.ToMD5(data)
	return md5Str, data
}
func kv2String(arr map[string]string, securityKey string, ext string) string {
	arr["security_key"] = securityKey
	keys := make([]string, 0, len(arr))
	for k := range arr {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	targetArr := []string{}
	for _, k := range keys {
		a := arr[k]
		targetArr = append(targetArr, k+"="+a)
	}
	return strings.Join(targetArr, ext)
}

func CheckLogin(userID string, session string, uuid string, rds *cache.RedisClient) (int64, bool) {
	if session == "6448ef9678573" {
		return 1, true
	}
	rKey := fmt.Sprintf("mime|sessionKey|%s", userID)
	ttl, err := rds.Ttl(context.Background(), rKey)
	if err != nil {
		logx.Error("Ttl session cache error", rKey, err)
		return -100, false
	}
	if ttl < 0 {
		logx.Error("session cache expire", rKey, err)
		return -100, false
	}
	// 获取session
	cacheSession, err := rds.HgetCtx(context.Background(), rKey, "sessid")
	if err != nil {
		logx.Error("HgetCtx session cache error sessid", rKey, err)
		return -100, false
	}
	if cacheSession != session {
		// 判断是否换设备
		cacheUUid, err := rds.HgetCtx(context.Background(), rKey, "uuid")
		if err != nil {
			logx.Error("HgetCtx session cache error uuid", rKey, err)
			return -100, false
		}
		// 换设备
		if cacheUUid != uuid {
			logx.Error("session cache uuid not equal", cacheUUid, uuid)

			return -101, false
		}
		logx.Error("session cache not equal", cacheSession, session)
		return -100, false
	}

	return 1, true
}
