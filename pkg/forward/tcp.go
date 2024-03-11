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
	done := make(chan *ForwardingError, 2)

	go func() {
		_, err := io.Copy(destinationConn, sourceConn)
		done <- NewError(err, sourceConn, destinationConn, true)
	}()

	go func() {
		_, err := io.Copy(sourceConn, destinationConn)
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

func (f *TCPDataForwarder) ForwardWithTrafficControl(sourceConn, destinationConn net.Conn) error {
	done := make(chan *ForwardingError, 2)

	go func() {
		destConnReader := shapeio.NewReader(destinationConn)
		destConnReader.SetRateLimit(float64(1024 * f.BandwidthLimit))
		_, err := io.Copy(sourceConn, destConnReader)
		done <- NewError(err, sourceConn, destinationConn, true)
	}()

	go func() {
		sourceConnReader := shapeio.NewReader(sourceConn)
		sourceConnReader.SetRateLimit(float64(1024 * f.BandwidthLimit))
		_, err := io.Copy(destinationConn, sourceConnReader)
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
