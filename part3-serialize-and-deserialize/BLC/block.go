package BLC

import (
	"bytes"
	"encoding/gob"
	"log"
)

type Block struct {
	Data  []byte // 交易数据
	Nonce int    // 挖矿时满足Nonce条件的随机值
}

// 将区块序列化成字节数组
func (block *Block) Serialize() []byte {

	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(block)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

// 将区块序列化成字节数组
func Deserialize(blockbytes []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(blockbytes))
	err := decoder.Decode(&block)
	if err != nil {
		log.Fatal(err)
	}
	return &block

}
