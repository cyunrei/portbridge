package forwarder

import (
	"github.com/fujiwara/shapeio"
	"io"
	"net"
)

const DefaultBandwidthLimit uint64 = 0

type TCPDataForwarder struct {
	BandwidthLimit uint64
}

func NewTCPDataForwarder() *TCPDataForwarder {
	return &TCPDataForwarder{BandwidthLimit: DefaultBandwidthLimit}
}

func (f *TCPDataForwarder) Forward(srcConn, dstConn net.Conn) error {
	if f.BandwidthLimit != DefaultBandwidthLimit {
		return f.ForwardWithTrafficControl(srcConn, dstConn)
	} else {
		return f.ForwardWithNormal(srcConn, dstConn)
	}
}

func (f *TCPDataForwarder) ForwardWithNormal(srcConn, dstConn net.Conn) error {
	return f.forwardData(srcConn, dstConn, 0)
}

func (f *TCPDataForwarder) ForwardWithTrafficControl(srcConn, dstConn net.Conn) error {
	return f.forwardData(srcConn, dstConn, float64(f.BandwidthLimit))
}

func (f *TCPDataForwarder) forwardData(srcConn, dstConn net.Conn, rateLimit float64) error {
	done := make(chan error, 2)

	go func() {
		var reader io.Reader
		if rateLimit > 0 {
			destConnReader := shapeio.NewReader(dstConn)
			destConnReader.SetRateLimit(1024 * rateLimit)
			reader = destConnReader
		} else {
			reader = dstConn
		}

		_, err := io.Copy(srcConn, reader)
		done <- err
	}()

	go func() {
		var reader io.Reader
		if rateLimit > 0 {
			srcConnReader := shapeio.NewReader(srcConn)
			srcConnReader.SetRateLimit(1024 * rateLimit)
			reader = srcConnReader
		} else {
			reader = srcConn
		}

		_, err := io.Copy(dstConn, reader)
		done <- err
	}()

	for i := 0; i < 2; i++ {
		err := <-done
		if err != nil {
			return err
		}
	}

	return nil
}

func (f *TCPDataForwarder) SetBandwidthLimit(limit uint64) *TCPDataForwarder {
	f.BandwidthLimit = limit
	return f
}
