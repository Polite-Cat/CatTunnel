package app

import (
	"fmt"
	"github.com/inhies/go-bytesize"
	"sync/atomic"
)

// GetReadBytes returns the number of bytes read
func (cat *Cat) GetReadBytes() uint64 {
	return atomic.LoadUint64(&cat.TotalReadBytes)
}

// GetWrittenBytes returns the number of bytes written
func (cat *Cat) GetWrittenBytes() uint64 {
	return atomic.LoadUint64(&cat.TotalWrittenBytes)
}

func (cat *Cat) PrintBytes(serverMode bool) string {
	if serverMode {
		return fmt.Sprintf("download %v upload %v", bytesize.New(float64(cat.GetWrittenBytes())).String(), bytesize.New(float64(cat.GetReadBytes())).String())
	}
	return fmt.Sprintf("download %v upload %v", bytesize.New(float64(cat.GetReadBytes())).String(), bytesize.New(float64(cat.GetWrittenBytes())).String())
}
