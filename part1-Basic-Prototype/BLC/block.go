package BLC

import (
	"bytes"
	"crypto/sha256"
	"strconv"
	"time"
)

// Block is BLC basic structure
type Block struct {
	TimeStamp     int64  // 时间戳
	Data          []byte // 交易数据
	PrevBlockHash []byte // 上一区块hash
	Hash          []byte // 当前区块hash
}

func (block *Block) SetHash() {
	// 1. 将时间戳转化为字节数组
	// （1）将int64转字符串
	// 第二个参数的范围 2~36，代表进制
	// （2）将时间戳转字节数据
	timestamp := []byte(strconv.FormatInt(block.TimeStamp, 10))
	// 2. 将除了Hash以外的其他属性，一字节数组的形式全品接起来
	headers := bytes.Join([][]byte{block.PrevBlockHash, block.Data, timestamp}, []byte{})
	// 3. 将拼接起来数据进行256hash
	hash := sha256.Sum256(headers)
	// 4. 将hash付给Hash属性字节
	block.Hash = hash[:]
}

// NewBlock is factory method to make a Block ref
func NewBlock(data string, prevBlockHash []byte) *Block {

	// 创建区块
	block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}}
	// 设置当前区块hash
	block.SetHash()

	return block
}

// 创建创世区块
func NewGenesisBlock() *Block {
	return NewBlock("Genenius Block", []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
}
