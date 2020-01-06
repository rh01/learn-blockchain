package BLC

import (
	"bytes"
	"encoding/binary"
	"log"
)

// 将int64转为字节数组
func IntToHex(target int64) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, target)
	if err != nil {
		log.Panic(err)
	}

	return buf.Bytes()
}
