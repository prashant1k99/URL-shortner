package main

import (
	"fmt"

	"github.com/prashant1k99/URL-Shortner/db"
)

func main() {
	fmt.Println("Hey, Hello World!")
	db.ConnectDB()

	_, err := db.GetCollection("users")
	if err != nil {
		panic(err)
	}
	defer db.DisconnectDB()
}
