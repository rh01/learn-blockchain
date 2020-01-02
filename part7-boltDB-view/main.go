package main

import (
	"fmt"
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
		log.Fatalf("database %v not found", blocksBuckets)
	}
	defer db.Close()

	// 查询数据
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBuckets))
		value := b.Get([]byte("shh"))
		if value == nil {
			return fmt.Errorf("No value found")
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

}
