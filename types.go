package main

import "net"

type TCPDataForwarder interface {
	Forward(sourceConn, destinationConn net.Conn)
}

type UDPDataForwarder interface {
	Forward(sourceConn, destinationConn net.UDPConn)
}
