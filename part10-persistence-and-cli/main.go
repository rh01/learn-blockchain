package main

import (
	"fmt"

	"github.com/rh01/learn-blockchain/part10-persistence-and-cli/blc"
)

func main() {
	blockchain := blc.NewBlockChain()
	fmt.Println(blockchain)

	fmt.Printf("Tip: %x\n", blockchain.Tip)
}
