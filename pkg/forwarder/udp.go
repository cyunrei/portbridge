package forwarder

import (
	"context"
	"errors"
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
	src      net.UDPConn
	dst      net.UDPConn
	bufSize  uint64
	bwLimit  uint64
	deadline time.Duration
}

func NewUDPDataForwarder() *UDPDataForwarder {
	return &UDPDataForwarder{
		bwLimit:  DefaultBandwidthLimit,
		bufSize:  DefaultUDPBufferSize,
		deadline: time.Duration(DefaultUDPDeadlineTime),
	}
}

func (f *UDPDataForwarder) Start() error {
	return f.forward()
}

func (f *UDPDataForwarder) forward() error {
	var limiter *rate.Limiter
	if f.bwLimit != DefaultBandwidthLimit {
		limiter = rate.NewLimiter(rate.Limit(f.bwLimit*1024), int(f.bwLimit*1024))
	}
	pool, _ := ants.NewPool(DefaultGoroutinePoolSize)
	defer pool.Release()

	srcConnBuf := make([]byte, f.bufSize)
	dstConnBuf := make([]byte, f.bufSize)

	for {
		f.src.SetReadDeadline(time.Now().Add(f.deadline * time.Second))
		n, srcAddr, err := f.src.ReadFromUDP(srcConnBuf)
		if err != nil {
			continue
		}
		if limiter != nil {
			err := limiter.WaitN(context.Background(), n)
			if err != nil {
				continue
			}
		}

		pool.Submit(func() {
			_, err := f.dst.Write(srcConnBuf)
			if err != nil {
				return
			}

			f.dst.SetReadDeadline(time.Now().Add(f.deadline * time.Second))
			m, _, err := f.dst.ReadFromUDP(dstConnBuf)
			if err != nil {
				return
			}

			if limiter != nil {
				err := limiter.WaitN(context.Background(), m)
				if err != nil {
					return
				}
			}

			_, err = f.src.WriteToUDP(dstConnBuf[:m], srcAddr)
			if err != nil {
				return
			}
		})
	}
}

func (f *UDPDataForwarder) reply() error {
	return errors.New("UDP does not require a separate reply method due to its forwarding mechanism")
}

func (f *UDPDataForwarder) WithSrc(src net.UDPConn) *UDPDataForwarder {
	f.src = src
	return f
}

func (f *UDPDataForwarder) WithDst(dst net.UDPConn) *UDPDataForwarder {
	f.dst = dst
	return f
}

func (f *UDPDataForwarder) WithBandwidthLimit(bandwidth uint64) *UDPDataForwarder {
	f.bwLimit = bandwidth
	return f
}

func (f *UDPDataForwarder) WithBufferSize(size uint64) *UDPDataForwarder {
	f.bufSize = size
	return f
}

func (f *UDPDataForwarder) WithDeadlineTime(second uint64) *UDPDataForwarder {
	f.deadline = time.Duration(second)
	return f
}
