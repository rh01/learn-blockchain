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
		log.Panic(err)
	}
	defer db.Close()

	// 插入更新数据库
	err = db.Update(func(tx *bolt.Tx) error {

		// 查看仓库是否存在
		b := tx.Bucket([]byte(blocksBuckets))
		// 若没有存在
		if b == nil {
			fmt.Println("No existing blc found. Creating")
		}

		// 创建表
		b, err = tx.CreateBucket([]byte(blocksBuckets))
		if err != nil {
			log.Fatal(err)
		}

		// 存储数据 kv
		err = b.Put([]byte("shh"), []byte("https://41sh.cn"))
		if err != nil {
			log.Fatal(err)
		}

		err = b.Put([]byte("shh2"), []byte("https://shenhengheng.cn"))
		if err != nil {
			log.Fatal(err)
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}
}
