package BLC

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
)

const (
	targetBits = 20
	maxNonce   = math.MaxInt64
)

// ProofOfWork 表示一次工作量证明
type ProofOfWork struct {
	block  *Block   // 当前验证的区块
	target *big.Int // 存储大数，作为挖矿的难度
}

// NewProofOfWork 工厂方法，表示新生成一个POW对象
func NewProofOfWork(block *Block) *ProofOfWork {
	target := big.NewInt(1)

	// 左移位
	target.Lsh(target, uint(256-targetBits))

	pow := &ProofOfWork{block, target}
	return pow
}

// 数据拼接
func (p *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join([][]byte{
		p.block.PrevBlockHash,
		p.block.Data,
		IntToHex(p.block.TimeStamp),
		IntToHex(int64(targetBits)),
		IntToHex(int64(nonce)),
	},
		[]byte{},
	)

	return data
}

func (p *ProofOfWork) Run() ([]byte, int) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 9
	fmt.Printf("Mining the block containing \"%s\" \n", p.block.Data)

	for nonce < maxNonce {
		data := p.prepareData(nonce)

		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)

		hashInt.SetBytes(hash[:])
		// hashInt < p.target -1
		//         =          0
		// hashInt > p.target 1
		if hashInt.Cmp(p.target) == -1 {
			break
		} else {
			nonce++
		}
	}

	fmt.Println("")
	fmt.Println("")

	return hash[:], nonce
}

func (p *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := p.prepareData(p.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValidate := hashInt.Cmp(p.target) == -1
	return isValidate
}
