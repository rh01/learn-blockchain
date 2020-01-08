// Copyright 2019 The darwin Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package BLC

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"strconv"
	"time"
)

// Block is BLC basic structure
type Block struct {
	TimeStamp     int64          // 时间戳
	Transactions  []*Transaction // 每一个区块包含很多交易数据
	PrevBlockHash []byte         // 上一区块hash
	Hash          []byte         // 当前区块hash
	Nonce         int            // 挖矿时满足Nonce条件的随机值
}

// NewBlock is factory method to make a Block ref
func NewBlock(transactions []*Transaction, prevBlockHash []byte) *Block {

	// 创建区块
	block := &Block{time.Now().Unix(), transactions, prevBlockHash, []byte{}, 0}

	// 将block作为参数，创一pow对象
	pow := NewProofOfWork(block)

	// Run() 执行一次工作量证明
	hash, nonce := pow.Run()

	// 设置区块的Hash
	block.Hash = hash[:]

	// 设置Nonce值
	block.Nonce = nonce

	// 修改为工作能力证明，该部分的区块生成应该是通过工作力证明得到的

	return block
}

// HashTransactions 将transaction序列化为字节数组
// 将区块里面所有交易的ID拼接在一起，并且生成Hash
func (block *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range block.Transactions {
		txHashes = append(txHashes, tx.ID)
	}

	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))
	return txHash[:]
}

// SetHash 设置当前区块的哈希
func (block *Block) SetHash() {
	// 1. 将时间戳转化为字节数组
	// （1）将int64转字符串
	// 第二个参数的范围 2~36，代表进制
	// （2）将时间戳转字节数据
	timestamp := []byte(strconv.FormatInt(block.TimeStamp, 10))

	// TODO: 这一步是我的理解，后续会改动
	var txHashes [][]byte
	for _, tx := range block.Transactions {
		txHashes = append(txHashes, tx.ID)
	}
	// 2. 将除了Hash以外的其他属性，一字节数组的形式全品接起来
	headers := bytes.Join([][]byte{block.PrevBlockHash, bytes.Join(txHashes, []byte{}), timestamp}, []byte{})
	// 3. 将拼接起来数据进行256hash
	hash := sha256.Sum256(headers)
	// 4. 将hash付给Hash属性字节
	block.Hash = hash[:]
}

// NewGenesisBlock 创建创世区块
func NewGenesisBlock(coinbase *Transaction) *Block {
	return NewBlock([]*Transaction{coinbase}, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
}

// Serialize 将区块序列化成字节数组
func (block *Block) Serialize() []byte {

	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(block)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

// Deserialize 将区块序列化成字节数组
func Deserialize(blockbytes []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(blockbytes))
	err := decoder.Decode(&block)
	if err != nil {
		log.Fatal(err)
	}
	return &block

}
