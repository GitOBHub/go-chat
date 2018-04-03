package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"server/chat/chat"
	"server/chat/database"
	"server/chat/protocol"
)

var (
	port    = flag.String("port", "5000", "port")
	clients = make(map[string]*chat.Connection, 10)
	mu      sync.Mutex
)

const passwd = "root:MYSQLtianwenjie@/chat"

func main() {
	flag.Parse()
	laddr := ":" + *port
	fmt.Printf("pid: %d    port: %s\n", os.Getpid(), *port)
	listener, err := net.Listen("tcp4", laddr)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := database.OpenMysql(passwd); err != nil {
		log.Fatal(err)
	}
	eventLoop(listener)
}

func eventLoop(ln net.Listener) {
	var nConn int
	for {
		c, err := ln.Accept()
		if err != nil {
			log.Print(err)
			continue
		}

		nConn++
		conn := &chat.Connection{Conn: c, Number: nConn}
		log.Printf("connection#%d is up", conn.Number)
		go handleConn(conn)
	}
}

//func register(data string) {

//}

func handleConn(conn *chat.Connection) {
	defer func() {
		conn.Close()
		log.Printf("connection#%d is down", conn.Number)
		log.Print(clients)
	}()

	if !confirmUser(conn) {
		return
	}
	mu.Lock()
	datas := database.MessagePreserved(conn.ID)
	mu.Unlock()
	if datas != nil {
		for _, data := range datas {
			conn.SendData(&data)
		}
	} else {
		//		log.Fatal("get failed")
	}
	for {
		data := conn.ReadData()
		if data == nil {
			break
		}
		if data.Type == protocol.Error {
			conn.SendError("Bad request")
			break
		}
		mu.Lock()
		if data.Type == protocol.Other && data.Content == "IsIDExist" {
			if !database.IsIDExist(data.Receiver.ID) {
				conn.SendErrorf("ID %s does not exist", data.Receiver.ID)
				mu.Unlock()
				continue
			}
			conn.SendOther("IsIDExist")
		}
		client, ok := clients[data.Receiver.ID]
		if !ok {
			conn.SendErrorf("%s is offline", data.Receiver.ID)
			if data.Type == protocol.Normal {
				database.PreserveMessage(data)
			}
			mu.Unlock()
			continue
		}
		mu.Unlock()
		if data.Type == protocol.Normal {
			data.Time = time.Now().Format("15:04:05")
			client.SendData(data)
		}
	}
	mu.Lock()
	delete(clients, conn.ID)
	mu.Unlock()
}

func confirmUser(conn *chat.Connection) bool {
	for {
		data := conn.ReadData()
		if data == nil {
			return false
		}
		if data.Type != protocol.Other {
			conn.SendError("Bad request!")
			return false
		}

		if data.Content == "login" {
			mu.Lock()
			_, ok := clients[data.ID]
			if ok {
				mu.Unlock()
				conn.SendErrorf("ID %s is online", data.ID)
				continue
			}
			userData := database.UserData(data.ID)
			if userData == nil {
				mu.Unlock()
				conn.SendErrorf("ID %s does not exist", data.ID)
				continue
			}
			clients[data.ID] = conn
			mu.Unlock()
			conn.User = *userData
			conn.SendOther("login")
			log.Printf("ID %s login", conn.ID)
			break
		} else if data.Content == "signup" {
			conn.User = data.User
			/*	conn.SendError("Entry closed")
				log.Printf("ID %s signup", conn.ID)
				continue*/
			mu.Lock()
			_, ok := clients[conn.ID]
			if ok {
				mu.Unlock()
				conn.SendErrorf("ID %s is already taken", conn.ID)
				continue
			}
			err := database.Register(&conn.User)
			if err != nil {
				mu.Unlock()
				conn.SendErrorf("ID %s is already taken", conn.ID)
				continue
			}
			clients[conn.ID] = conn
			mu.Unlock()
			conn.SendOther("signup")
			log.Printf("ID %s sign up", conn.ID)
			break
		} else {
			conn.SendError("Bad Request!")
			return false
		}
	}
	return true
}
