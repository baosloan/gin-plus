package hash

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5(text string) string {
	//创建一个md5哈希对象
	hash := md5.New()
	//将需要加密的字符串写入哈希对象
	hash.Write([]byte(text))
	//计算并返回哈希值，结果是一个字节切片，将字节切片转换为十六进制字符串表示
	return hex.EncodeToString(hash.Sum(nil))
}
