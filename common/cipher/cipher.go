package cipher

import "fmt"

var _key = []byte("c2ltb24yMDIzLXZ0dW4=")

func GenerateKey(key string) {
	_key = []byte(key)
	_key2 := []byte(fmt.Sprintf("%d", len(key)))
	_key = XOR(_key2)
}

func XOR(src []byte) []byte {
	_kLen := len(_key)
	for i := 0; i < len(src); i++ {
		src[i] ^= _key[i%_kLen]
	}
	return src
}
