package forwarder

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net"
)

type Forwarder struct {
	srcAddr  string
	dstAddr  string
	protocol string
	df       DataForwarder
}

func NewForwarder() *Forwarder {
	return &Forwarder{}
}

func (f *Forwarder) WithSourceAddr(srcAddr string) *Forwarder {
	f.srcAddr = srcAddr
	return f
}

func (f *Forwarder) WithDestinationAddr(dstAddr string) *Forwarder {
	f.dstAddr = dstAddr
	return f
}

func (f *Forwarder) WithProtocol(protocol string) *Forwarder {
	f.protocol = protocol
	return f
}

func (f *Forwarder) WithDataForwarder(df DataForwarder) *Forwarder {
	f.df = df
	return f
}

func (f *Forwarder) Start() error {
	var err error
	switch f.protocol {
	case "tcp":
		err = f.startTCP()
	case "udp":
		err = f.startUDP()
	default:
		return errors.New("unsupported protocol: " + f.protocol)
	}
	return err
}

func (f *Forwarder) startTCP() error {
	localListener, err := net.Listen("tcp", f.srcAddr)
	if err != nil {
		return fmt.Errorf("unable to bind to local TCP address: %s", err)
	}
	defer localListener.Close()

	log.Printf("TCP Port forwarding is active. Forwarding from %s to %s", f.srcAddr, f.dstAddr)

	for {
		localConn, err := localListener.Accept()
		if err != nil {
			log.Errorf("Error accepting TCP connection: %s", err)
			continue
		}

		log.Printf("TCP connection established from %s", localConn.RemoteAddr())

		remoteConn, err := net.Dial("tcp", f.dstAddr)
		if err != nil {
			log.Warnf("Unable to connect to remote TCP address: %s", err)
			localConn.Close()
			log.Printf("TCP connection disconnected from %s", localConn.RemoteAddr())
			continue
		}
		go func() {
			df := f.df.(*TCPDataForwarder).WithSrc(localConn).WithDst(remoteConn)
			df.Start()
			localConn.Close()
			log.Printf("TCP connection disconnected from %s", localConn.RemoteAddr())
			remoteConn.Close()
		}()
	}
}

func (f *Forwarder) startUDP() error {
	localUDPAddr, err := net.ResolveUDPAddr("udp", f.srcAddr)
	if err != nil {
		return fmt.Errorf("error resolving local UDP address: %s", err)
	}
	localConn, err := net.ListenUDP("udp", localUDPAddr)
	if err != nil {
		return fmt.Errorf("unable to bind to local UDP address: %s", err)
	}
	defer localConn.Close()

	log.Printf("UDP Port forwarding is active. Forwarding from %s to %s", f.srcAddr, f.dstAddr)

	remoteUDPAddr, err := net.ResolveUDPAddr("udp", f.dstAddr)
	if err != nil {
		return fmt.Errorf("error resolving remote UDP address: %s", err)
	}
	remoteConn, err := net.DialUDP("udp", nil, remoteUDPAddr)
	if err != nil {
		return fmt.Errorf("unable to connect to remote UDP address: %s", err)
	}
	defer remoteConn.Close()

	df := f.df.(*UDPDataForwarder).WithSrc(*localConn).WithDst(*remoteConn)
	df.Start()

	return nil
}
