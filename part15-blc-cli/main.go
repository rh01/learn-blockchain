package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/rh01/blockchain/part15-blc-cli/BLC"
)

const (
	blockTableName = "blocks"
)

func main() {

	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printchainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	addBlockData := addBlockCmd.String("data", "", "New Block data")
	switch os.Args[1] {
	case "addblock":
		err := addBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		fmt.Println("No addblock and printchain cmd")
		// 下面就是处理的逻辑，主要是检查数据库文件为空，若为空挖矿
		_, err := os.Stat("blockchain.db")
		if err != nil {
			blc := BLC.NewBlockChain()
			blc.
		}

		os.Exit(1)
	}
	// blockchain := BLC.NewBlockChain()

	// blockchain.AddBlock("XiaoMing Send to Xiaohua 20 BTC")
	// blockchain.AddBlock("XiaoMing Send to Xiaohua 20 BTC")
	// blockchain.AddBlock("XiaoMing Send to Xiaohua 20 BTC")
	// blockchain.AddBlock("XiaoMing Send to Xiaohua 20 BTC")
	// fmt.Println(blockchain)

	// var blockIterator *BLC.BlockChainIterator
	// blockIterator = blockchain.Iterator()
	// var hashInt big.Int
	// // cmp := big.NewInt(1)
	// // cmp.Lsh(x, n)
	// for {

	// 	err := blockIterator.DB.View(func(tx *bolt.Tx) error {
	// 		b := tx.Bucket([]byte(blockTableName))
	// 		block := BLC.Deserialize(b.Get(blockIterator.CurrentHash))
	// 		fmt.Printf("Data: %s\n", string(block.Data))
	// 		fmt.Printf("PrevBlockHash: %x\n", block.PrevBlockHash)
	// 		fmt.Printf("TimeStamp: %v\n", time.Unix(block.TimeStamp, 0).Format("2006-01-02 03:04:05 PM")) // Format 参数不能随意修改
	// 		fmt.Printf("Hash: %x\n", (block.Hash))

	// 		fmt.Println("")
	// 		return nil
	// 	})
	// 	if err != nil {
	// 		log.Panic(err)
	// 	}

	// 	// fmt.Printf("%x\n", blockIterator.CurrentHash)

	// 	blockIterator = blockIterator.Next()

	// 	hashInt.SetBytes(blockIterator.CurrentHash)

	// 	if hashInt.Cmp(big.NewInt(0)) == 0 {
	// 		break
	// 	}
	// }
}
