package utils

import (
	"net"
)

// GetServerIp:调用net.Interfaces实现获取本机的外部IP地址
// 问题：我在本地多网卡机器上，运行分布式场景，此函数返回的ip有误导致rpc连接失败。 遂google结果如下：
// 1、https://www.jianshu.com/p/301aabc06972
// 2、https://www.cnblogs.com/chaselogs/p/11301940.html
func GetServerIp() string {
	ip, err := externalIP()
	if err != nil {
		return ""
	}
	return ip.String()
}

// 通过net.Interfaces获取外部IP
func externalIP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			ip := getIpFromAddr(addr)
			if ip == nil {
				continue
			}
			return ip, nil
		}
	}
	return nil, err
}

func getIpFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil // not an ipv4 address
	}

	return ip
}
