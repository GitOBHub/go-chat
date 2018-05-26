package main

import (
	//	"flag"
	"fmt"
	"log"
	"os"

	"go-chat/conf"
	"go-chat/database"
	"go-chat/server"
)

//var port = flag.String("port", "5000", "port")
//var dsn = flag.String("dsn", "", "data source name")

func main() {
	//	flag.Parse()
	laddr := conf.Server.Addr

	db, err := database.OpenMySQL(conf.Server.DSN)
	if err != nil {
		log.Fatal(err)
	}

	srv := server.NewChatServer(laddr, db)
	fmt.Printf("pid: %d    port: %s\n", os.Getpid(), conf.Server.Port)
	log.Fatal(srv.ListenAndServe())
}
