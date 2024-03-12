package forward

import (
	"github.com/fujiwara/shapeio"
	"io"
	"net"
)

type TCPDataForwarder struct {
	BandwidthLimit uint64
}

func NewTCPDataForwarder() *TCPDataForwarder {
	return &TCPDataForwarder{BandwidthLimit: DefaultBandwidthLimit}
}

func (f *TCPDataForwarder) Forward(sourceConn, destinationConn net.Conn) error {
	if f.BandwidthLimit != DefaultBandwidthLimit {
		return f.ForwardWithTrafficControl(sourceConn, destinationConn)
	} else {
		return f.ForwardWithNormal(sourceConn, destinationConn)
	}
}

func (f *TCPDataForwarder) ForwardWithNormal(sourceConn, destinationConn net.Conn) error {
	return f.forwardData(sourceConn, destinationConn, 0)
}

func (f *TCPDataForwarder) ForwardWithTrafficControl(sourceConn, destinationConn net.Conn) error {
	return f.forwardData(sourceConn, destinationConn, float64(f.BandwidthLimit))
}

func (f *TCPDataForwarder) forwardData(sourceConn, destinationConn net.Conn, rateLimit float64) error {
	done := make(chan *ForwardingError, 2)

	go func() {
		var reader io.Reader
		if rateLimit > 0 {
			destConnReader := shapeio.NewReader(destinationConn)
			destConnReader.SetRateLimit(1024 * rateLimit)
			reader = destConnReader
		} else {
			reader = destinationConn
		}

		_, err := io.Copy(sourceConn, reader)
		done <- NewError(err, sourceConn, destinationConn, true)
	}()

	go func() {
		var reader io.Reader
		if rateLimit > 0 {
			sourceConnReader := shapeio.NewReader(sourceConn)
			sourceConnReader.SetRateLimit(1024 * rateLimit)
			reader = sourceConnReader
		} else {
			reader = sourceConn
		}

		_, err := io.Copy(destinationConn, reader)
		done <- NewError(err, sourceConn, destinationConn, false)
	}()

	for i := 0; i < 2; i++ {
		e := <-done
		if e != nil && e.Err != nil {
			return e
		}
	}

	return nil
}

func (f *TCPDataForwarder) SetBandwidthLimit(bandwidthLimit uint64) *TCPDataForwarder {
	f.BandwidthLimit = bandwidthLimit
	return f
}
