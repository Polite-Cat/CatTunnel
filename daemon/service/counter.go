package service

import (
	"fmt"
	"io/ioutil"
	"sync/atomic"
)

var TotalReadByte uint64 = 0

var TotalWriteByte uint64 = 0

func IncrReadByte(n int) {
	atomic.AddUint64(&TotalReadByte, uint64(n))
}

func IncrWriteByte(n int) {
	atomic.AddUint64(&TotalWriteByte, uint64(n))
}
func init() {
	//go func() {
	//	for {
	//		readByte := atomic.LoadUint64(&TotalReadByte)
	//		writeByte := atomic.LoadUint64(&TotalWriteByte)
	//
	//		// 将数据写入文件
	//		err := writeDataToFile(readByte, writeByte)
	//		if err != nil {
	//			fmt.Println("无法写入数据文件：", err)
	//		}
	//
	//		// 等待 5 秒
	//		time.Sleep(5 * time.Second)
	//	}
	//}()
}

func writeDataToFile(readByte, writeByte uint64) error {
	data := fmt.Sprintf("TotalReadByte: %d\nTotalWriteByte: %d\n", readByte, writeByte)
	filePath := "/var/run/vtun"

	return ioutil.WriteFile(filePath, []byte(data), 0644)
}
