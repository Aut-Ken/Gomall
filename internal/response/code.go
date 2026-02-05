package response

// ============================================
// 通用状态码 (0-999)
// ============================================
const (
	CodeSuccess         = 0     // 成功
	CodeBadRequest      = 400   // 参数错误
	CodeUnauthorized    = 401   // 未登录
	CodeForbidden       = 403   // 无权限
	CodeNotFound        = 404   // 不存在
	CodeMethodNotAllow  = 405   // 方法不允许
	CodeConflict        = 409   // 冲突
	CodeTooManyRequests = 429   // 请求过于频繁
	CodeServerError     = 500   // 系统错误
	CodeServiceUnavailable = 503 // 服务不可用
)

// ============================================
// 用户模块错误码 (10001-10099)
// ============================================
const (
	// 用户相关 10001-10010
	CodeUserNotFound       = 10001 // 用户不存在
	CodeUserAlreadyExist   = 10002 // 用户已存在
	CodeUserDisabled      = 10003 // 用户已被禁用
	CodeUserPasswordError = 10004 // 密码错误
	CodeUserTokenExpired  = 10005 // Token已过期
	CodeUserTokenInvalid  = 10006 // Token无效
	CodeUserNoPermission  = 10007 // 无权限操作
	CodeUserLoginRequired = 10008 // 需要登录
	CodeUserParamError    = 10009 // 用户参数错误
	CodeUserNotAdmin      = 10010 // 非管理员用户
)

// ============================================
// 商品模块错误码 (20001-20099)
// ============================================
const (
	// 商品相关 20001-20010
	CodeProductNotFound     = 20001 // 商品不存在
	CodeProductOffShelf     = 20002 // 商品已下架
	CodeProductStockEmpty   = 20003 // 商品库存不足
	CodeProductStockNotEnough = 20004 // 商品库存不足
	CodeProductCategoryError = 20005 // 商品分类错误
	CodeProductParamError   = 20006 // 商品参数错误
	CodeProductCreateFailed = 20007 // 商品创建失败
	CodeProductUpdateFailed = 20008 // 商品更新失败
	CodeProductDeleteFailed = 20009 // 商品删除失败
	CodeProductStatusError = 20010 // 商品状态错误
)

// ============================================
// 订单模块错误码 (30001-30099)
// ============================================
const (
	// 订单相关 30001-30010
	CodeOrderNotFound       = 30001 // 订单不存在
	CodeOrderCreateFailed   = 30002 // 订单创建失败
	CodeOrderStatusError    = 30003 // 订单状态错误
	CodeOrderPayFailed      = 30004 // 订单支付失败
	CodeOrderCancelFailed   = 30005 // 订单取消失败
	CodeOrderExpired        = 30006 // 订单已过期
	CodeOrderPaid           = 30007 // 订单已支付
	CodeOrderNotPaid        = 30008 // 订单未支付
	CodeOrderParamError     = 30009 // 订单参数错误
	CodeOrderNoPaymentMethod = 30010 // 无支付方式
)

// ============================================
// 支付模块错误码 (40001-40099)
// ============================================
const (
	// 支付相关 40001-40010
	CodePayCreateFailed    = 40001 // 支付创建失败
	CodePayAmountError     = 40002 // 支付金额错误
	CodePayTimeout         = 40003 // 支付超时
	CodePayCancelFailed    = 40004 // 支付取消失败
	CodePayRefundFailed    = 40005 // 退款失败
	CodePayQueryFailed     = 40006 // 支付查询失败
	CodePaySignError       = 40007 // 签名错误
	CodePayChannelError    = 40008 // 支付渠道错误
	CodePayParamError      = 40009 // 支付参数错误
	CodePayNotFound        = 40010 // 支付记录不存在
)

// ============================================
// 购物车模块错误码 (50001-50099)
// ============================================
const (
	// 购物车相关 50001-50010
	CodeCartNotFound       = 50001 // 购物车为空
	CodeCartItemNotFound   = 50002 // 购物车商品不存在
	CodeCartAddFailed      = 50003 // 添加购物车失败
	CodeCartUpdateFailed   = 50004 // 更新购物车失败
	CodeCartDeleteFailed   = 50005 // 删除购物车失败
	CodeCartClearFailed    = 50006 // 清空购物车失败
	CodeCartStockNotEnough = 50007 // 库存不足
	CodeCartParamError     = 50008 // 购物车参数错误
	CodeCartProductOffShelf = 50009 // 商品已下架
	CodeCartProductNotExist = 50010 // 商品不存在
)

// ============================================
// 秒杀模块错误码 (60001-60099)
// ============================================
const (
	// 秒杀相关 60001-60010
	CodeSeckillNotStart    = 60001 // 秒杀未开始
	CodeSeckillEnded       = 60002 // 秒杀已结束
	CodeSeckillStockEmpty  = 60003 // 秒杀库存不足
	CodeSeckillRepeatBuy   = 60004 // 重复购买
	CodeSeckillCreateOrderFail = 60005 // 秒杀订单创建失败
	CodeSeckillParamError  = 60006 // 秒杀参数错误
	CodeSeckillLimitBuy    = 60007 // 超过购买限制
	CodeSeckillNotFound    = 60008 // 秒杀活动不存在
	CodeSeckillEndedOrNotStart = 60009 // 秒杀活动未开始或已结束
	CodeSeckillHighRequest = 60010 // 请求过于频繁
)

