package response

// 通用错误码
const UniformErrorCode = -1        // 前端统一提示错误码
const OK = 1                       // 成功返回
const RecordNotFound = 2           // 记录不存在
const TokenExpiredErrorCode = -100 // Session过期
const UuidChanged = -101           // 更换了登录设备
