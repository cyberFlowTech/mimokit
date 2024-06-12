package utils

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// apikey映射关系
var (
	ApiKeyMap = map[string]string{
		"a_1648893980": "*&^%as#@$%#%",   //全平台
		"a_1648893915": "*&^%as#$&#@%",   //安卓机
		"i_1648893994": "*&^$as%%$&#@",   //IOS
		"w_1669703667": "!#Q&xxr!%!*@",   //web PC
		"p_1603348974": "*&^%mime@#^&!*", //后台apikey
		"p_1652180136": "!&@hmil*^!!!@",  //正式服udp
		"p_1652180080": "!&@@hmel&!!!@@", //开发服udp
		"p_1619681152": "*&^%$!@#^sa&!*", //IM udp 验证码
		"m_1652180201": "!&A&mmr!!!@@",   //消息服务
		"m_1665407443": "!@#$%^#^&$##",   //国际Exceptions(全平台)
		"p_1669894266": "!@!%%^#^@$#%@",  //error 404 udp
		"m_1670923199": "*@&%^^#^$$*%@",  //APP客服
		"a_1680591588": "*&^%aw#$@#@%",   //原生安卓机
		"i_1681907958": "*&^%as#$@T@%",   //原生ios机
		"a_1681997958": "*g!%asd#@Ts%",   // 环信切换 uniapp 安卓机
		"i_1681997960": "!&^%as&$ST@#",   // 环信切换 uniapp ios机
	}
)

func GetPlatformFromApi(api string) string {
	if _, ok := ApiKeyMap[api]; ok {
		strs := strings.Split(api, "_")
		//10 => 'iPhone',     //苹果
		//	20 => 'Android',    //安卓
		//	30 => 'PC Web'          //pc web端
		switch strs[0] {
		case "i":
			return "10"
		case "a":
			return "20"
		case "w":
			return "30"
		}
	} else {
		return "99"
	}
	return "99"
}

const (
	PERIOD         = 72000           // 签名超时时间,单位秒
	ClubEncryptKey = "!1@q#w2e$%^#@" // 部落相关id加密key
)

func GetStatefulSetIndex() (uint, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return 0, fmt.Errorf("failed to get hostname: %v", err)
	}
	fmt.Println("Hostname:", hostname)

	if len(hostname) < 1 {
		return 0, fmt.Errorf("hostname is empty")
	}

	// 检查主机名是否符合 StatefulSet 的命名规范
	parts := strings.Split(hostname, "-")
	if len(parts) < 2 {
		// 假设是本地或单节点环境
		fmt.Println("Running in local or single-node environment, returning index 0")
		return 0, nil
	}

	// 尝试将最后一部分解析为整数
	index, err := strconv.ParseUint(parts[len(parts)-1], 10, 64)
	if err != nil {
		// 如果解析失败，假设是本地环境
		fmt.Println("Failed to parse index, assuming local environment:", err)
		return 0, nil
	}

	return uint(index), nil
}
