package main

import (
	"fmt"

	"github.com/prashant1k99/URL-Shortner/db"
)

func main() {
	fmt.Println("Hey, Hello World!")
	db.ConnectDB()
	defer db.DisconnectDB()
}
