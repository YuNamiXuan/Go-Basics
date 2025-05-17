package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:Cyx20030922.@tcp(127.0.0.1:3306)/bluebell?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	rows, err := db.Query("SELECT title, content FROM post")
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var titile string
		var content string
		err = rows.Scan(&titile, &content)
		if err != nil {
			panic(err.Error())
		}
		fmt.Println(titile, content)
	}
}
