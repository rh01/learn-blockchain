// Copyright 2019 The darwin Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package BLC

import (
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/boltdb/bolt"
)

const (
	dbName              = "blockchain.db" // 数据库的名字
	blockTableName      = "blocks"        // 表的名字
	genesisCoinbaseData = "xxxx"          // 创世区块的交易信息
	subsidy             = 10
)

// BlockChain 存储有序的区块
type BlockChain struct {
	Tip []byte   // 区块链中最后一个区块的哈希
	Db  *bolt.DB // 数据库
}

// 先找到包含当前用户未花费输出的所有交易的集合
// 返回交易的数组
func (blockchain *BlockChain) FindUnspentTranscation(address string) []Transaction {
	blockIterator := blockchain.Iterator()
	var hashInt big.Int

	for {

		err := blockIterator.DB.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(blockTableName))
			block := Deserialize(b.Get(blockIterator.CurrentHash))
			// fmt.Printf("Data: %s\n", string(block.Data))
			fmt.Printf("PrevBlockHash: %x\n", block.PrevBlockHash)
			fmt.Printf("TimeStamp: %v\n", time.Unix(block.TimeStamp, 0).Format("2006-01-02 03:04:05 PM")) // Format 参数不能随意修改
			fmt.Printf("Hash: %x\n", block.Hash)

			for _, transaction := range block.Transactions {
				fmt.Printf("Transaction ID: %x\n", transaction.ID)
			}

			fmt.Println("")
			return nil
		})
		if err != nil {
			log.Panic(err)
		}

		// fmt.Printf("%x\n", blockIterator.CurrentHash)

		blockIterator = blockIterator.Next()

		hashInt.SetBytes(blockIterator.CurrentHash)

		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break
		}
	}

	return nil
}

// Iterator 返回 BlockChainIterator 对象
func (blockchain *BlockChain) Iterator() *BlockChainIterator {

	return &BlockChainIterator{
		blockchain.Tip,
		blockchain.Db,
	}
}

// AddBlock 新增区块
//func (blockchain *BlockChain) AddBlock(transactions []*Transaction) {
//	// 1. 创建区块
//	newBlock := NewBlock(transactions, blockchain.Tip)
//	// 获取数据库表
//	err := blockchain.Db.Update(func(tx *bolt.Tx) error {
//		b := tx.Bucket([]byte(blockTableName))
//		// 存储新区块的数据
//		err := b.Put(newBlock.Hash, newBlock.Serialize())
//		if err != nil {
//			log.Panic(err)
//		}
//
//		// 更新l对应的Hash
//		err = b.Put([]byte("l"), newBlock.Hash)
//		if err != nil {
//			log.Panic(err)
//		}
//		// 更新区块链的Tip为最新的区块的Hash值
//		blockchain.Tip = newBlock.Hash
//		return nil
//	})
//	if err != nil {
//		log.Panic(err)
//	}
//}

// dbExists 判断数据库是否存在
func dbExists() bool {
	if _, err := os.Stat(dbName); err != nil {
		return false
	}

	return true
}

// CreateBlockchainWithGenesisBlock 创建一个带有创世区块节点的区块链
func CreateBlockchainWithGenesisBlock() *BlockChain {

	if dbExists() {

		fmt.Println("创世区块已经存在...")

		db, err := bolt.Open(dbName, 0600, nil)
		if err != nil {
			log.Fatal(err)
		}

		// 这一步是为了避免覆盖之前存储的区块信息，并重新开始后
		// 从最后一个区块开始计算
		var blockchain *BlockChain
		err = db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(blockTableName))
			hash := b.Get([]byte("l"))
			blockchain = &BlockChain{hash, db}

			return nil
		})
		// (2) 创建创世区块
		// 如果err为nil，则说明这个表里没有区块信息
		// 此时需要创建创世区块
		if err != nil {
			log.Panic(err)
		}

		return blockchain
	}

	// 1尝试打开数据库或者创建数据库
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	var blockHash []byte
	// 2调用db.Updata更新数据库
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b == nil {
			// (1) 判断表不存在则创建表
			b, err = tx.CreateBucket([]byte(blockTableName))
			if err != nil {
				log.Panic(err)
			}
		}

		if b != nil {
			// 创建创世区块
			// 创建创世区块的交易对象
			genesisCoinBase := NewCoinbaseTx("shh", genesisCoinbaseData)
			genesisBlock := NewGenesisBlock(genesisCoinBase)
			// (3) 将创世区块序列化
			// 把创世区块的Hash值作为key，Block的序列化数据
			// 作为value存储在表中，存储在数据库
			err := b.Put(genesisBlock.Hash, genesisBlock.Serialize())
			if err != nil {
				log.Panic(err)
			}

			// 存储最新的区块
			// (5) 设置一个key，l，将hash作为value再次存储在数据里面
			err = b.Put([]byte("l"), genesisBlock.Hash)
			if err != nil {
				log.Panic(err)
			}

			blockHash = genesisBlock.Hash
		}

		return nil

	})

	return &BlockChain{blockHash, db}
}

// NewBlockChain 初始化一个区块链
func NewBlockChain() *BlockChain {
	var tip []byte // 获取最后一个区块的hash值

	// 尝试打开或者创建数据库
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b == nil {
			b, err = tx.CreateBucket([]byte(blockTableName))
			if err != nil {
				log.Panic(err)
			}
			// 将创世区块序列化后的数据存储在表中
			genesisBlock := NewGenesisBlock(NewCoinbaseTx("sss", genesisCoinbaseData))

			err = b.Put(genesisBlock.Hash, genesisBlock.Serialize())
			if err != nil {
				log.Panic(err)
			}

			// 接下来存储Hash
			err = b.Put([]byte("l"), genesisBlock.Hash)
			if err != nil {
				log.Panic(err)
			}

			tip = genesisBlock.Hash
		} else {
			tip = b.Get([]byte("l"))
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return &BlockChain{tip, db}
}
