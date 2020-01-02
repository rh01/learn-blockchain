package main

import (
	"fmt"
	"math/big"

	"github.com/rh01/learn-blockchain/part13-iterate-db-block-chain/BLC"
)

func main() {
	blockchain := BLC.NewBlockChain()

	// blockchain.AddBlock("XiaoMing Send to Xiaohua 20 BTC")
	// blockchain.AddBlock("XiaoMing Send to Xiaohua 20 BTC")
	// blockchain.AddBlock("XiaoMing Send to Xiaohua 20 BTC")
	// blockchain.AddBlock("XiaoMing Send to Xiaohua 20 BTC")
	// fmt.Println(blockchain)

	var blockIterator *BLC.BlockChainIterator
	blockIterator = blockchain.Iterator()
	var hashInt big.Int
	// cmp := big.13(1)
	// cmp.Lsh(x, n)
	for {

		fmt.Printf("%x\n", blockIterator.CurrentHash)

		blockIterator = blockIterator.Next()

		hashInt.SetBytes(blockIterator.CurrentHash)

		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break
		}
	}
}
