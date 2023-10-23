package tcp

import (
	"github.com/networm6/PoliteCat/common/std"
	"github.com/networm6/PoliteCat/daemon/service"
	"github.com/networm6/PoliteCat/protocol"
	"io"
	"log"
	"net"
	"time"

	"github.com/networm6/PoliteCat/common/cipher"
	"github.com/songgao/water"
)

// StartServer Start tcp server
func StartServer(config service.Config) {
	log.Printf("vtun tcp server started on %v", config.LocalAddr)
	tunInterface := protocol.CreateTun(config)
	// server -> client
	go toClient(config, tunInterface)
	ln, err := net.Listen("tcp", config.LocalAddr)
	if err != nil {
		log.Println(err)
		return
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		// client -> server
		go toServer(config, conn, tunInterface)
	}
}

func toClient(config service.Config, tunInterface *water.Interface) {
	packet := make([]byte, config.MTU)
	for {
		num, err := tunInterface.Read(packet)
		if err != nil || err == io.EOF || num == 0 {
			continue
		}
		bytes := packet[:num]
		if key := std.GetDestinationKey(bytes); key != "" {
			if conn, ok := getCache().Get(key); ok {
				if config.Enc {
					bytes = cipher.XOR(bytes)
				}
				service.IncrWriteByte(num)
				_, _ = conn.(net.Conn).Write(bytes)
			}
		}
	}
}

func toServer(config service.Config, tcpConn net.Conn, tunInterface *water.Interface) {
	defer func(tcpConn net.Conn) {
		_ = tcpConn.Close()
	}(tcpConn)
	packet := make([]byte, config.MTU)
	for {
		_ = tcpConn.SetReadDeadline(time.Now().Add(time.Duration(config.Timeout) * time.Second))
		n, err := tcpConn.Read(packet)
		if err != nil || err == io.EOF {
			break
		}
		b := packet[:n]
		if config.Enc {
			b = cipher.XOR(b)
		}
		if key := std.GetSourceKey(b); key != "" {
			getCache().Set(key, tcpConn, 10*time.Minute)
			service.IncrReadByte(len(b))
			_, _ = tunInterface.Write(b)
		}
	}
}
