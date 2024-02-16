package server

import "github.com/gin-gonic/gin"

type LLMBOT interface {
	SpeakToBot(c *gin.Context, messageMapSlice []map[string]string) (respMessage string)
	SpeakToBot_server(c *gin.Context)
}
type LLMTransceiver interface {
	SpeakToLLM(c *gin.Context, messageMapSlice []map[string]string) (respMessage string)
}

type UserOnline interface {
	UserLogin(accIp, accPort string, appId uint32, userId string, addr string,
		loginTime uint64) (userOnline UserOnline)
	Heartbeat(currentTime uint64)
	LogOut()
	IsOnline() (online bool)
}
