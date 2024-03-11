package forward

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net"
)

type ForwardingConfig struct {
	SourceAddr       string
	DestinationAddr  string
	Protocol         string
	TCPDataForwarder TCPDataForwarder
	UDPDataForwarder UDPDataForwarder
}

func NewForwardingConfig() *ForwardingConfig {
	return &ForwardingConfig{}
}

func (f *ForwardingConfig) WithSourceAddr(sourceAddr string) *ForwardingConfig {
	f.SourceAddr = sourceAddr
	return f
}

func (f *ForwardingConfig) WithDestinationAddr(destinationAddr string) *ForwardingConfig {
	f.DestinationAddr = destinationAddr
	return f
}

func (f *ForwardingConfig) WithProtocol(protocol string) *ForwardingConfig {
	f.Protocol = protocol
	return f
}

func (f *ForwardingConfig) WithTCPDataForwarder(tcpDataForwarder TCPDataForwarder) *ForwardingConfig {
	f.TCPDataForwarder = tcpDataForwarder
	return f
}

func (f *ForwardingConfig) WithUDPDataForwarder(udpDataForwarder UDPDataForwarder) *ForwardingConfig {
	f.UDPDataForwarder = udpDataForwarder
	return f
}

func (f *ForwardingConfig) StartPortForwarding() error {
	var err error
	switch f.Protocol {
	case "tcp":
		err = f.startTCPPortForwarding()
	case "udp":
		err = f.startUDPPortForwarding()
	default:
		return errors.New("unsupported protocol: " + f.Protocol)
	}
	return err
}

func (f *ForwardingConfig) startTCPPortForwarding() error {
	localListener, err := net.Listen("tcp", f.SourceAddr)
	if err != nil {
		return fmt.Errorf("unable to bind to local TCP address: %s\n", err)
	}
	defer localListener.Close()

	log.Printf("TCP Port forwarding is active. Forwarding from %s to %s\n", f.SourceAddr, f.DestinationAddr)

	for {
		localConn, err := localListener.Accept()
		if err != nil {
			log.Errorf("Error accepting TCP connection: %s\n", err)
			continue
		}

		log.Printf("TCP connection established from %s\n", localConn.RemoteAddr())

		remoteConn, err := net.Dial("tcp", f.DestinationAddr)
		if err != nil {
			log.Warnf("Unable to connect to remote TCP address: %s\n", err)
			localConn.Close()
			log.Printf("TCP connection disconnected from %s\n", localConn.RemoteAddr())
			continue
		}
		go func() {
			f.TCPDataForwarder.Forward(localConn, remoteConn)
			localConn.Close()
			log.Printf("TCP connection disconnected from %s\n", localConn.RemoteAddr())
			remoteConn.Close()
		}()
	}
}

func (f *ForwardingConfig) startUDPPortForwarding() error {
	localUDPAddr, err := net.ResolveUDPAddr("udp", f.SourceAddr)
	if err != nil {
		return fmt.Errorf("error resolving local UDP address: %s\n", err)
	}
	localConn, err := net.ListenUDP("udp", localUDPAddr)
	if err != nil {
		return fmt.Errorf("unable to bind to local UDP address: %s\n", err)
	}
	defer localConn.Close()

	log.Printf("UDP Port forwarding is active. Forwarding from %s to %s\n", f.SourceAddr, f.DestinationAddr)

	remoteUDPAddr, err := net.ResolveUDPAddr("udp", f.DestinationAddr)
	if err != nil {
		return fmt.Errorf("error resolving remote UDP address: %s\n", err)
	}
	remoteConn, err := net.DialUDP("udp", nil, remoteUDPAddr)
	if err != nil {
		return fmt.Errorf("unable to connect to remote UDP address: %s\n", err)
	}
	defer remoteConn.Close()

	f.UDPDataForwarder.Forward(*localConn, *remoteConn)

	return nil
}
