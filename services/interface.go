package services

import "github.com/gin-gonic/gin"

type ChatSessionService interface {
	SendMessageToBot(c *gin.Context)
	InitSession(c *gin.Context)
	ChatSendMessageMQ(c *gin.Context)
	GetUserList(c *gin.Context)
	ChatSessionSendMessageAll(c *gin.Context)
	SignUp(c *gin.Context)
	Login(c *gin.Context)
}

type UserManagement interface {
	GetList(c *gin.Context)
	GetOnlineList(c *gin.Context)
}
