package structs

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MqttClientBuilder = func(mqttClient *MqttClient)

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
	var client = &MqttClient{
		Opts: mqtt.NewClientOptions(),
	}
	for i := range builder {
		builder[i](client)
	}
	if len(client.PubTopics) == 0 {
		return nil, fmt.Errorf("missing publish topics")
	}
	if len(client.SubTopics) == 0 {
		return nil, fmt.Errorf("missing subscribe topics")
	}
	return client, nil
}
