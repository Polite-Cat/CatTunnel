package tcp

import (
	"github.com/networm6/PoliteCat/daemon/service"
	"github.com/networm6/PoliteCat/protocol"
	"io"
	"log"
	"net"
	"time"

	"github.com/networm6/PoliteCat/common/cipher"
	"github.com/songgao/water"
)

// StartClient Start tcp client
func StartClient(config service.Config) {
	log.Printf("vtun tcp client started on %v", config.LocalAddr)
	tunInterface := protocol.CreateTun(config)
	go tunToTcp(config, tunInterface)
	for {
		conn, err := net.DialTimeout("tcp", config.ServerAddr, time.Duration(config.Timeout)*time.Second)
		if err != nil {
			time.Sleep(3 * time.Second)
			continue
		}
		getCache().Set(Key, conn, 24*time.Hour)
		tcpToTun(config, conn, tunInterface)
		getCache().Delete(Key)
	}
}

func tunToTcp(config service.Config, tunInterface *water.Interface) {
	packet := make([]byte, config.MTU)
	for {
		num, err := tunInterface.Read(packet)
		if err != nil || num == 0 {
			continue
		}
		if conn, ok := getCache().Get(Key); ok {
			bytes := packet[:num]
			if config.Enc {
				packet = cipher.XOR(packet)
			}
			conn := conn.(net.Conn)
			_ = conn.SetWriteDeadline(time.Now().Add(time.Duration(config.Timeout) * time.Second))
			_, _ = conn.Write(bytes)
		}
	}
}

func tcpToTun(config service.Config, tcpConn net.Conn, tunInterface *water.Interface) {
	defer func(tcpConn net.Conn) {
		_ = tcpConn.Close()
	}(tcpConn)
	packet := make([]byte, config.MTU)
	for {
		_ = tcpConn.SetReadDeadline(time.Now().Add(time.Duration(config.Timeout) * time.Second))
		num, err := tcpConn.Read(packet)
		if err != nil || err == io.EOF {
			break
		}
		bytes := packet[:num]
		if config.Enc {
			bytes = cipher.XOR(bytes)
		}
		_, err = tunInterface.Write(bytes)
		if err != nil {
			break
		}
	}
}
