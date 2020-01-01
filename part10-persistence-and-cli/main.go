package main

import (
	"blockchain/part10-persistence-and-cli/blc"
	"fmt"
)

func main() {
	blockchain := blc.NewBlockChain()
	fmt.Println(blockchain)

	fmt.Printf("Tip: %x\n", blockchain.Tip)
}
