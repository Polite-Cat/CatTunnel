package encrypt

var _key = []byte("fuck_key")

// SetMixinKey 设置密钥。
func SetMixinKey(key string) {
	_key = []byte(key)
}

// Xor 混淆你的数据。
func Xor(src []byte) []byte {
	var (
		_lenK = len(_key)
		_lenV = len(src)
	)

	for i := 0; i < _lenV; i++ {
		src[i] ^= _key[i%_lenK]
	}
	return src
}

// None 不加密
func None(src []byte) []byte {
	return src
}

// Reverse 反转
func Reverse(src []byte) []byte {
	var (
		_lenV = len(src)
	)
	for i, j := 0, _lenV-1; i < j; i, j = i+1, j-1 {
		src[i], src[j] = src[j], src[i]
	}
	return src
}

func Copy(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
