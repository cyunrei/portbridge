package forward

import (
	"github.com/fujiwara/shapeio"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
)

type SimpleTCPDataForwarder struct{}

func (f *SimpleTCPDataForwarder) Forward(sourceConn, destinationConn net.Conn) {
	go func() {
		_, err := io.Copy(sourceConn, destinationConn)
		if err != nil {
			log.Printf("Connection disconnted from %s\n", sourceConn.RemoteAddr())
		}
		sourceConn.Close()
		destinationConn.Close()
	}()
	go func() {
		io.Copy(destinationConn, sourceConn)
		sourceConn.Close()
		destinationConn.Close()
	}()
}

func NewSimpleTCPDataForwarder() *SimpleTCPDataForwarder {
	return &SimpleTCPDataForwarder{}
}

func NewTrafficControlTCPDataForwarder() *TrafficControlTCPDataForwarder {
	return &TrafficControlTCPDataForwarder{BandwidthLimit: DefaultTCPBandwidthLimit}
}

const DefaultTCPBandwidthLimit int64 = 0

type TrafficControlTCPDataForwarder struct {
	BandwidthLimit int64
}

func (f *TrafficControlTCPDataForwarder) SetBandwidthLimit(bandwidthLimit int64) *TrafficControlTCPDataForwarder {
	f.BandwidthLimit = bandwidthLimit
	return f
}

func (f *TrafficControlTCPDataForwarder) Forward(sourceConn, destinationConn net.Conn) {
	go func() {
		destConnReader := shapeio.NewReader(destinationConn)
		destConnReader.SetRateLimit(float64(1024 * f.BandwidthLimit))
		_, err := io.Copy(sourceConn, destConnReader)
		if err != nil {
			log.Printf("Connection disconnted from %s\n", sourceConn.RemoteAddr())
		}
		sourceConn.Close()
		destinationConn.Close()
	}()
	go func() {
		sourceConnReader := shapeio.NewReader(sourceConn)
		sourceConnReader.SetRateLimit(float64(1024 * f.BandwidthLimit))
		io.Copy(destinationConn, sourceConnReader)
		sourceConn.Close()
		destinationConn.Close()
	}()
}
