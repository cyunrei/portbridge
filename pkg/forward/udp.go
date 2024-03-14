package forward

import (
	"context"
	"golang.org/x/time/rate"
	"net"
	"time"
)

type UDPDataForwarder struct {
	BufferSize     uint64
	BandwidthLimit uint64
	Deadline       time.Duration
}

func NewUDPDataForwarder() *UDPDataForwarder {
	return &UDPDataForwarder{
		BandwidthLimit: DefaultBandwidthLimit,
		BufferSize:     DefaultUDPBufferSize,
		Deadline:       time.Duration(DefaultUDPDeadline),
	}
}

func (f *UDPDataForwarder) Forward(sourceConn, destinationConn net.Conn) error {
	if f.BandwidthLimit != DefaultBandwidthLimit {
		return f.ForwardWithTrafficControl(sourceConn, destinationConn)
	} else {
		return f.ForwardWithNormal(sourceConn, destinationConn)
	}
}

func (f *UDPDataForwarder) ForwardWithNormal(sourceConn, destinationConn net.Conn) error {
	sourceUDPConn, _ := sourceConn.(*net.UDPConn)
	destinationUDPConn, _ := destinationConn.(*net.UDPConn)
	f.forwardData(sourceUDPConn, destinationUDPConn, nil)
	return nil
}

func (f *UDPDataForwarder) ForwardWithTrafficControl(sourceConn, destinationConn net.Conn) error {
	sourceUDPConn, _ := sourceConn.(*net.UDPConn)
	destinationUDPConn, _ := destinationConn.(*net.UDPConn)

	limiter := rate.NewLimiter(rate.Limit(f.BandwidthLimit*1024/8), int(f.BandwidthLimit*1024/8))
	f.forwardData(sourceUDPConn, destinationUDPConn, limiter)
	return nil
}

func (f *UDPDataForwarder) forwardData(sourceUDPConn, destinationUDPConn *net.UDPConn, limiter *rate.Limiter) {
	sourceConnBuffer := make([]byte, f.BufferSize)
	for {
		sourceUDPConn.SetReadDeadline(time.Now().Add(f.Deadline * time.Second))
		n, sourceConnAddr, err := sourceUDPConn.ReadFromUDP(sourceConnBuffer)
		if err != nil {
			continue
		}

		data := make([]byte, n)
		copy(data, sourceConnBuffer[:n])

		go func(data []byte, sourceConnAddr *net.UDPAddr) {
			if limiter != nil {
				err := limiter.WaitN(context.Background(), n)
				if err != nil {
					return
				}
			}

			_, err := destinationUDPConn.Write(data)
			if err != nil {
				return
			}

			destinationConnBuffer := make([]byte, f.BufferSize)
			destinationUDPConn.SetReadDeadline(time.Now().Add(f.Deadline * time.Second))
			m, _, err := destinationUDPConn.ReadFromUDP(destinationConnBuffer)
			if err != nil {
				return
			}

			if limiter != nil {
				err := limiter.WaitN(context.Background(), m)
				if err != nil {
					return
				}
			}

			_, err = sourceUDPConn.WriteToUDP(destinationConnBuffer[:m], sourceConnAddr)
			if err != nil {
				return
			}
		}(data, sourceConnAddr)
	}
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
	f.Deadline = time.Duration(second)
	return f
}
