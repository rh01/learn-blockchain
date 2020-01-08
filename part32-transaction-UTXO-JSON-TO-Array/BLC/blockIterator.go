// Copyright 2019 The darwin Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package BLC

import (
	"log"

	"github.com/boltdb/bolt"
)

// BlockChainIterator 迭代器
type BlockChainIterator struct {
	CurrentHash []byte // 当前正在遍历区块的哈希
	DB          *bolt.DB
}

// Next
func (bci *BlockChainIterator) Next() *BlockChainIterator {
	// 具体的步骤如下：
	// 首先打开数据库
	var currentBlock *Block
	var nextHash []byte
	err := bci.DB.View(func(tx *bolt.Tx) error {
		// 获取当前表
		b := tx.Bucket([]byte(blockTableName))

		// 获取最新区块的序列化后的表示
		currentBlockBytes := b.Get(bci.CurrentHash)

		// 使用反序列化方法得到Block结构的区块
		currentBlock = Deserialize(currentBlockBytes)

		nextHash = currentBlock.PrevBlockHash
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return &BlockChainIterator{nextHash, bci.DB}
}
