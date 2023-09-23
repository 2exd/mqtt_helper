package structs

import (
	"mqtt_helper/log"
	"sync"
	"time"
)

// ClientInfo 保存客户端的信息
type ClientInfo struct {
	LastOnline time.Time
	IsOnline   bool
}

var (
	once              sync.Once
	MapMutex          sync.Mutex
	clientMapInstance *ClientMap
)

type ClientMap struct {
	Data map[string]ClientInfo
}

// CheckOnline 检查客户端在线情况
func CheckOnline() {
	cMap := GetClientMapInstance()
	for {
		// 如果客户端超过30秒没有发送消息，则标记为离线
		time.Sleep(5 * time.Second) // 每5秒检查一次
		MapMutex.Lock()
		for id, info := range cMap.Data {
			if !info.IsOnline {
				delete(cMap.Data, id)
				log.Logger.Infof("offline client %s is deleted", id)
			}
			if time.Since(info.LastOnline) > 30*time.Second && info.IsOnline {
				info.IsOnline = false
				cMap.Data[id] = info
				log.Logger.Infof("client %s offline", id)
			}
		}
		MapMutex.Unlock()
	}
}

func AddOrUpdateClient(clientName string) {
	// map 添加 client
	cMap := GetClientMapInstance()
	MapMutex.Lock()
	info, ok := cMap.Data[clientName]
	if !ok {
		// add
		newInfo := &ClientInfo{
			LastOnline: time.Now(),
			IsOnline:   true,
		}
		cMap.Data[clientName] = *newInfo
	} else {
		// update
		info.LastOnline = time.Now()
		info.IsOnline = true
		cMap.Data[clientName] = info
	}
	MapMutex.Unlock()
}

func DeleteClient(clientName string) {
	cMap := GetClientMapInstance()
	MapMutex.Lock()
	_, ok := cMap.Data[clientName]
	if !ok {

	} else {
		delete(cMap.Data, clientName)
	}
	MapMutex.Unlock()
}

// GetClientMapInstance 获取单例实例的函数
func GetClientMapInstance() *ClientMap {
	once.Do(func() {
		clientMapInstance = &ClientMap{
			Data: make(map[string]ClientInfo, 10),
		}
	})

	return clientMapInstance
}
