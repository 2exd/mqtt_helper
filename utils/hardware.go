package utils

import (
	"fmt"
	"net"
	"strings"
)

func GetNonLoopbackValidIPv4() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range interfaces {
		// 过滤掉回环接口
		if iface.Flags&net.FlagLoopback == 0 {
			addrs, err := iface.Addrs()
			if err != nil {
				return "", err
			}

			for _, addr := range addrs {
				ipNet, ok := addr.(*net.IPNet)
				if ok && !ipNet.IP.IsLoopback() {
					ip := ipNet.IP
					// 检查是否为 169.254.x.x
					if !ip.IsGlobalUnicast() || strings.HasPrefix(ip.String(), "169.254.") {
						continue
					}
					if ip.To4() != nil {
						return ip.String(), nil
					}
				}
			}
		}
	}
	return "", fmt.Errorf("No non-loopback IPv4 address found")
}
