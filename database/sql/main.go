package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// Define a struct to hold row data
type myRow struct {
	id   int
	name string
}

func main() {
	// Establish database connection
	// Format: "username:password@protocol(address:port)/database?parameters"
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/database?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println("Connect to database failed: ", err)
		return
	}
	defer db.Close()

	// Execute SQL query to select data
	rows, err := db.Query("SELECT id, name FROM your_table")
	if err != nil {
		fmt.Println("Execute SQL failed: ", err)
		return
	}
	defer rows.Close()

	// Iterate through each row in the result set
	for rows.Next() {
		var row myRow

		// Scan copies the current row's columns into the struct fields
		// Note: The order and number of arguments must match the SELECT columns
		err = rows.Scan(&row.id, &row.name)
		if err != nil {
			fmt.Println("Copy columns into struct fields failed: ", err)
			return
		}
		fmt.Println(row)
	}

	// Check for any errors that occurred during iteration
	if err := rows.Err(); err != nil {
		fmt.Println("Iterate rows failed: ", err)
		return
	}
}