// ============================================
// 文件上传模块错误码 (70001-70099)
// ============================================
const (
	// 文件上传相关 70001-70010
	CodeUploadFileEmpty    = 70001 // 上传文件为空
	CodeUploadFileTooLarge = 70002 // 文件过大
	CodeUploadFileTypeError = 70003 // 文件类型不支持
	CodeUploadSaveFailed   = 70004 // 文件保存失败
	CodeUploadNotFound     = 70005 // 文件不存在
	CodeUploadDeleteFailed = 70006 // 文件删除失败
	CodeUploadTooManyFiles = 70007 // 上传文件数量过多
	CodeUploadParamError   = 70008 // 上传参数错误
	CodeUploadServerError  = 70009 // 文件服务器错误
	CodeUploadNotAllowed   = 70010 // 无上传权限
)

// ============================================
// 错误码映射（错误码 -> 错误消息）
// ============================================
var codeMsgMap = map[int]string{
	// 通用
	CodeSuccess:           "成功",
	CodeBadRequest:        "参数错误",
	CodeUnauthorized:     "未登录或登录已过期",
	CodeForbidden:        "没有权限",
	CodeNotFound:          "资源不存在",
	CodeMethodNotAllow:    "方法不允许",
	CodeConflict:          "请求冲突",
	CodeTooManyRequests:   "请求过于频繁",
	CodeServerError:       "系统错误",
	CodeServiceUnavailable: "服务不可用",

	// 用户
	CodeUserNotFound:       "用户不存在",
	CodeUserAlreadyExist:   "用户已存在",
	CodeUserDisabled:      "用户已被禁用",
	CodeUserPasswordError: "密码错误",
	CodeUserTokenExpired:  "Token已过期",
	CodeUserTokenInvalid:  "Token无效",
	CodeUserNoPermission:  "无权限操作",
	CodeUserLoginRequired: "需要登录",
	CodeUserParamError:    "用户参数错误",
	CodeUserNotAdmin:      "非管理员用户",

	// 商品
	CodeProductNotFound:     "商品不存在",
	CodeProductOffShelf:     "商品已下架",
	CodeProductStockEmpty:   "商品库存不足",
	CodeProductStockNotEnough: "商品库存不足",
	CodeProductCategoryError: "商品分类错误",
	CodeProductParamError:   "商品参数错误",
	CodeProductCreateFailed: "商品创建失败",
	CodeProductUpdateFailed: "商品更新失败",
	CodeProductDeleteFailed: "商品删除失败",
	CodeProductStatusError: "商品状态错误",

	// 订单
	CodeOrderNotFound:       "订单不存在",
	CodeOrderCreateFailed:   "订单创建失败",
	CodeOrderStatusError:    "订单状态错误",
	CodeOrderPayFailed:      "订单支付失败",
	CodeOrderCancelFailed:   "订单取消失败",
	CodeOrderExpired:        "订单已过期",
	CodeOrderPaid:           "订单已支付",
	CodeOrderNotPaid:        "订单未支付",
	CodeOrderParamError:     "订单参数错误",
	CodeOrderNoPaymentMethod: "无支付方式",

	// 支付
	CodePayCreateFailed:    "支付创建失败",
	CodePayAmountError:     "支付金额错误",
	CodePayTimeout:         "支付超时",
	CodePayCancelFailed:    "支付取消失败",
	CodePayRefundFailed:    "退款失败",
	CodePayQueryFailed:     "支付查询失败",
	CodePaySignError:       "签名错误",
	CodePayChannelError:    "支付渠道错误",
	CodePayParamError:      "支付参数错误",
	CodePayNotFound:        "支付记录不存在",

	// 购物车
	CodeCartNotFound:       "购物车为空",
	CodeCartItemNotFound:   "购物车商品不存在",
	CodeCartAddFailed:      "添加购物车失败",
	CodeCartUpdateFailed:   "更新购物车失败",
	CodeCartDeleteFailed:   "删除购物车失败",
	CodeCartClearFailed:    "清空购物车失败",
	CodeCartStockNotEnough: "库存不足",
	CodeCartParamError:     "购物车参数错误",
	CodeCartProductOffShelf: "商品已下架",
	CodeCartProductNotExist: "商品不存在",

	// 秒杀
	CodeSeckillNotStart:    "秒杀未开始",
	CodeSeckillEnded:       "秒杀已结束",
	CodeSeckillStockEmpty:  "秒杀库存不足",
	CodeSeckillRepeatBuy:   "重复购买",
	CodeSeckillCreateOrderFail: "秒杀订单创建失败",
	CodeSeckillParamError:  "秒杀参数错误",
	CodeSeckillLimitBuy:    "超过购买限制",
	CodeSeckillNotFound:    "秒杀活动不存在",
	CodeSeckillEndedOrNotStart: "秒杀活动未开始或已结束",
	CodeSeckillHighRequest: "请求过于频繁",

	// 文件上传
	CodeUploadFileEmpty:    "上传文件为空",
	CodeUploadFileTooLarge: "文件过大",
	CodeUploadFileTypeError: "文件类型不支持",
	CodeUploadSaveFailed:   "文件保存失败",
	CodeUploadNotFound:     "文件不存在",
	CodeUploadDeleteFailed: "文件删除失败",
	CodeUploadTooManyFiles: "上传文件数量过多",
	CodeUploadParamError:   "上传参数错误",
	CodeUploadServerError:  "文件服务器错误",
	CodeUploadNotAllowed:   "无上传权限",
}

// GetCodeMsg 根据错误码获取错误消息
func GetCodeMsg(code int) string {
	if msg, ok := codeMsgMap[code]; ok {
		return msg
	}
	return "未知错误"
}
