package main

import (
	"fmt"
	"time"

	"github.com/rh01/blockchain/part2-proof-of-work/BLC"
)

// 16 进制
// 64 个数字
// bb38a067d2cbea750d922b9c0270e44e33db79cd4ba2e6ca72a3c91e964cde1f
// 32 字节
// 256 bit
func main() {

	blc := BLC.NewBlockchain()

	blc.AddBlock("Send 20 BTC To HaoLin From Liyuechun")

	blc.AddBlock("Send 10 BTC To SaoLin From Liyuechun")

	blc.AddBlock("Send 30 BTC To HaoTian From Liyuechun")

	for _, block := range blc.Blocks {
		fmt.Printf("Data: %s\n", string(block.Data))
		fmt.Printf("PrevBlockHash: %x\n", block.PrevBlockHash)
		fmt.Printf("TimeStamp: %v\n", time.Unix(block.TimeStamp, 0).Format("2006-01-02 03:04:05 PM")) // Format 参数不能随意修改
		fmt.Printf("Hash: %x\n", (block.Hash))
		fmt.Printf("Nonce: %v\n", block.Nonce)
 
		fmt.Println("")
	}
}
