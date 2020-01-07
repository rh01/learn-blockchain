package main

import (
	"github.com/rh01/learn-blockchain/part30-transaction/BLC"
)

func main() {
	//// 创建区块链
	//blockchain := BLC.NewBlockChain()

	// 创建CLI对象
	cli := BLC.CLI{}

	// 调用CLI的RUn方法
	cli.Run()
}
