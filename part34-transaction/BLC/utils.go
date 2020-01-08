// Copyright 2019 The darwin Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package BLC

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"log"
)

// 将int64转为字节数组
func IntToHex(target int64) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, target)
	if err != nil {
		log.Panic(err)
	}

	return buf.Bytes()
}

// 将json转array
func JSONToArray(jsonStr string) []string {
	var array []string
	if err := json.Unmarshal([]byte(jsonStr), &array); err != nil {
		log.Panic(err)
	}

	return array
}
