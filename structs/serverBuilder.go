package structs

type MqttServerBuilder = func(mqttServer *MqttServer)

func NewMqttServer(builder ...MqttClientBuilder) (*MqttServer, error) {
	Mc, _ := NewMqttClient(builder...)
	Ms = &MqttServer{
		MqttClient: Mc,
	}
	return Ms, nil
}
