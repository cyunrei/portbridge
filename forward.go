package main

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
)

type ForwardingConfig struct {
	SourceAddr       string
	DestinationAddr  string
	Protocol         string
	TCPDataForwarder TCPDataForwarder
}

func startPortForwarding(c ForwardingConfig) error {
	var err error
	switch c.Protocol {
	case "tcp":
		err = startTCPPortForwarding(c.SourceAddr, c.DestinationAddr, c.TCPDataForwarder)
	case "udp":
		err = startUDPPortForwarding(c.SourceAddr, c.DestinationAddr)
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

		log.Printf("New TCP connection established from %s\n", localConn.RemoteAddr())

		remoteConn, err := net.Dial("tcp", destinationAddr)
		if err != nil {
			log.Warnf("Unable to connect to TCP source address: %s\n", err)
			localConn.Close()
			continue
		}

		go forwarder.Forward(localConn, remoteConn)

	}
}

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

func startUDPPortForwarding(sourceAddr, destinationAddr string) error {
	localUDPAddr, err := net.ResolveUDPAddr("udp", sourceAddr)
	if err != nil {
		return fmt.Errorf("error resolving local address: %s\n", err)
	}
	remoteUDPAddr, err := net.ResolveUDPAddr("udp", destinationAddr)
	if err != nil {
		return fmt.Errorf("error resolving remote address: %s\n", err)
	}

	conn, err := net.ListenUDP("udp", localUDPAddr)
	if err != nil {
		return fmt.Errorf("unable to bind to local UDP address: %s\n", err)
	}
	defer conn.Close()

	log.Printf("UDP Port forwarding is active. Forwarding from %s to %s\n", sourceAddr, destinationAddr)

	forwardUDPData(conn, remoteUDPAddr)

	return nil
}

func forwardUDPData(conn *net.UDPConn, remoteAddr *net.UDPAddr) {
	bufferSize := 1024

	remoteConn, err := net.DialUDP("udp", nil, remoteAddr)
	if err != nil {
		log.Fatalf("Error establishing remote connection: %s\n", err)
	}
	defer remoteConn.Close()

	connBuffer := make([]byte, bufferSize)

	type dataWithAddr struct {
		data []byte
		addr *net.UDPAddr
	}

	dataChannel := make(chan dataWithAddr)

	readDataFromConn := func() {
		for {
			n, connAddr, err := conn.ReadFromUDP(connBuffer)
			if err != nil {
				log.Errorf("Error reading from UDP: %s\n", err)
				continue
			}
			dataChannel <- dataWithAddr{connBuffer[:n], connAddr}
		}
	}

	go readDataFromConn()

	writeDataToRemote := func(data []byte, addr *net.UDPAddr) {
		_, err := remoteConn.Write(data)
		if err != nil {
			log.Errorf("Error writing to remote UDP: %s\n", err)
			return
		}

		remoteBuffer := make([]byte, bufferSize)
		m, err := remoteConn.Read(remoteBuffer)
		if err != nil {
			log.Errorf("Error reading from remote UDP: %s\n", err)
			return
		}

		_, err = conn.WriteToUDP(remoteBuffer[:m], addr)
		if err != nil {
			log.Errorf("Error writing to UDP: %s\n", err)
		}
	}

	for {
		dataWithAddr := <-dataChannel
		go writeDataToRemote(dataWithAddr.data, dataWithAddr.addr)
	}
}
