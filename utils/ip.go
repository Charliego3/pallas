package utils

import (
	"net"
)

func IntranetIp() (net.IP, bool) {
	if ip, ok := ExternalIp(); ok {
		return ip, true
	}

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, false
	}

	for _, addr := range addrs {
		ip := addr.(*net.IPNet).IP
		if ip == nil || ip.IsLoopback() {
			continue
		}

		ip = ip.To4()
		if ip == nil {
			continue
		}
		return ip, true
	}
	return nil, false
}

func ExternalIp() (net.IP, bool) {
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		return nil, false
	}
	addr := conn.LocalAddr().(*net.UDPAddr)
	return addr.IP, true
}
