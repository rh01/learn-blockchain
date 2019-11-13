package main

import (
	"fmt"

	"github.com/rh01/blockchain/part3-serialize-and-deserialize/BLC"
)

func main() {
	block := BLC.Block{[]byte("send 3 BTC to ZhangBo Grom Shh"), 0}

	fmt.Printf("%s\n", block.Data)
	fmt.Printf("%d\n", block.Nonce)
	fmt.Println("")

	bytes := block.Serialize()
	fmt.Println(bytes)

	blockDes := BLC.Deserialize(bytes)
	fmt.Printf("%s\n", blockDes.Data)
	fmt.Printf("%d\n", blockDes.Nonce)
	fmt.Println("")

}
