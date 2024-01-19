package utility

import "net"

func RandomTCPListener() (net.Listener, error) {
	listener, err := net.Listen("tcp", ":0")
	return listener, err
}
