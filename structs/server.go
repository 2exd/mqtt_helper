package structs

import (
	"context"
	"encoding/json"
	"mqtt_helper/conf"
	"mqtt_helper/constants"
	"mqtt_helper/log"
	"mqtt_helper/utils"
	"time"
)

var lastMod time.Time

var broadcast = make(chan []byte, 10)

type MqttServer struct {
	*MqttClient
}

var Ms *MqttServer

func (c *MqttServer) publishCode(data string) {
	text := &MqttMessage{
		MsgType:  MessageTransfer,
		OpCode:   SendCode,
		Username: c.Username,
		IP:       c.Ip,
		Data:     data,
	}
	jsonMarshal, _ := json.Marshal(text)
	token := c.Client.Publish(constants.SERVER_ALL_CLIENT, byte(c.Qos), false, jsonMarshal)
	log.Logger.Infof("Send topic %s, msg is: %s", constants.SERVER_ALL_CLIENT, jsonMarshal)
	token.Wait()
}

func (c *MqttServer) Run(ctx context.Context) error {
	config = conf.GetConfig()
	c.ClientInit()
	c.ConnectBroker()
	c.SubscribeTopics()

	lastMod = time.Now()

	fileTicker := time.NewTicker(5 * time.Second)
	defer fileTicker.Stop()
	var err error

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
		case <-fileTicker.C:
			lastMod, err = utils.ReadFileIfModified(lastMod, constants.TEMP_FILE)
			if err != nil {
				log.Logger.Error(err)
			}
		case message := <-broadcast:
			// 广播推送消息
			go c.publishCode(string(message))
		case <-ctx.Done():
			// 跳出 retry 循环
			break loop
		}
	}
	return nil
}
