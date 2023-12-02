package serverimpl

import (
	"FantasticLife/server"
	"github.com/gin-gonic/gin"
)

type GptConn struct {
	Key       string
	EndPoint  string
	AppSecret string
}
type GptBot struct {
	conn    *GptConn
	chatMap map[string]string
}

func (b *GptBot) BOTChat(c *gin.Context) {

}
func (b *GptBot) BOTRemember(c *gin.Context) {

}

//func NewGptConn(Key, EndPoint, Appsecret string) *GptConn {
//	return &GptConn{
//		Key:       Key,
//		EndPoint:  EndPoint,
//		AppSecret: Appsecret,
//	}
//}

func NewGptBot(pConn *GptConn) (server.BOT, error) {
	return &GptBot{
		conn: pConn,
	}, nil
}
