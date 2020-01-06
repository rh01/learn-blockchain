// Copyright 2019 The darwin Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package BLC

import (
	"encoding/hex"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"math/big"
	"os"
)

const (
	dbName              = "blockchain.db" // 数据库的名字
	blockTableName      = "blocks"        // 表的名字
	genesisCoinbaseData = "2009-01-03 xxxxx"          // 创世区块的交易信息
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
	// 存储未花费输出的交易
	var unspentTXs []Transaction
	// 存储区块当中所有的交易输入
	spentTXOs := make(map[string][]int)

	blockIterator := blockchain.Iterator()
	var hashInt big.Int

	for {

		err := blockIterator.DB.View(func(tx *bolt.Tx) error {
			// 获取表
			b := tx.Bucket([]byte(blockTableName))
			//反序列化
			block := Deserialize(b.Get(blockIterator.CurrentHash))

			for _, transaction := range block.Transactions {
				// fmt.Printf("Transaction ID: %x\n", transaction.ID)
				//将byte array 类型转为string类型
				txID := hex.EncodeToString(transaction.ID)

			Outputs:
				for outIdx, out := range transaction.Vout {
					//	判断是否被花费
					if spentTXOs[txID] != nil {
						for _, spentOut := range spentTXOs[txID] {
							// 相等说明，当前的输出在这个tx中已被花费
							if spentOut == outIdx {
								continue Outputs
							}
						}
					}
					// 若未花费，添加到未花费的
					if out.CanBeUnlockedWith(address) {
						unspentTXs = append(unspentTXs, *transaction)
					}
				}

				//
				if transaction.isCoinbase() == false {
					for _, in := range transaction.Vin {
						if in.CanUnlockOutputWith(address) {
							inTxID := hex.EncodeToString(in.Txid)
							spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Vout)
						}
					}
				}
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

	return unspentTXs
}

// 通过transaction数组查找可用的未花费的输出信息
//  16 10
// 查找可用的交易信息
func (blockchain *BlockChain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	// 字典，存储交易id，Vout中的未花费TXOutput的index
	unspentOutputs := make(map[string][]int)

	// 查看未花费
	unspentTXs := blockchain.FindUnspentTranscation(address)

	accumulated := 0 // 统计未花费对应TXoutputs的总量

Work:
	// 遍历交易数组
	for _, tx := range unspentTXs {
		txID := hex.EncodeToString(tx.ID)
		// 遍历交易里面的Vout
		for outIdx, output := range tx.Vout {
			if output.CanBeUnlockedWith(address) && accumulated < amount {
				accumulated += output.Value
				unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)

				if accumulated >= amount {
					break Work
				}
			}
		}
	}
	// 12, {"1111":[1,2,3]}
	return accumulated, unspentOutputs
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
			genesisBlock := NewGenesisBlock(NewCoinbaseTx("shh", genesisCoinbaseData))

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

// MineBlock 挖矿
// 根据交易的数组，打包新的区块
func (blockchain *BlockChain) MineBlock(txs []*Transaction)  {
	err := blockchain.Db.Update(func(tx *bolt.Tx) error {
		// 新建区块
		// 拿到新的区块
		newBlock := NewBlock(txs, blockchain.Tip)

		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			// key: newBlock Hash
			// value: newBlock.Serialized
			err := b.Put(newBlock.Hash, newBlock.Serialize())
			if err != nil {
				log.Panic(err)
			}
			// key: l
			// value: newBlockHash
			err = b.Put([]byte("l"), newBlock.Hash)
			if err != nil {
				log.Panic(err)
			}
			// 更新区块的最新的hash
			blockchain.Tip = newBlock.Hash
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	// 将区块存储到区块连中
}