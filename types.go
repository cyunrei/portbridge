package main

import "net"

type TCPDataForwarder interface {
	Forward(sourceConn, destinationConn net.Conn)
}
