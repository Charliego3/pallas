package utils

import "net"

func RandomTCPListener() (net.Listener, error) {
	listner, err := net.Listen("tcp", ":0")
	return listner, err
}
