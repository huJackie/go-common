package code

var (
	//全局通用
	ServeErr      = add(500, "服务器异常错误")
	OK            = add(10000, "OK")
	ParamsInvalid = add(10001, "参数错误")
	Retry         = add(10002, "请重试")
	RequestErr    = add(10003, "请求错误")
	NotFound      = add(10004, "数据不存在")
	ApiNotFound   = add(10005, "接口不存在")
	FileFormat    = add(10006, "文件格式错误")
	Denied        = add(10007, "无权限")
	JwtExpire     = add(10008, "重新登录")
	JwtInvalid    = add(10009, "token无效")
	JwtEmpty      = add(10010, "未登录")
	AccountErr    = add(10011, "账号或密码错误")
	AccountForbid = add(10012, "账号状态异常")
	DataExist     = add(10013, "数据已存在")
	AppNotExist   = add(10014, "应用不存在")
)
