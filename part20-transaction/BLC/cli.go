// Copyright 2019 The darwin Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package BLC

import (
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/boltdb/bolt"
)

// CLI作为客户端，命令行终端程序
type CLI struct {
	Blockchain *BlockChain // 区块链
}

// ValidateArgs 用来验证命令行参数个数
func (cli *CLI) ValidateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

// printUsage 打印参数信息
func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("\taddblock -data DATA -- 交易数据.")
	fmt.Println("\tprintchain -- 输出区块信息.")
}

// printChain 输出区块链中的所有区块的信息
func (cli *CLI) printChain() {
	var blockIterator *BlockChainIterator
	blockIterator = cli.Blockchain.Iterator()
	var hashInt big.Int
	// cmp := big.NewInt(1)
	// cmp.Lsh(x, n)
	for {

		err := blockIterator.DB.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(blockTableName))
			block := Deserialize(b.Get(blockIterator.CurrentHash))
			// fmt.Printf("Data: %s\n", string(block.Data))
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

func (cli *CLI) addBlock(data string) {
	// cli.Blockchain.AddBlock(tx)
	fmt.Println("新增了一个区块")
}

// Run 方法用来添加flag等相关的操作
func (cli *CLI) Run() {
	// 验证命令行参数
	cli.ValidateArgs()

	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printBlockchainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	addBlockData := addBlockCmd.String("data", "", "add Block data field")

	switch os.Args[1] {
	case "addblock":
		err := addBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
		// fmt.Println("xxx")
	default:
		cli.printUsage()
		os.Exit(1)
	}

	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			cli.printUsage()
			os.Exit(1)
		}
		// fmt.Println("Data: " + *addBlockData)
		cli.addBlock(*addBlockData)
	}

	if printBlockchainCmd.Parsed() {
		cli.printChain()
		//fmt.Println("xxxxxx")
	}
}
