package server

import "github.com/gin-gonic/gin"

type BOT interface {
	BOTChat(c *gin.Context)
	BOTRemember(c *gin.Context)
}
