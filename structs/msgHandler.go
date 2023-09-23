package structs

import (
	"github.com/atotto/clipboard"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"mqtt_helper/constants"
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
	Login      OpCode = 1
	Logout     OpCode = 2
	Ping       OpCode = 3
	Pong       OpCode = 4
	ServerDown OpCode = 5

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

var responseCh = make(chan struct{})

func ConnectionControl(msg MqttMessage) {
	nameAtIP := constants.GetNameAtIP(msg.Username, msg.IP)
	nameIP := constants.GetNameIP(msg.Username, msg.IP)
	switch msg.OpCode {
	case Login:
		cMap := GetClientMapInstance()
		info, ok := cMap.Data[nameIP]
		// 订阅私有话题
		if !ok || !info.IsOnline {
			log.Logger.Infof("client login from %s", nameAtIP)
		}
		if !ok {
			SubscribePrivateTopic(nameIP, GetMqttServerInstance())
		}
		// map 添加 client
		AddOrUpdateClient(nameIP)

	case Logout:
		log.Logger.Infof("client logout from %s.", nameAtIP)
		// map 删除 client
		DeleteClient(nameIP)
		// 取消订阅私有话题
		UnsubscribePrivateTopic(nameIP, GetMqttServerInstance())

	case Ping:
		log.Logger.Debugf("receive ping from %s.", nameAtIP)
		AddOrUpdateClient(nameIP)
		ms := GetMqttServerInstance()
		go ms.PublishPong(nameIP)

	case Pong:
		log.Logger.Debugf("receive pong from %s.", nameAtIP)
		// 发送回应通知到通道
		responseCh <- struct{}{}

	case ServerDown:
		log.Logger.Infof("Server down, receive msg from %s.", nameAtIP)

	default:
		log.Logger.Infof("暂不支持的消息类型！type=%d, opCode = %d", ConnectControl, msg.OpCode)
	}
}

func MsgTransfer(msg MqttMessage) {
	nameAtIP := constants.GetNameAtIP(msg.Username, msg.IP)
	switch msg.OpCode {
	case SendScreenshot:
		log.Logger.Infof("receive screenshot from %s, size = %d", nameAtIP, len(msg.Data))
		rgba, _ := utils.DecodeImageFromBase64(msg.Data)
		utils.SaveScreen(rgba, constants.GetNameIP(msg.Username, msg.IP))

	case SendCode:
		log.Logger.Infof("receive code from %s. message:\n%s", nameAtIP, msg.Data)
		clipboard.WriteAll(msg.Data)

	case SendClipboard:
		log.Logger.Infof("receive code from %s. message:\n%s", nameAtIP, msg.Data)
		fileName, err := utils.SaveClipBoard(msg.Data, constants.GetNameIP(msg.Username, msg.IP))
		if err != nil {
			log.Logger.Error(err)
		}
		log.Logger.Infof("append to %s\n", fileName)

	default:
		log.Logger.Infof("暂不支持的消息类型！from %s, type=%d, opCode = %d", nameAtIP, MessageTransfer, msg.OpCode)
	}
}

func SubscribePrivateTopic(nameIP string, s *MqttServer) {
	topic := constants.GetNameIPTopic(nameIP)
	log.Logger.Info("Subscribe topics is: ", topic)
	if token := s.Client.Subscribe(topic, byte(config.MqttQos), func(client mqtt.Client, msg mqtt.Message) {
		log.Logger.Debugf("Message arrived! TOPIC: %s, MSG: %s", msg.Topic(), msg.Payload())
		go parseMessage(msg)
	}); token.Wait() && token.Error() != nil {
		log.Logger.Error(token.Error())
	}
	log.Logger.Info("Subscription succeeded! Subscribe Topic is: ", topic)
}

func UnsubscribePrivateTopic(nameIP string, s *MqttServer) {
	topic := constants.GetNameIPTopic(nameIP)
	log.Logger.Info("Unsubscribe topics is: ", topic)
	if token := s.Client.Unsubscribe(topic); token.Wait() && token.Error() != nil {
		log.Logger.Error(token.Error())
	}
	log.Logger.Info("Unsubscription succeeded! Subscribe Topic is: ", topic)
}
