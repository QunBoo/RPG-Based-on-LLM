package utils

import "encoding/json"

const (
	MessageTypeText = "text"

	MessageCmdMsg   = "msg"
	MessageCmdEnter = "enter"
	MessageCmdExit  = "exit"
)

// 消息的定义
type Message struct {
	Target string `json:"target"` // 目标
	Type   string `json:"type"`   // 消息类型 text/img/
	Msg    string `json:"msg"`    // 消息内容
	From   string `json:"from"`   // 发送者
}

func NewTestMsg(from string, Msg string) (message *Message) {

	message = &Message{
		Type: MessageTypeText,
		From: from,
		Msg:  Msg,
	}

	return
}

func getTextMsgData(cmd, uuId, msgId, message string) string {
	textMsg := NewTestMsg(uuId, message)
	head := NewResponseHead(msgId, cmd, OK, "Ok", textMsg)

	return head.String()
}

// 文本消息
func GetMsgData(uuId, msgId, cmd, message string) string {

	return getTextMsgData(cmd, uuId, msgId, message)
}

// 文本消息
func GetTextMsgData(uuId, msgId, message string) string {

	return getTextMsgData("msg", uuId, msgId, message)
}

// 用户进入消息
func GetTextMsgDataEnter(uuId, msgId, message string) string {

	return getTextMsgData("enter", uuId, msgId, message)
}

// 用户退出消息
func GetTextMsgDataExit(uuId, msgId, message string) string {

	return getTextMsgData("exit", uuId, msgId, message)
}

/************************  响应数据  **************************/
type Head struct {
	Seq      string    `json:"seq"`      // 消息的Id
	Cmd      string    `json:"cmd"`      // 消息的cmd 动作
	Response *Response `json:"response"` // 消息体
}

type Response struct {
	Code    uint32      `json:"code"`
	CodeMsg string      `json:"codeMsg"`
	Data    interface{} `json:"data"` // 数据 json
}

// push 数据结构体
type PushMsg struct {
	Seq  string `json:"seq"`
	Uuid uint64 `json:"uuid"`
	Type string `json:"type"`
	Msg  string `json:"msg"`
}

// 设置返回消息
func NewResponseHead(seq string, cmd string, code uint32, codeMsg string, data interface{}) *Head {
	response := NewResponse(code, codeMsg, data)

	return &Head{Seq: seq, Cmd: cmd, Response: response}
}

func (h *Head) String() (headStr string) {
	headBytes, _ := json.Marshal(h)
	headStr = string(headBytes)

	return
}

func NewResponse(code uint32, codeMsg string, data interface{}) *Response {
	return &Response{Code: code, CodeMsg: codeMsg, Data: data}
}
