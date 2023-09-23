package constants

import "strings"

var (
	// VERSION 版本
	VERSION      = "0.2.1"
	CLIENT_TOPIC = "client/to/server"
	SERVER_TOPIC = "server/to/client"

	SERVER_ALL_CLIENT = "server/all/client"

	SERVER = "server"
	CLIENT = "client"

	CONFIG = "config"
	YAML   = "yaml"

	TEMP_FILE = "./temp.txt"
)

var (
	OP_CODE  = "opCode"
	MSG_TYPE = "msgType"
	USERNAME = "username"
	IP       = "ip"
	DATA     = "data"

	LOGIN       = "login"
	PING        = "ping"
	PONG        = "pong"
	LOGOUT      = "logout"
	SERVER_DOWN = "server down"
)

var (
	MENU = "menu"
	ZERO = "0"
)

func GetNameIP(name, ip string) string {
	var sb strings.Builder
	sb.WriteString(name)
	sb.WriteString("_")
	sb.WriteString(ip)
	// eg:client_192.168.150.1
	return sb.String()
}

func GetNameAtIP(name, ip string) string {
	var sb strings.Builder
	sb.WriteString(name)
	sb.WriteString("@")
	sb.WriteString(ip)
	// eg:client@192.168.150.1
	return sb.String()
}

func GetNameIPTopic(nameIp string) string {
	var sb strings.Builder
	sb.WriteString(CLIENT)
	sb.WriteString("/")
	sb.WriteString(nameIp)
	sb.WriteString("/")
	sb.WriteString(SERVER)
	// eg:client/client_192.168.150.1/server
	return sb.String()
}

func GetClientNameIPTopic(nameIp string) string {
	var sb strings.Builder
	sb.WriteString(SERVER_TOPIC)
	sb.WriteString("/")
	sb.WriteString(nameIp)
	// eg:server/to/client/client_192.168.150.1
	return sb.String()
}
