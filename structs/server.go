package structs

import (
	"context"
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"mqtt_helper/conf"
	"mqtt_helper/constants"
	"mqtt_helper/log"
	"mqtt_helper/utils"
	"time"
)

var (
	lastMod   time.Time
	broadcast = make(chan []byte, 10)
	// 接收控制台输入
	consoleInput = make(chan string)
)

var onMessageArriveServer mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	// MsgArr <- true
	log.Logger.Debugf("Message arrived! TOPIC: %s, MSG: %s", msg.Topic(), msg.Payload())
	go parseMessage(msg)
}

func (s *MqttServer) PublishCode(data string) {
	text := &MqttMessage{
		MsgType:  MessageTransfer,
		OpCode:   SendCode,
		Username: s.Username,
		IP:       s.Ip,
		Data:     data,
	}
	jsonMarshal, _ := json.Marshal(text)
	token := s.Client.Publish(constants.SERVER_ALL_CLIENT, byte(s.Qos), false, jsonMarshal)
	log.Logger.Infof("Send topic %s, msg is: %s", constants.SERVER_ALL_CLIENT, jsonMarshal)
	token.Wait()
}

func (s *MqttServer) PublishPong(baseTopic string) {
	clientPongTopic := constants.GetClientNameIPTopic(baseTopic)
	text := &MqttMessage{
		MsgType:  ConnectControl,
		OpCode:   Pong,
		Username: s.Username,
		IP:       s.Ip,
		Data:     constants.PONG,
	}
	jsonMarshal, _ := json.Marshal(text)
	token := s.Client.Publish(clientPongTopic, byte(s.Qos), false, jsonMarshal)
	log.Logger.Debugf("Pong msg. Send topic %s, OpCode is %d", clientPongTopic, Pong)
	token.Wait()
}

func (s *MqttServer) PublishServerDown() {
	cMap := GetClientMapInstance()
	for k, _ := range cMap.Data {
		clientPongTopic := constants.GetClientNameIPTopic(k)
		text := &MqttMessage{
			MsgType:  ConnectControl,
			OpCode:   ServerDown,
			Username: s.Username,
			IP:       s.Ip,
			Data:     constants.SERVER_DOWN,
		}
		jsonMarshal, _ := json.Marshal(text)
		token := s.Client.Publish(clientPongTopic, byte(s.Qos), false, jsonMarshal)
		log.Logger.Infof("server down msg. Send topic %s, OpCode is %d", clientPongTopic, ServerDown)
		token.Wait()
	}
}

func (s *MqttServer) Run(ctx context.Context) error {
	config = conf.GetConfig()
	s.ClientInit()
	s.ConnectBroker()
	s.SubscribeTopics()

	lastMod = time.Now()

	fileTicker := time.NewTicker(5 * time.Second)
	defer fileTicker.Stop()
	var err error

	go CheckOnline()

	go func() {
		for {
			select {
			case data := <-utils.FileChange:
				// 将消息发送到广播通道，以便它可以被广播到所有客户端
				broadcast <- data
			}
		}
	}()

loop:
	for {
		select {
		case <-ctx.Done():
			// 跳出 retry 循环
			break loop
		case <-fileTicker.C:
			// 文件监控
			lastMod, err = utils.ReadFileIfModified(lastMod, constants.TEMP_FILE)
			if err != nil {
				log.Logger.Error(err)
			}
		case message := <-broadcast:
			// 广播推送消息
			go s.PublishCode(string(message))
		case console := <-consoleInput:
			// 处理控制台消息
			log.Logger.Info(console)
			go processConsole(console)
		}
	}
	return nil
}

func processConsole(str string) {
	switch str {
	case constants.ZERO, constants.MENU:

	}
}

func menu() {

}
