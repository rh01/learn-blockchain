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
	// fmt.Println("\taddblock -data DATA -- 交易数据.")
	fmt.Println("\tgetbalance -address ADDRESS -- Get balance of ADDRESS.")
	fmt.Println("\tcreateblockchain -address ADDRESS  -- Create a blockchain and send genesis block reward to ADDRESS .")
	fmt.Println("\tprintchain -- 输出区块信息.")
	fmt.Println("\tsendd -from FROM -to TO -amount AMOUNT - Send amount of coin from FROM to TO.")
}

// printChain 输出区块链中的所有区块的信息
func (cli *CLI) printChain() {
	// 判断数据库是否存在
	if dbExists() == false {
		cli.printUsage()
		return
	}

	blockchain := GetBlockchain()
	defer blockchain.Db.Close()

	var blockIterator *BlockChainIterator
	blockIterator = blockchain.Iterator()
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
			for _, tx := range block.Transactions {
				fmt.Println(tx)
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
}

func (cli *CLI) sendToken() {
	// 1.10 -> liyuechun
	// 2.3 -> 转给haolin

	// 新建一个交易
	//tx := Transaction{
	//	ID:   nil,
	//	Vin:  nil,
	//	Vout: nil,
	//}

	var inputs []*Transaction

	blockchain := GetBlockchain()
	defer blockchain.Db.Close()

	// 新建交易
	tx1 := NewUTXOTransaction("shh", "xiaoqiang", 2, blockchain, inputs)

	fmt.Println("第一笔交易")
	fmt.Println(tx1)

	inputs = append(inputs, tx1)
	tx2 := NewUTXOTransaction("shh", "xiaoming", 3, blockchain, inputs)

	fmt.Println("第二笔交易")
	fmt.Println(tx2)

	inputs = append(inputs, tx2)
	tx3 := NewUTXOTransaction("shh", "daa", 3, blockchain, inputs)

	fmt.Println("第三笔交易")
	fmt.Println(tx3)
	// 挖矿
	// cli.Blockchain.MineBlock([]*Transaction{tx1, tx2, tx3})

}

func (cli *CLI) addBlock(data string) {
	cli.sendToken()
}

// Run 方法用来添加flag等相关的操作
func (cli *CLI) Run() {
	// 验证命令行参数,如果没有参数，将会打印Usage
	cli.ValidateArgs()

	printBlockchainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	addBlockData := addBlockCmd.String("data", "", "add Block data field")

	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	genenisAddress := createBlockchainCmd.String("address", "", "创建创世区块，并将区块打包到数据库中")

	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	balanceAddress := getBalanceCmd.String("address", "", "查询余额")

	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	fromAddress := sendCmd.String("from", "", "源地址...")
	toAddress := sendCmd.String("to", "", "目标地址....")
	amount := sendCmd.String("amount", "", "转账的数量....")

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
	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
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

	if createBlockchainCmd.Parsed() {
		if *genenisAddress == "" {
			cli.printUsage()
			os.Exit(1)
		}

		cli.createBlockchain(*genenisAddress)
	}

	if getBalanceCmd.Parsed() {
		if *balanceAddress == "" {
			cli.printUsage()
			os.Exit(1)
		}
		fmt.Printf("查询 %s 的余额...\n", *balanceAddress)
	}

	if sendCmd.Parsed() {
		if *fromAddress == "" || *toAddress == "" || *amount == "" {
			cli.printUsage()
			os.Exit(1)
		}

		fmt.Printf("send -from %s -to %s -amount %s\n", *fromAddress, *toAddress, *amount)

		sendFrom := JSONToArray(*fromAddress)
		sendTo := JSONToArray(*toAddress)
		sendAmount := JSONToArray(*amount)

		fmt.Println("from:")
		fmt.Println(sendFrom)

		fmt.Println("to:")
		fmt.Println(sendTo)

		fmt.Println("amount:")
		fmt.Println(sendAmount)

	}

	if printBlockchainCmd.Parsed() {
		cli.printChain()
		//fmt.Println("xxxxxx")
	}
}

func (cli *CLI) createBlockchain(address string) {
	if dbExists() {
		fmt.Println("创世区块已经存在....")
		os.Exit(1)
	}

	CreateBlockchain(address)
}
