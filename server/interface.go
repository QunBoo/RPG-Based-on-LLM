package server

import "github.com/gin-gonic/gin"

type BOT interface {
	//SpeakToBot(c *gin.Context, messageMap map[string]string)
	InitBot(c *gin.Context)
	SpeakToBot(c *gin.Context, message map[string]string)
	SpeakToBot_server(c *gin.Context)
}
