package forward

import (
	"errors"
	"net"
	"time"
)

type UDPDataForwarder struct {
	BufferSize     uint64
	BandwidthLimit uint64
	DeadlineSecond time.Duration
}

func NewUDPDataForwarder() *UDPDataForwarder {
	return &UDPDataForwarder{
		BandwidthLimit: DefaultBandwidthLimit,
		BufferSize:     DefaultUDPBufferSize,
		DeadlineSecond: time.Duration(DefaultUDPDeadlineSecond),
	}
}

func (f *UDPDataForwarder) Forward(sourceConn, destinationConn net.Conn) error {
	return f.ForwardWithNormal(sourceConn, destinationConn)
}

func (f *UDPDataForwarder) ForwardWithNormal(sourceConn, destinationConn net.Conn) error {
	sourceUDPConn, _ := sourceConn.(*net.UDPConn)
	destinationUDPConn, _ := destinationConn.(*net.UDPConn)
	sourceConnBuffer := make([]byte, f.BufferSize)
	for {
		sourceConn.SetReadDeadline(time.Now().Add(f.DeadlineSecond * time.Second))
		n, sourceConnAddr, err := sourceUDPConn.ReadFromUDP(sourceConnBuffer)
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
			m, _, err := destinationUDPConn.ReadFromUDP(destinationConnBuffer)
			var netErr net.Error
			if errors.As(err, &netErr) && netErr.Timeout() {
				return
			}

			_, err = sourceUDPConn.WriteToUDP(destinationConnBuffer[:m], sourceConnAddr)
			if err != nil {
				return
			}

		}(data, sourceConnAddr)
	}
}

func (f *UDPDataForwarder) ForwardWithTrafficControl(sourceConn, destinationConn net.Conn) error {
	return f.ForwardWithNormal(sourceConn, destinationConn)
}

func (f *UDPDataForwarder) SetBandwidthLimit(bandwidthLimit uint64) *UDPDataForwarder {
	f.BandwidthLimit = bandwidthLimit
	return f
}

func (f *UDPDataForwarder) SetBufferSize(size uint64) *UDPDataForwarder {
	f.BufferSize = size
	return f
}

func (f *UDPDataForwarder) SetDeadlineSecond(second uint64) *UDPDataForwarder {
	f.DeadlineSecond = time.Duration(second)
	return f
}
