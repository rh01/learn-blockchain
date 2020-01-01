package main

import (
	"github.com/rh01/learn-blockchain/part12-iterate-db-block-chain/BLC"
)

func main() {
	blockchain := BLC.NewBlockChain()

	blockchain.AddBlock("XiaoMing Send to Xiaohua 20 BTC")
	blockchain.AddBlock("XiaoMing Send to Xiaohua 20 BTC")
	blockchain.AddBlock("XiaoMing Send to Xiaohua 20 BTC")
	blockchain.AddBlock("XiaoMing Send to Xiaohua 20 BTC")
	// fmt.Println(blockchain)

	// fmt.Printf("Tip: %x\n", blockchain.Tip)
}
