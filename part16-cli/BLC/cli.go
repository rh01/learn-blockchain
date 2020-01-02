// Copyright 2019 The darwin Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package BLC

import (
	"flag"
	"fmt"
	"log"
	"os"
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

// Run 方法用来添加flag等相关的操作
func (cli *CLI) Run() {
	// 验证命令行参数
	cli.ValidateArgs()

	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printBlockchainCmd := flag.NewFlagSet("printblockchain", flag.ExitOnError)

	addBlockData := addBlockCmd.String("data", "", "add Block data field")

	switch os.Args[1] {
	case "addblock":
		err := addBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printblockchain":
		err := printBlockchainCmd.Parse(os.Args[2:])
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
		fmt.Println("Data: " + *addBlockData)
	}

	if printBlockchainCmd.Parsed() {
		fmt.Println("xxxxxx")
	}
}
