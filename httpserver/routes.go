package httpserver

import (
	"github.com/gin-gonic/gin"
	"xj_web_server/httpserver/activity"
	"xj_web_server/httpserver/agent"
	"xj_web_server/httpserver/exchange"
	"xj_web_server/httpserver/game"
	"xj_web_server/httpserver/handle"
	"xj_web_server/httpserver/index"
	"xj_web_server/httpserver/news"
	"xj_web_server/httpserver/rank"
	"xj_web_server/httpserver/servermiddleware"
	"xj_web_server/httpserver/sign"
	"xj_web_server/httpserver/user"
	"xj_web_server/httpserver/wss"
	"xj_web_server/module/apigw"
)

func initRoutes(router *gin.Engine) {
	router.GET("/version", handle.Version)
	router.GET("/ping", handle.Ping)
	router.GET("/v1/check/time", handle.CheckTime)

	//websocket,
	router.GET("/v1/wss", wss.Wss)
	router.GET("/v1/client", wss.Client)

	router.POST("/v1/game/version/set", game.SetGameVersion)

	v1 := router.Group("/v1", servermiddleware.Base())
	{

		//获取登录服务器地址
		v1.POST("/host", index.GetHost)
		//获取大厅服务器地址
		v1.POST("/hall/host", index.GetHallHost)

		//注册
		v1.POST("/registered", user.Registered)
		//获取验证码
		v1.POST("/getcode", user.GetCode)
		//忘记密码
		v1.POST("/forgot/pwd", user.ForgotPwd)

		//获取用户信息
		userPath := v1.Group("/user", servermiddleware.BaseAuth())
		{
			userPath.POST("/msg", user.GetUserMsg)
			//战绩
			userPath.POST("/record/list", user.RecordList)
			//修改用户信息
			userPath.POST("/update/information", user.UpdateInformation)
			// 救济金规则列表
			userPath.POST("/benefits", user.BaseEnsures)
			//领取救济金
			userPath.POST("/benefits/receive", user.TakeBaseEnsure)
			// 绑定银行卡
			userPath.POST("/binding/bankcard", user.BindingBankcard)
			//绑定手机号码
			userPath.POST("/binding/mobile", user.BindingMobile)
		}

		//代理
		agentPath := v1.Group("/agent", servermiddleware.BaseAuth())
		{
			//我的推广
			agentPath.POST("/team/mypromote", agent.MyPromote)
			//直属
			agentPath.POST("/team/dire", agent.TeamDire)
			//返佣详情
			agentPath.POST("/team/level", agent.GetAgentRotyaltyLevel)
			//配置
			agentPath.POST("/team/config", agent.GetAgentRotyaltyConfig)
			//业绩
			agentPath.POST("/team/achievement", agent.TeamAchievement)
			//业绩来源
			agentPath.POST("/team/achievement/form", agent.TeamAchievementForm)
			//agentPath.POST("/team/promote", agent.TeamPromote)
			//领取佣金
			agentPath.POST("/team/agent/take", agent.TeamTakeAgentRoyalty)
			//佣金记录
			agentPath.POST("/team/agent/take/record", agent.TeamTakeAgentRecord)

			//一级代理
			agentPath.POST("/team/one", agent.One)

			//agentPath.POST("/team/record", agent.TeamList)
			//agentPath.POST("/daily/knots", agent.DailyKnots)
			//agentPath.POST("/user/info", agent.UserGameInfo)
			//agentPath.POST("/diamond/logs", agent.DiamondChangeLog)
		}

		//启动页数据
		v1.POST("/index", index.Star)

		//签到
		signPath := v1.Group("/sign", servermiddleware.BaseAuth())
		{
			//签到
			signPath.POST("/receive", sign.SigIn)
			//签到列表
			signPath.POST("/list", sign.SigList)

		}

		wheelPath := v1.Group("/wheel", servermiddleware.BaseAuth())
		{
			//大转盘列表
			wheelPath.POST("/rules", sign.BigWheelRules)
			//大转盘抽奖
			wheelPath.POST("/turntable", sign.BigWheelTurntable)
			//大转盘记录
			wheelPath.POST("/record", sign.BigWheelRecord)

		}

		rankPath := v1.Group("/rank", servermiddleware.BaseAuth())
		{
			//排行榜
			rankPath.POST("/list", rank.RanksList)
		}

		exchangePath := v1.Group("/exchange", servermiddleware.BaseAuth())
		{
			//获取兑换配置
			exchangePath.POST("/config", exchange.Config)
			//钻石兑换
			exchangePath.POST("/diamond", exchange.DiamondExchange)
			//钻石兑换记录
			exchangePath.POST("/diamond/record", exchange.DiamondExchangeRecord)
		}

		//公告列表
		v1.POST("/news/info", news.NewsInfo)
		//跑马灯
		v1.POST("/prize/info", news.RecordPrizeInfo)

		//活动
		v1.POST("/activitylist", activity.ActivityList)

		v1.POST("/game/version/list", game.GetAllGameVersionV2)

	}

	v2 := router.Group("/v2", servermiddleware.Base())
	{
		v2.POST("/game/version/list", game.GetAllGameVersionV2)
	}
	//微服务 api 网关
	router.POST("v1/user/signup", apigw.DoSignUpHandler)
	//获取服务器列表
	router.POST("v1/public/gethost", apigw.GetHostHandler)
}
