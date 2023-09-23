package structs

import (
	"github.com/atotto/clipboard"
	"mqtt_helper/log"
	"mqtt_helper/utils"
)

// MessageType 表示消息类型
type MessageType int

const (
	ConnectControl  MessageType = 1
	MessageTransfer MessageType = 2
)

// OpCode 表示操作代码
type OpCode int

const (
	// 连接控制消息操作代码
	Login  OpCode = 1
	Logout OpCode = 2
	Ping   OpCode = 3
	Pong   OpCode = 4

	// 消息传输消息操作代码
	SendScreenshot OpCode = 1
	SendCode       OpCode = 2
	SendClipboard  OpCode = 3
)

// MqttMessage 表示交互协议中的消息
type MqttMessage struct {
	MsgType  MessageType `json:"msgType"`
	OpCode   OpCode      `json:"opCode"`
	Username string      `json:"username"`
	IP       string      `json:"ip"`
	Data     string      `json:"data"`
}

func ConnectionControl(msg MqttMessage) {

}

func MsgTransfer(msg MqttMessage) {
	switch msg.OpCode {
	case SendScreenshot:
		log.Logger.Infof("receive screenshot from %s@%s, size = %d", msg.Username, msg.IP, len(msg.Data))
		rgba, _ := utils.DecodeImageFromBase64(msg.Data)
		utils.SaveScreen(rgba, msg.Username+"_"+msg.IP)
	case SendCode:
		log.Logger.Infof("receive code from %s@%s. message:\n%s", msg.Username, msg.IP, msg.Data)
		clipboard.WriteAll(msg.Data)
	case SendClipboard:
		log.Logger.Infof("receive code from %s@%s. message:\n%s", msg.Username, msg.IP, msg.Data)
		fileName, err := utils.SaveClipBoard(msg.Data, msg.Username+"_"+msg.IP)
		if err != nil {
			log.Logger.Error(err)
		}
		log.Logger.Infof("%s\n", fileName)

	default:
		log.Logger.Infof("暂不支持的消息类型！opCode = %d", msg.OpCode)
	}
}
