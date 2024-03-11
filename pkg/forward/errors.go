package forward

import (
	"fmt"
	"net"
)

type ForwardingError struct {
	Err         error
	Source      net.Conn
	Destination net.Conn
	IsSourceErr bool
}

func (fe *ForwardingError) Error() string {
	direction := "destination"
	if fe.IsSourceErr {
		direction = "source"
	}
	return fmt.Sprintf("forward error on %s connection: %v", direction, fe.Err)
}

func NewError(err error, src, dst net.Conn, isSrcErr bool) *ForwardingError {
	return &ForwardingError{
		Err:         err,
		Source:      src,
		Destination: dst,
		IsSourceErr: isSrcErr,
	}
}
