package sd

import (
	"fmt"
	"net"
	"testing"
)

func TestIP(t *testing.T) {
	// 获取本机所有网络接口的 IP 地址
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	// 输出 IP 地址
	ipAddrs := make([]string, 0)
	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		if ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			fmt.Println(ipNet.IP.String())
			ipAddrs = append(ipAddrs, ipNet.IP.String())
		}
	}
	fmt.Printf(ipAddrs[0])
}
