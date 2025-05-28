package apis

import (
	"entry_task/middlewares"
	"entry_task/view/http/apis/message_apis"
	"entry_task/view/http/apis/user_apis"
	"github.com/gin-gonic/gin"
)

func SetUserApi(e *gin.RouterGroup) {
	userGroup := e.Group("/user")
	{
		userGroup.POST("/signup", user_apis.UserSignUp)
		userGroup.POST("/login", user_apis.UserLogIn)
	}

}

func SetMessageApi(e *gin.RouterGroup) {
	messageGroup := e.Group("/message").Use(middlewares.AuthMiddleware())
	{
		// 通用 查询
		messageGroup.GET("/:messageId/:commentPage", message_apis.ViewContent)
		messageGroup.GET("/explore/:postPage", message_apis.MessageExplore)
		// 管理员 创建消息&更新消息&删除消息
		messageGroup.POST("/new", message_apis.MessageNew)
		messageGroup.PUT("/:messageId/update", message_apis.MessageUpdate)
		messageGroup.PUT("/:messageId/delete", message_apis.RemoveMessage)
		// 普通用户 评论&TAG趋势
		messageGroup.POST("/:messageId/comment", message_apis.MessageComment)
		messageGroup.GET("/trend", message_apis.MessageTrend)
	}
}
