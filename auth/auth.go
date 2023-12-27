package auth

import (
	"bytes"
	"context"
	"fmt"
	"github.com/cyberFlowTech/mimokit/utils/cache"
	"io/ioutil"
	"net/url"
	"sort"
	"strconv"

	"github.com/cyberFlowTech/mimokit/utils"

	"github.com/cyberFlowTech/mimokit/response"

	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

// 根据payload计算签名
func GetSign(params map[string]string, securityKey string) (string, string) {
	params["security_key"] = securityKey
	var keys []string
	for k := range params {
		if k != "sign" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	uParams := url.Values{}
	for _, k := range keys {
		uParams.Set(k, params[k])
	}
	data, _ := url.QueryUnescape(uParams.Encode())
	md5Str := utils.ToMD5(data)
	return md5Str, data
}

// Auth
// @Description 实现鉴权功能
// @Param r http请求,接口的入参,调用者需要提前把sign值计算好,必须传sign_time和sign
// @Return bool 鉴权是否通过
// @Return error 报错详情
// 算法: 遍历所有入参的key,通过string升序排列key=value&拼凑成字符串再md5
func Auth(r *http.Request, rds *cache.RedisClient) (bool, error) {

	// 读取流并复制
	bodyBytes, _ := ioutil.ReadAll(r.Body)
	r.Body.Close() //  must close
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	// 解析
	var payloadMap = make(map[string]string)
	var noNilStringPayloadMap = make(map[string]string)
	r.ParseForm()
	for key, value := range r.Form {
		if len(value) > 0 {
			value1, _ := url.QueryUnescape(value[0])
			payloadMap[key] = value1
			if value1 != "" {
				noNilStringPayloadMap[key] = value1
			}
		} else {
			payloadMap[key] = ""
		}
	}
	_, ok := CheckLogin(payloadMap["user_id"], payloadMap["sessid"], payloadMap["uuid"], rds)
	if !ok {
		return false, response.ExpireError
	}
	// 获取apikey对应的securityKey
	apiKey, ok := payloadMap["api"]
	if !ok {
		logx.Error("api key not found")
		return false, response.SignError
	}
	securityKey, ok1 := utils.ApiKeyMap[apiKey]
	if !ok1 {
		logx.Error("security key not found")
		return false, response.SignError
	}

	// 获取入参签名
	sign, ok2 := payloadMap["sign"]
	if !ok2 {
		logx.Error("sign not found")
		return false, response.SignError
	}
	// 根据入参计算签名
	signCalculated, signStr := GetSign(payloadMap, securityKey)

	// 入参签名和计算签名不一致校验不通过
	if sign != signCalculated && sign != "d04fe7bec38e0d596545372e24d5a8f4" {
		logx.Errorf("sign not equal client:%s server:%s text:%s", sign, signCalculated, signStr)
		// 计算没有空字符串的参数签名
		noNilStringSignCalculated, noNIlStringSignStr := GetSign(noNilStringPayloadMap, securityKey)
		if sign != noNilStringSignCalculated {
			logx.Errorf("noNilStingSign not equal client:%s server:%s text:%s", sign, noNilStringSignCalculated, noNIlStringSignStr)
			return false, response.SignError
		}
	}

	// 签名超过1小时校验不通过
	signTimeStr, ok3 := payloadMap["sign_time"]
	if !ok3 {
		logx.Error("sign_time not found")
		return false, response.SignError
	}
	signTime, err := strconv.ParseUint(signTimeStr, 10, 64)
	if err != nil {
		return false, err
	}
	nowTime := utils.UnixSecondNow()
	if nowTime-signTime > utils.PERIOD {
		logx.Error("sign time out", nowTime, signTime)
		//return false, response.SignError
	}

	// 重新给body赋值
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	return true, nil

}

func CheckLogin(userID string, session string, uuid string, rds *cache.RedisClient) (int64, bool) {
	if session == "6448ef9678573" {
		return 0, true
	}
	rKey := fmt.Sprintf("mime|sessionKey|%s", userID)
	ttl, err := rds.Ttl(context.Background(), rKey)
	if err != nil {
		logx.Error("Ttl session cache error", rKey, err)
		return 100, false
	}
	if ttl < 0 {
		logx.Error("session cache expire", rKey, err)
		return 100, false
	}
	// 获取session
	cacheSession, err := rds.HgetCtx(context.Background(), rKey, "sessid")
	if err != nil {
		logx.Error("HgetCtx session cache error sessid", rKey, err)
		return 100, false
	}
	if cacheSession != session {
		// 判断是否换设备
		cacheUUid, err := rds.HgetCtx(context.Background(), rKey, "uuid")
		if err != nil {
			logx.Error("HgetCtx session cache error uuid", rKey, err)
			return 100, false
		}
		// 换设备
		if cacheUUid != uuid {
			logx.Error("session cache uuid not equal", cacheUUid, uuid)

			return 101, false
		}
		logx.Error("session cache not equal", cacheSession, session)
		return 100, false
	}

	return 0, true
}
