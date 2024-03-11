package forward

import "net"

type DataForwarder interface {
	Forward(sourceConn, destinationConn net.Conn) error
	ForwardWithNormal(sourceConn, destinationConn net.Conn) error
	ForwardWithTrafficControl(sourceConn, destinationConn net.Conn) error
}
