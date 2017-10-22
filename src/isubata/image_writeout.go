package main

import (
	crand "crypto/rand"
	"encoding/binary"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"math/rand"
	"os"
	"time"
)

var (
	db *sqlx.DB
)

func main() {
	seedBuf := make([]byte, 8)
	crand.Read(seedBuf)
	rand.Seed(int64(binary.LittleEndian.Uint64(seedBuf)))

	db_host := os.Getenv("ISUBATA_DB_HOST")
	if db_host == "" {
		db_host = "127.0.0.1"
	}
	db_port := os.Getenv("ISUBATA_DB_PORT")
	if db_port == "" {
		db_port = "3306"
	}
	db_user := os.Getenv("ISUBATA_DB_USER")
	if db_user == "" {
		db_user = "root"
	}
	db_password := os.Getenv("ISUBATA_DB_PASSWORD")
	if db_password != "" {
		db_password = ":" + db_password
	}

	dsn := fmt.Sprintf("%s%s@tcp(%s:%s)/isubata?parseTime=true&loc=Local&charset=utf8mb4",
		db_user, db_password, db_host, db_port)

	log.Printf("Connecting to db: %q", dsn)
	db, _ = sqlx.Connect("mysql", dsn)
	for {
		err := db.Ping()
		if err == nil {
			break
		}
		log.Println(err)
		time.Sleep(time.Second * 3)
	}

	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(5 * time.Minute)
	log.Printf("Succeeded to connect db.")

	// 書き出し開始
	rows, err := db.Query("SELECT name, data FROM image")
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var name string
		var data []byte
		rows.Scan(&name, &data)
		file, _ := os.OpenFile("/mnt/images/"+name, os.O_CREATE|os.O_WRONLY, 0666)
		file.Write(data)
		file.Close()
	}
	rows.Close()
}
