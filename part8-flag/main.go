// Copyright 2019 The darwin Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {

	// s.Args 提供了原始命令行参数访问功能，注意，切片的第一个
	// 参数是该程序的名称
	argsWithProg := os.Args
	argsWithoutProg := os.Args[1:]

	// 你也可以使用标准的索引位置方式取得单个参数的值
	arg := os.Args[3]
	fmt.Println(argsWithProg)
	fmt.Println(argsWithoutProg)
	fmt.Println(arg)

	flagString := flag.String("printchain", "", "输出所有的区块信息.....")

	flagInt := flag.Int("number", 6, "输出一个整数....")

	flagBool := flag.Bool("open", false, "判断真假....")

	flag.Parse()
	fmt.Printf("%s\n", *flagString)
	fmt.Printf("%d\n", *flagInt)
	fmt.Printf("%v\n", *flagBool)
	fmt.Printf("tail %v\n", flag.Args())

	flag.Usage()
}
