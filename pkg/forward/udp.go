package forward

import (
	"errors"
	"net"
	"time"
)

const DefaultUDPBufferSize uint64 = 1024
const DefaultUDPDeadlineSecond uint64 = 5

type SimpleUDPDataForwarder struct {
	BufferSize     uint64
	DeadlineSecond time.Duration
}

func NewSimpleUDPDataForwarder() *SimpleUDPDataForwarder {
	return &SimpleUDPDataForwarder{
		BufferSize:     DefaultUDPBufferSize,
		DeadlineSecond: time.Duration(DefaultUDPDeadlineSecond),
	}
}

func (f *SimpleUDPDataForwarder) SetBufferSize(size uint64) *SimpleUDPDataForwarder {
	f.BufferSize = size
	return f
}

func (f *SimpleUDPDataForwarder) SetDeadlineSecond(second uint64) *SimpleUDPDataForwarder {
	f.DeadlineSecond = time.Duration(second)
	return f
}

func (f *SimpleUDPDataForwarder) Forward(sourceConn, destinationConn net.UDPConn) {
	sourceConnBuffer := make([]byte, f.BufferSize)
	for {
		sourceConn.SetReadDeadline(time.Now().Add(f.DeadlineSecond * time.Second))
		n, sourceConnAddr, err := sourceConn.ReadFromUDP(sourceConnBuffer)
		if err != nil {
			continue
		}

		data := make([]byte, n)
		copy(data, sourceConnBuffer[:n])

		go func(data []byte, sourceConnAddr *net.UDPAddr) {
			_, err := destinationConn.Write(data)
			if err != nil {
				return
			}

			destinationConnBuffer := make([]byte, f.BufferSize)
			destinationConn.SetReadDeadline(time.Now().Add(f.DeadlineSecond * time.Second))
			m, _, err := destinationConn.ReadFromUDP(destinationConnBuffer)
			var netErr net.Error
			if errors.As(err, &netErr) && netErr.Timeout() {
				return
			}

			_, err = sourceConn.WriteToUDP(destinationConnBuffer[:m], sourceConnAddr)
			if err != nil {
				return
			}

		}(data, sourceConnAddr)
	}
}
