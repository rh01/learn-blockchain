package main

import (
	"log"

	"github.com/boltdb/bolt"
)

// 数据库名字
const dbFile = "blc.db"

// 仓库
const blocksBuckets = "blocks"

func main() {
	// 数据库的创建
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()
	
}
