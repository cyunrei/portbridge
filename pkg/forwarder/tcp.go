package forwarder

import (
	"github.com/fujiwara/shapeio"
	"io"
	"net"
)

const DefaultBandwidthLimit uint64 = 0

type TCPDataForwarder struct {
	src     net.Conn
	dst     net.Conn
	bwLimit uint64
}

func NewTCPDataForwarder() *TCPDataForwarder {
	return &TCPDataForwarder{bwLimit: DefaultBandwidthLimit}
}

func (f *TCPDataForwarder) Start() error {
	done := make(chan error, 2)

	go func() {
		err := f.forward()
		if err != nil {
			done <- err
		}
	}()

	go func() {
		err := f.reply()
		if err != nil {
			done <- err
		}
	}()

	for i := 0; i < 2; i++ {
		err := <-done
		if err != nil {
			return err
		}
	}

	return nil
}

func (f *TCPDataForwarder) forward() error {
	var reader io.Reader
	if f.bwLimit > 0 {
		destConnReader := shapeio.NewReader(f.dst)
		destConnReader.SetRateLimit(float64(1024 * f.bwLimit))
		reader = destConnReader
	} else {
		reader = f.dst
	}
	_, err := io.Copy(f.src, reader)
	return err
}

func (f *TCPDataForwarder) reply() error {
	var reader io.Reader
	if f.bwLimit > 0 {
		srcConnReader := shapeio.NewReader(f.src)
		srcConnReader.SetRateLimit(float64(1024 * f.bwLimit))
		reader = srcConnReader
	} else {
		reader = f.src
	}
	_, err := io.Copy(f.dst, reader)
	return err
}

func (f *TCPDataForwarder) WithSrc(src net.Conn) *TCPDataForwarder {
	f.src = src
	return f
}

func (f *TCPDataForwarder) WithDst(dst net.Conn) *TCPDataForwarder {
	f.dst = dst
	return f
}

func (f *TCPDataForwarder) WithBandwidthLimit(bandwidth uint64) *TCPDataForwarder {
	f.bwLimit = bandwidth
	return f
}
