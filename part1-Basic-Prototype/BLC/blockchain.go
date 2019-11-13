package BLC

type BlockChain struct {
	Blocks []*Block // 存储有序的区块
}

// 新增区块
func (blockchain *BlockChain) AddBlock(data string) {
	// 1. 获取区块链中最后一个区块的hash值
	S := len(blockchain.Blocks)
	prevBlockHash := blockchain.Blocks[S-1].Hash
	block := NewBlock(data, prevBlockHash)
	// 2. 添加
	blockchain.Blocks = append(blockchain.Blocks, block)
} 

// 创建一个带有创世区块节点的区块链
func NewBlockchain() *BlockChain {
	return &BlockChain{[]*Block{NewGenesisBlock()}}
}
