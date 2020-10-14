package util

const (
	SuccessCode    = 0
	ErrorLackCode  = 1
	ErrorSqlCode   = 2
	ErrorRidesCode = 3

	//参数签名秘钥
	DesKey = "r5k1*8a$@8dc!dytkcs2dqz!"

	//redis key
	RedisKeyRegisteredCode       = "user:registered:code:"        //注册验证码
	RedisKeyRegisteredCodeNumber = "user:registered:code:number:" //注册验证码次数
	RedisKeyForgotCode           = "user:forgot:code:"            //忘记密码验证码
	RedisKeyForgotCodeNumber     = "user:forgot:code:number:"     //忘记密码验证码次数
	RedisKeyBindingCode          = "user:binding:code:"           //绑定手机验证码
	RedisKeyBindingCodeNumber    = "user:binding:code:number:"    //绑定手机验证码次数

	RedisKeyToken = "user:login:token:" //token 缓存

	RedisKeyLoginServer    = "user:login:server:list:" //登录服务器
	RedisKeyHallServerList = "user:hall:server:list:"  //大厅服务器列表

	RedisKeyTeam = "user:team:" //team缓存

	RedisKeyTeamDire            = ":dire"             //直属查询
	RedisKeyTeamAchievement     = ":achievement"      //业绩查询
	RedisKeyTeamAchievementForm = ":achievement:form" //业绩来源

	//time
	FormatTime      = "15:04:05"            //时间格式
	FormatDate      = "2006-01-02"          //日期格式
	FormatDateTime  = "2006-01-02 15:04:05" //完整时间格式
	FormatDateTime2 = "2006-01-02 15:04"    //完整时间格式

)
