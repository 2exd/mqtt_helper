package structs

import (
	"bytes"
	"context"
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/tidwall/gjson"
	"image/png"
	"math/rand"
	"mqtt_helper/conf"
	"mqtt_helper/constants"
	"mqtt_helper/log"
	"mqtt_helper/utils"
	"os"
	"strings"
	"time"
)

var config *conf.Config

var Mc *MqttClient

type MqttClient struct {
	Client    mqtt.Client
	Opts      *mqtt.ClientOptions
	PubTopics []string
	SubTopics []string
	Ip        string
	Username  string
	Qos       int
}

var connectionLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	log.Logger.Errorf("onConnectionLost was called with error: %s", err)
	// 直接退出
	// os.Exit(1)
}

var reconnectingHandler mqtt.ReconnectHandler = func(client mqtt.Client, options *mqtt.ClientOptions) {
	log.Logger.Errorf("...... mqtt reconnecting ......")
}

var onMessageArrive mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	// MsgArr <- true
	log.Logger.Debugf("Message arrived! TOPIC: %s, MSG: %s", msg.Topic(), msg.Payload())
	go parseMessage(msg, Mc)
}

func parseMessage(msg mqtt.Message, c *MqttClient) {
	// 判断 json 是否合法
	if !gjson.Valid(string(msg.Payload())) {
		log.Logger.Errorf("json valid error!")
		return
	}

	// 使用 gjson 解析 JSON
	result := gjson.Parse(string(msg.Payload()))

	// 创建 Message 结构体并赋值
	mqttMsg := MqttMessage{
		MsgType:  MessageType(result.Get(constants.MSG_TYPE).Int()),
		OpCode:   OpCode(result.Get(constants.OP_CODE).Int()),
		Username: result.Get(constants.USERNAME).String(),
		IP:       result.Get(constants.IP).String(),
		Data:     result.Get(constants.DATA).String(),
	}
	log.Logger.Debugf("op code is %d", mqttMsg.MsgType)
	switch mqttMsg.MsgType {
	case 1:
		ConnectionControl(mqttMsg)
	case 2:
		go MsgTransfer(mqttMsg)
	default:
	}
}

// ConnectBroker 连接 mqtt 服务器
func (c *MqttClient) ConnectBroker() {
	if token := c.Client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	log.Logger.Debug("Client successfully connected to the broker!")
}

// ClientInit 设置 mqtt client 参数
func (c *MqttClient) ClientInit() {
	rand.Seed(time.Now().UnixNano())

	c.Opts.AddBroker(config.MqttAddr)

	c.Opts.SetKeepAlive(5 * time.Second)
	c.Opts.SetPingTimeout(30 * time.Second)

	ip, err := utils.GetNonLoopbackValidIPv4()
	if err != nil {
		log.Logger.Error(err)
		os.Exit(1)
	}

	c.Ip = ip
	c.Username = config.MqttUser
	c.Qos = config.MqttQos

	c.Opts.SetClientID(config.MqttUser + "@" + ip)

	// c.Opts.SetUsername(config.MqttUser)
	// c.Opts.SetPassword(config.MqttPassword)

	c.Opts.SetAutoReconnect(true)
	c.Opts.SetMaxReconnectInterval(10 * time.Second)
	c.Opts.SetConnectionLostHandler(connectionLostHandler)
	c.Opts.SetReconnectingHandler(reconnectingHandler)

	c.Opts.SetResumeSubs(true)
	c.Client = mqtt.NewClient(c.Opts)
	log.Logger.Infof("Client ID: %s", c.Opts.ClientID)
	log.Logger.Debug("Client initialization succeeded!")
}

// SubscribeTopics 订阅主题
func (c *MqttClient) SubscribeTopics() {
	log.Logger.Info("Subscribe topics is: ", c.SubTopics)
	topic := &strings.Builder{}
	for i := 0; i < len(c.SubTopics); i++ {
		switch c.SubTopics[i] {
		case constants.SERVER_TOPIC:
			// server/to/client/IP
			topic.WriteString(c.SubTopics[i])
			topic.WriteString("/")
			topic.WriteString(c.Ip)
			if token := c.Client.Subscribe(topic.String(), byte(config.MqttQos), onMessageArrive); token.Wait() && token.Error() != nil {
				log.Logger.Error(token.Error())
				os.Exit(1)
			}
			log.Logger.Info("Subscription succeeded! Subscribe Topic is: ", topic)
			if token := c.Client.Subscribe(constants.SERVER_ALL_CLIENT, byte(config.MqttQos), onMessageArrive); token.Wait() && token.Error() != nil {
				log.Logger.Error(token.Error())
				os.Exit(1)
			}
			log.Logger.Info("Subscription succeeded! Subscribe Topic is: ", constants.SERVER_ALL_CLIENT)
		default:
			if token := c.Client.Subscribe(c.SubTopics[i], byte(config.MqttQos), onMessageArrive); token.Wait() && token.Error() != nil {
				log.Logger.Error(token.Error())
				os.Exit(1)
			}
			log.Logger.Info("Subscription succeeded! Subscribe Topic is: ", c.SubTopics[i])
		}
	}
}

func (c *MqttClient) publishScreenshot(data []byte) {
	stringData := utils.ToBase64(data)
	text := &MqttMessage{
		MsgType:  MessageTransfer,
		OpCode:   SendScreenshot,
		Username: c.Username,
		IP:       c.Ip,
		Data:     stringData,
	}
	jsonMarshal, _ := json.Marshal(text)
	token := c.Client.Publish(constants.CLIENT_TOPIC, byte(c.Qos), false, jsonMarshal)
	log.Logger.Infof("Send topic %s, msg is picture, size=%d", constants.CLIENT_TOPIC, len(data))
	token.Wait()
}

func (c *MqttClient) publishCode(data string) {
	text := &MqttMessage{
		MsgType:  MessageTransfer,
		OpCode:   SendCode,
		Username: c.Username,
		IP:       c.Ip,
		Data:     data,
	}
	jsonMarshal, _ := json.Marshal(text)
	token := c.Client.Publish(constants.CLIENT_TOPIC, byte(c.Qos), false, jsonMarshal)
	log.Logger.Infof("Send topic %s, msg is: %s", constants.CLIENT_TOPIC, jsonMarshal)
	token.Wait()
}

func (c *MqttClient) Run(ctx context.Context) error {
	config = conf.GetConfig()
	c.ClientInit()
	c.ConnectBroker()
	c.SubscribeTopics()

	go utils.GetScreen()

loop:
	for {
		select {
		case img := <-utils.ImgChan:

			// 将 *image.RGBA 转换为字节数组
			var data bytes.Buffer
			err := png.Encode(&data, img)
			if err != nil {
				log.Logger.Fatal(err)
			}
			// data := utils.EncodeToBytes(img)
			go c.publishScreenshot(data.Bytes())
		case <-ctx.Done():
			// 跳出 retry 循环
			break loop
		}
	}
	return nil
}