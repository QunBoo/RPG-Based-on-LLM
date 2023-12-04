package server

import "github.com/gin-gonic/gin"

type BOT interface {
	//SpeakToBot(c *gin.Context, messageMap map[string]string)
	SpeakToBot(c *gin.Context, message map[string]string)
	SpeakToBot_server(c *gin.Context)
}
type LLMBOT interface {
	SpeakToBot(c *gin.Context, messageMapSlice []map[string]string) (respMessage string)
	SpeakToBot_server(c *gin.Context)
}
type LLMTransceiver interface {
	SpeakToLLM(c *gin.Context, messageMapSlice []map[string]string) (respMessage string)
}
