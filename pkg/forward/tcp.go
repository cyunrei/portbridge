package forward

import (
	"github.com/fujiwara/shapeio"
	"io"
	"net"
)

type SimpleTCPDataForwarder struct{}

func (f *SimpleTCPDataForwarder) Forward(sourceConn, destinationConn net.Conn) error {
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

func NewSimpleTCPDataForwarder() *SimpleTCPDataForwarder {
	return &SimpleTCPDataForwarder{}
}

func NewTrafficControlTCPDataForwarder() *TrafficControlTCPDataForwarder {
	return &TrafficControlTCPDataForwarder{BandwidthLimit: DefaultTCPBandwidthLimit}
}

const DefaultTCPBandwidthLimit uint64 = 0

type TrafficControlTCPDataForwarder struct {
	BandwidthLimit uint64
}

func (f *TrafficControlTCPDataForwarder) SetBandwidthLimit(bandwidthLimit uint64) *TrafficControlTCPDataForwarder {
	f.BandwidthLimit = bandwidthLimit
	return f
}

func (f *TrafficControlTCPDataForwarder) Forward(sourceConn, destinationConn net.Conn) error {
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
