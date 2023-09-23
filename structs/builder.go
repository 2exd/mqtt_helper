package structs

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var (
	ms *MqttServer
)

type MqttServer struct {
	*MqttClient
}
type MqttClientBuilder = func(mqttClient *MqttClient)

type MqttServerBuilder = func(mqttServer *MqttServer)

func NewMqttServer(builder ...MqttClientBuilder) (*MqttServer, error) {
	mc, _ := NewMqttClient(builder...)
	ms = &MqttServer{
		MqttClient: mc,
	}
	return ms, nil
}

func PubTopics(pubTopics []string) MqttClientBuilder {
	return func(client *MqttClient) {
		client.PubTopics = pubTopics
	}
}

func SubTopics(subTopics []string) MqttClientBuilder {
	return func(client *MqttClient) {
		client.SubTopics = subTopics
	}
}

func NewMqttClient(builder ...MqttClientBuilder) (*MqttClient, error) {
	mc = &MqttClient{
		Opts: mqtt.NewClientOptions(),
	}
	for i := range builder {
		builder[i](mc)
	}
	if len(mc.PubTopics) == 0 {
		return nil, fmt.Errorf("missing publish topics")
	}
	if len(mc.SubTopics) == 0 {
		return nil, fmt.Errorf("missing subscribe topics")
	}
	return mc, nil
}

func GetMqttServerInstance() *MqttServer {
	return ms
}

func GetMqttClientInstance() *MqttClient {
	return mc
}
