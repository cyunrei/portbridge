package forwarder

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net"
)

type Forwarder struct {
	SourceAddr      string
	DestinationAddr string
	Protocol        string
	DataForwarder   DataForwarder
}

func NewForwarder() *Forwarder {
	return &Forwarder{}
}

func (f *Forwarder) WithSourceAddr(sourceAddr string) *Forwarder {
	f.SourceAddr = sourceAddr
	return f
}

func (f *Forwarder) WithDestinationAddr(destinationAddr string) *Forwarder {
	f.DestinationAddr = destinationAddr
	return f
}

func (f *Forwarder) WithProtocol(protocol string) *Forwarder {
	f.Protocol = protocol
	return f
}

func (f *Forwarder) WithDataForwarder(dataForwarder DataForwarder) *Forwarder {
	f.DataForwarder = dataForwarder
	return f
}

func (f *Forwarder) Start() error {
	var err error
	switch f.Protocol {
	case "tcp":
		err = f.startTCP()
	case "udp":
		err = f.startUDP()
	default:
		return errors.New("unsupported protocol: " + f.Protocol)
	}
	return err
}

func (f *Forwarder) startTCP() error {
	localListener, err := net.Listen("tcp", f.SourceAddr)
	if err != nil {
		return fmt.Errorf("unable to bind to local TCP address: %s", err)
	}
	defer localListener.Close()

	log.Printf("TCP Port forwarding is active. Forwarding from %s to %s", f.SourceAddr, f.DestinationAddr)

	for {
		localConn, err := localListener.Accept()
		if err != nil {
			log.Errorf("Error accepting TCP connection: %s", err)
			continue
		}

		log.Printf("TCP connection established from %s", localConn.RemoteAddr())

		remoteConn, err := net.Dial("tcp", f.DestinationAddr)
		if err != nil {
			log.Warnf("Unable to connect to remote TCP address: %s", err)
			localConn.Close()
			log.Printf("TCP connection disconnected from %s", localConn.RemoteAddr())
			continue
		}
		go func() {
			f.DataForwarder.Forward(localConn, remoteConn)
			localConn.Close()
			log.Printf("TCP connection disconnected from %s", localConn.RemoteAddr())
			remoteConn.Close()
		}()
	}
}

func (f *Forwarder) startUDP() error {
	localUDPAddr, err := net.ResolveUDPAddr("udp", f.SourceAddr)
	if err != nil {
		return fmt.Errorf("error resolving local UDP address: %s", err)
	}
	localConn, err := net.ListenUDP("udp", localUDPAddr)
	if err != nil {
		return fmt.Errorf("unable to bind to local UDP address: %s", err)
	}
	defer localConn.Close()

	log.Printf("UDP Port forwarding is active. Forwarding from %s to %s", f.SourceAddr, f.DestinationAddr)

	remoteUDPAddr, err := net.ResolveUDPAddr("udp", f.DestinationAddr)
	if err != nil {
		return fmt.Errorf("error resolving remote UDP address: %s", err)
	}
	remoteConn, err := net.DialUDP("udp", nil, remoteUDPAddr)
	if err != nil {
		return fmt.Errorf("unable to connect to remote UDP address: %s", err)
	}
	defer remoteConn.Close()

	f.DataForwarder.Forward(&*localConn, &*remoteConn)

	return nil
}
