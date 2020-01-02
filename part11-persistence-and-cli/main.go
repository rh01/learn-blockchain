package main

import (
	"fmt"

	"github.com/rh01/learn-blockchain/part11-persistence-and-cli/BLC"
)

func main() {
	blockchain := BLC.NewBlockChain()
	fmt.Println(blockchain)

	fmt.Printf("Tip: %x\n", blockchain.Tip)

	blockchain.AddBlock("XiaoMing Send to Xiaohua 20 BTC")
	// fmt.Println(blockchain)

	fmt.Printf("Tip: %x\n", blockchain.Tip)
}
