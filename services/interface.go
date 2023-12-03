package services

import "github.com/gin-gonic/gin"

type TalkFunction interface {
	SendMessageToBot(c *gin.Context)
	InitSession(c *gin.Context)
}

type UserManagement interface {
	GetList(c *gin.Context)
	GetOnlineList(c *gin.Context)
}
