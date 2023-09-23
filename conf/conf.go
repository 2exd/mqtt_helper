package conf

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"sync"
)

type ConfigStruct struct {
	MqttAddr     string
	MqttQos      int
	MqttUser     string
	MqttPassword string
	SubTopics    []string
	PubTopics    []string
}

var (
	AppConfig *ConfigStruct
	once      sync.Once
)

func initializeConfig() {
	// 从配置文件中读取配置项
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	// 初始化 AppConfig
	AppConfig = &ConfigStruct{
		MqttAddr:     viper.GetString("mqtt.address"),
		MqttQos:      viper.GetInt("mqtt.qos"),
		MqttUser:     viper.GetString("mqtt.user"),
		MqttPassword: viper.GetString("mqtt.password"),
		SubTopics:    viper.GetStringSlice("mqtt.subTopics"),
		PubTopics:    viper.GetStringSlice("mqtt.pubTopics"),
	}
}

func GetConfig() *ConfigStruct {
	once.Do(func() {
		initializeConfig()
	})
	return AppConfig
}

func PrintConfig() {
	log.Printf("MqttAddr: %s\n", AppConfig.MqttAddr)

	fmt.Printf("MqttAddr: %s\n", AppConfig.MqttAddr)
	fmt.Printf("MqttQos: %d\n", AppConfig.MqttQos)
	fmt.Printf("MqttUser: %s\n", AppConfig.MqttUser)
	fmt.Printf("MqttPassword: %s\n", AppConfig.MqttPassword)
	fmt.Printf("SubTopics: %v\n", AppConfig.SubTopics)
	fmt.Printf("PubTopics: %v\n", AppConfig.PubTopics)
}
