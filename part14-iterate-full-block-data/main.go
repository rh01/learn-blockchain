package main

import (
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/boltdb/bolt"
	"github.com/rh01/learn-blockchain/part14-iterate-full-block-data/BLC"
)

const (
	blockTableName = "blocks"
)

func main() {
	blockchain := BLC.NewBlockChain()

	// blockchain.AddBlock("XiaoMing Send to Xiaohua 20 BTC")
	// blockchain.AddBlock("XiaoMing Send to Xiaohua 20 BTC")
	// blockchain.AddBlock("XiaoMing Send to Xiaohua 20 BTC")
	// blockchain.AddBlock("XiaoMing Send to Xiaohua 20 BTC")
	// fmt.Println(blockchain)

	var blockIterator *BLC.BlockChainIterator
	blockIterator = blockchain.Iterator()
	var hashInt big.Int
	// cmp := big.NewInt(1)
	// cmp.Lsh(x, n)
	for {

		err := blockIterator.DB.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(blockTableName))
			block := BLC.Deserialize(b.Get(blockIterator.CurrentHash))
			fmt.Printf("Data: %s\n", string(block.Data))
			fmt.Printf("PrevBlockHash: %x\n", block.PrevBlockHash)
			fmt.Printf("TimeStamp: %v\n", time.Unix(block.TimeStamp, 0).Format("2006-01-02 03:04:05 PM")) // Format 参数不能随意修改
			fmt.Printf("Hash: %x\n", (block.Hash))

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
}
