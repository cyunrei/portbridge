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

func (c *ForwardingConfig) WithSourceAddr(sourceAddr string) *ForwardingConfig {
	c.SourceAddr = sourceAddr
	return c
}

func (c *ForwardingConfig) WithDestinationAddr(destinationAddr string) *ForwardingConfig {
	c.DestinationAddr = destinationAddr
	return c
}

func (c *ForwardingConfig) WithProtocol(protocol string) *ForwardingConfig {
	c.Protocol = protocol
	return c
}

func (c *ForwardingConfig) WithTCPDataForwarder(tcpDataForwarder TCPDataForwarder) *ForwardingConfig {
	c.TCPDataForwarder = tcpDataForwarder
	return c
}

func (c *ForwardingConfig) WithUDPDataForwarder(udpDataForwarder UDPDataForwarder) *ForwardingConfig {
	c.UDPDataForwarder = udpDataForwarder
	return c
}

func (c *ForwardingConfig) StartPortForwarding() error {
	var err error
	switch c.Protocol {
	case "tcp":
		err = startTCPPortForwarding(c.SourceAddr, c.DestinationAddr, c.TCPDataForwarder)
	case "udp":
		err = startUDPPortForwarding(c.SourceAddr, c.DestinationAddr, c.UDPDataForwarder)
	default:
		return errors.New("unsupported protocol: " + c.Protocol)
	}
	return err
}

func startTCPPortForwarding(sourceAddr, destinationAddr string, forwarder TCPDataForwarder) error {
	localListener, err := net.Listen("tcp", sourceAddr)
	if err != nil {
		return fmt.Errorf("unable to bind to local TCP address: %s\n", err)
	}
	defer localListener.Close()

	log.Printf("TCP Port forwarding is active. Forwarding from %s to %s\n", sourceAddr, destinationAddr)

	for {
		localConn, err := localListener.Accept()
		if err != nil {
			log.Errorf("Error accepting TCP connection: %s\n", err)
			continue
		}

		log.Printf("TCP connection established from %s\n", localConn.RemoteAddr())

		remoteConn, err := net.Dial("tcp", destinationAddr)
		if err != nil {
			log.Warnf("Unable to connect to remote TCP address: %s\n", err)
			localConn.Close()
			log.Printf("TCP connection disconnected from %s\n", localConn.RemoteAddr())
			continue
		}
		go func() {
			forwarder.Forward(localConn, remoteConn)
			localConn.Close()
			log.Printf("TCP connection disconnected from %s\n", localConn.RemoteAddr())
			remoteConn.Close()
		}()
	}
}

func startUDPPortForwarding(sourceAddr, destinationAddr string, forwarder UDPDataForwarder) error {
	localUDPAddr, err := net.ResolveUDPAddr("udp", sourceAddr)
	if err != nil {
		return fmt.Errorf("error resolving local UDP address: %s\n", err)
	}
	localConn, err := net.ListenUDP("udp", localUDPAddr)
	if err != nil {
		return fmt.Errorf("unable to bind to local UDP address: %s\n", err)
	}
	defer localConn.Close()

	log.Printf("UDP Port forwarding is active. Forwarding from %s to %s\n", sourceAddr, destinationAddr)

	remoteUDPAddr, err := net.ResolveUDPAddr("udp", destinationAddr)
	if err != nil {
		return fmt.Errorf("error resolving remote UDP address: %s\n", err)
	}
	remoteConn, err := net.DialUDP("udp", nil, remoteUDPAddr)
	if err != nil {
		return fmt.Errorf("unable to connect to remote UDP address: %s\n", err)
	}
	defer remoteConn.Close()

	forwarder.Forward(*localConn, *remoteConn)

	return nil
}
