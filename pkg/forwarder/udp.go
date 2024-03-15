package forwarder

import (
	"context"
	"github.com/panjf2000/ants/v2"
	"golang.org/x/time/rate"
	"net"
	"time"
)

const (
	DefaultUDPBufferSize     uint64 = 1024
	DefaultUDPDeadlineTime   uint64 = 5
	DefaultGoroutinePoolSize        = 1000
)

type UDPDataForwarder struct {
	BufferSize     uint64
	BandwidthLimit uint64
	DeadlineTime   time.Duration
}

func NewUDPDataForwarder() *UDPDataForwarder {
	return &UDPDataForwarder{
		BandwidthLimit: DefaultBandwidthLimit,
		BufferSize:     DefaultUDPBufferSize,
		DeadlineTime:   time.Duration(DefaultUDPDeadlineTime),
	}
}

func (f *UDPDataForwarder) Forward(srcConn, dstConn net.Conn) error {
	if f.BandwidthLimit != DefaultBandwidthLimit {
		return f.ForwardWithTrafficControl(srcConn, dstConn)
	} else {
		return f.ForwardWithNormal(srcConn, dstConn)
	}
}

func (f *UDPDataForwarder) ForwardWithNormal(srcConn, dstConn net.Conn) error {
	f.forwardData(srcConn.(*net.UDPConn), dstConn.(*net.UDPConn), nil)
	return nil
}

func (f *UDPDataForwarder) ForwardWithTrafficControl(srcConn, dstConn net.Conn) error {
	limiter := rate.NewLimiter(rate.Limit(f.BandwidthLimit*1024/8), int(f.BandwidthLimit*1024/8))
	f.forwardData(srcConn.(*net.UDPConn), dstConn.(*net.UDPConn), limiter)
	return nil
}

func (f *UDPDataForwarder) forwardData(srcConn, dstConn *net.UDPConn, limiter *rate.Limiter) {
	pool, _ := ants.NewPool(DefaultGoroutinePoolSize)
	defer pool.Release()

	srcConnBuf := make([]byte, f.BufferSize)
	dstConnBuf := make([]byte, f.BufferSize)

	for {
		srcConn.SetReadDeadline(time.Now().Add(f.DeadlineTime * time.Second))
		n, srcAddr, err := srcConn.ReadFromUDP(srcConnBuf)
		if err != nil {
			continue
		}

		pool.Submit(func() {
			if limiter != nil {
				err := limiter.WaitN(context.Background(), n)
				if err != nil {
					return
				}
			}

			_, err := dstConn.Write(srcConnBuf)
			if err != nil {
				return
			}

			dstConn.SetReadDeadline(time.Now().Add(f.DeadlineTime * time.Second))
			m, _, err := dstConn.ReadFromUDP(dstConnBuf)
			if err != nil {
				return
			}

			if limiter != nil {
				err := limiter.WaitN(context.Background(), m)
				if err != nil {
					return
				}
			}

			_, err = srcConn.WriteToUDP(dstConnBuf[:m], srcAddr)
			if err != nil {
				return
			}
		})
	}
}

func (f *UDPDataForwarder) SetBandwidthLimit(limit uint64) *UDPDataForwarder {
	f.BandwidthLimit = limit
	return f
}

func (f *UDPDataForwarder) SetBufferSize(size uint64) *UDPDataForwarder {
	f.BufferSize = size
	return f
}

func (f *UDPDataForwarder) SetDeadlineTime(second uint64) *UDPDataForwarder {
	f.DeadlineTime = time.Duration(second)
	return f
}
