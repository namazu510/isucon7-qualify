package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"os"
	"time"
	_ "github.com/go-sql-driver/mysql"
)

var (
	db *sqlx.DB
)

func main() {
	db_host := "127.0.0.1"
	db_port := "3306"
	db_user := "isucon"
	db_password := ":isucon"
	dsn := fmt.Sprintf("%s%s@tcp(%s:%s)/isubata?parseTime=true&loc=Local&charset=utf8mb4",
		db_user, db_password, db_host, db_port)

	log.Printf("Connecting to db: %q", dsn)
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		panic(err)
	}
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
		file, err := os.OpenFile("/srv/images/"+name, os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			panic(err)
		}
		file.Write(data)
		file.Close()
	}
	rows.Close()
}
