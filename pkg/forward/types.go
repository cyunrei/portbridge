package forward

import "net"

type TCPDataForwarder interface {
	Forward(sourceConn, destinationConn net.Conn) error
}

type UDPDataForwarder interface {
	Forward(sourceConn, destinationConn net.UDPConn)
}
