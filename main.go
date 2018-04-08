package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"go-chat/database"
	"go-chat/server"
)

var port = flag.String("port", "5000", "port")
var dsn = flag.String("dsn", "", "data source name")

func main() {
	flag.Parse()
	laddr := ":" + *port

	db, err := database.OpenMySQL(*dsn)
	if err != nil {
		log.Fatal(err)
	}

	srv := server.NewChatServer(laddr, db)
	fmt.Printf("pid: %d    port: %s\n", os.Getpid(), *port)
	log.Fatal(srv.ListenAndServe())
}
