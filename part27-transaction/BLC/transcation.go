// Copyright 2019 The darwin Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package BLC

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

// Transaction 表示比特币的一笔交易
type Transaction struct {
	ID   []byte     // 交易的ID
	Vin  []TXInput  // 交易的输入
	Vout []TXOutput // 交易的输出,一笔交易有很多交易输出
}

// NewCoinbaseTx 创世区块交易信息的初始化,这是一个特殊的交易
// coinbase的交易
func NewCoinbaseTx(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to %s", to)
	}

	txin := TXInput{
		Txid:      []byte{}, // 交易的id，创世区块的交易id为nil
		Vout:      -1,       // 创世区块的输出交易的索引为-1，即不存在
		ScriptSig: data,     // 创世区块的交易信息
	}

	txout := TXOutput{
		Value:        subsidy,
		ScriptPubKey: to,
	}

	tx := Transaction{
		ID:   nil,
		Vin:  []TXInput{txin},
		Vout: []TXOutput{txout},
	}

	tx.SetID()

	return &tx
}

// 新建一个转账交易
func NewUTXOTransaction(from, to string, amount int, bc *BlockChain, txs []*Transaction) *Transaction {
	// 输入
	var inputs []TXInput
	// 输出
	var outputs []TXOutput

	// 找到有效的可用的交易输出数据模型
	// 查询未花费的输出
	// 10TXOTra
	// map[0ffec609bafae305b2f9be3bae96ff6c7e78a5a8fb999bf4c7210c127f6fe62e:[0]]
	acc, validOutputs := bc.FindSpendableOutputs(from, amount, txs)
	if acc < amount {
		log.Panic("ERROR: not enough founds")
	}

	// 建立输入
	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			log.Panic(err)
		}

		for _, out := range outs {
			// 创建一个输入
			input := TXInput{
				Txid:      txID,
				Vout:      out,
				ScriptSig: from,
			}
			// 将输入添加到inputs数组中去
			inputs = append(inputs, input)
		}
	}

	// 建立输出,转账
	output := TXOutput{
		Value:        amount,
		ScriptPubKey: to,
	}
	outputs = append(outputs, output)
	// 建立输出，找零
	output = TXOutput{
		Value:        acc - amount,
		ScriptPubKey: from,
	}

	outputs = append(outputs, output)

	// 创建交易
	tx := Transaction{
		ID:   nil,
		Vin:  inputs,
		Vout: outputs,
	}
	tx.SetID()

	return &tx
}

// SetID 设置交易的ID
func (tx *Transaction) SetID() {
	var encoder bytes.Buffer
	var hash [32]byte

	enc := gob.NewEncoder(&encoder)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	// 将序列化以后的字节数组生成256hash
	hash = sha256.Sum256(encoder.Bytes())
	tx.ID = hash[:]
}

// 新建新的UTXO交易，转账
// 1. 先找到包含当前用户未花费的所有交易的集合
// 2. 找到用户足够的余额所对应的未花费输出
// 未花费输出：TXOutput没有对应的TXInput
// 3. 12, {"1111":[1,2,3]}
// 4. 新建一个输入
// 5. 新建输出
//   1. TXOuput:
//   2. TXOuput:

// 1. 判断当前交易是否是CoinbaseTX
func (tx *Transaction) isCoinbase() bool {
	return len(tx.Vin) == 1 &&
		tx.Vin[0].Vout == -1 && len(tx.Vin[0].Txid) == 0
}

// TXInput 表示一笔交易的输入
// 假设交易输入对应交易输出，那么如何绑定这两者？
// ID -> transcation (Found)
// Transcation -> Vout.index -> TXOutput (Found)
//
type TXInput struct {
	Txid      []byte // 交易的ID
	Vout      int    // 存储TXOutput里面的索引
	ScriptSig string // 存储TXInput
}

// CanUnlockOutputWith 检查账号地址
func (in *TXInput) CanUnlockOutputWith(unlockingData string) bool {
	return in.ScriptSig == unlockingData
}

// TXOutput 表示一笔交易的输出
// 交易的输出怎么理解？转账，必须有转账的数额和转账的地址
type TXOutput struct {
	Value        int    // 分
	ScriptPubKey string // 钱包的地址
}

// CanBeUnlockedWith 检查是否能够检索账号
func (out *TXOutput) CanBeUnlockedWith(unlockingData string) bool {
	return out.ScriptPubKey == unlockingData
}
