package main

import (
	"fmt"

	"github.com/rh01/learn-blockchain/part10-persistence-and-cli/BLC"
)

func main() {
	blockchain := BLC.NewBlockChain()
	fmt.Println(blockchain)

	fmt.Printf("Tip: %x\n", blockchain.Tip)
}
