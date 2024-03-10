package forward

import (
	log "github.com/sirupsen/logrus"
	"net"
)

const DefaultUDPBufferSize uint64 = 1024

type SimpleUDPDataForwarder struct {
	BufferSize uint64
}

func NewSimpleUDPDataForwarder() *SimpleUDPDataForwarder {
	return &SimpleUDPDataForwarder{
		BufferSize: DefaultUDPBufferSize,
	}
}

func (f *SimpleUDPDataForwarder) SetBufferSize(size uint64) *SimpleUDPDataForwarder {
	f.BufferSize = size
	return f
}

func (f *SimpleUDPDataForwarder) Forward(conn, remoteConn net.UDPConn) {
	connBuffer := make([]byte, f.BufferSize)

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

		remoteBuffer := make([]byte, f.BufferSize)
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
