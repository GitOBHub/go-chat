package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"server/chat/database"
	"server/chat/server/server"
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
	srv, err := server.NewServer(laddr, db)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("pid: %d    port: %s\n", os.Getpid(), *port)
	srv.Serve()
}
