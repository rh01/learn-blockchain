// Copyright 2019 The darwin Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package BLC

// Transaction 表示比特币的一笔交易
type Transaction struct {
	ID   []byte     // 交易的ID
	Vin  []TXInput  // 交易的输入
	Vout []TXOutput // 交易的输出,一笔交易有很多交易输出
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

// TXOutput 表示一笔交易的输出
// 交易的输出怎么理解？转账，必须有转账的数额和转账的地址
type TXOutput struct {
	Value        int    // 分
	ScriptPubKey string // 钱包的地址
}
