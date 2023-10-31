package tunnel

import (
	"github.com/networm6/CatTunnel/common/encrypt"
	"github.com/networm6/CatTunnel/common/tools"
	"sync/atomic"
)

func (tun *Tunnel) readFromTunnel() {
	packet := make([]byte, tun._bufferSize)
	for tools.ContextOpened(*tun.LifeCtx) {
		num, err := tun._tunInterface.Read(packet)
		tun.incrWrittenBytes(num)
		if err != nil {
			continue
		}
		mixinPacket := tun._mixinFunc(packet[:num])
		tun.OutputStream <- encrypt.Copy(mixinPacket)
	}
}

func (tun *Tunnel) writeToTunnel() {
	for tools.ContextOpened(*tun.LifeCtx) {
		//startTime := time.Now().UnixNano()
		packet := <-tun.InputStream
		mixinPacket := tun._mixinFunc(packet)
		num, err := tun._tunInterface.Write(mixinPacket)
		if err != nil {
			continue
		}
		tun.incrReadBytes(num)
		//endTime := time.Now().UnixNano()
		//nanoSeconds := float64(endTime - startTime) // ns
		//fmt.Printf("writeToTunnel %d cost time(ms) :%v;\n", num, nanoSeconds/1e6)
	}
}

func (tun *Tunnel) incrReadBytes(n int) {
	atomic.AddUint64(tun._totalReadBytes, uint64(n))
}

func (tun *Tunnel) incrWrittenBytes(n int) {
	atomic.AddUint64(tun._totalWrittenBytes, uint64(n))
}
