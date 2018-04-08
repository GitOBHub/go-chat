package server

import (
	"log"
	"sync"
	"time"

	"github.com/GitOBHub/net/conns"
	"github.com/GitOBHub/net/server"
	"go-chat/chat"
	"go-chat/database"
	"go-chat/protocol"
)

var (
	clients map[string]*chat.Connection
	connIDs map[int]string
	db      *database.DB
	mu      sync.Mutex
)

func NewChatServer(addr string, d *database.DB) *server.Server {
	s := server.NewServer(addr)
	clients = make(map[string]*chat.Connection, 10)
	connIDs = make(map[int]string, 10)
	db = d
	s.MessageHandleFunc(HandleMessage)
	s.ConnectionHandleFunc(HandleConnection)
	return s
}

func HandleMessage(c *conns.Connection, b []byte) {
	log.Print("Enter handleMessage")
	conn := chat.NewConn(c)
	data := protocol.DecodeData(b)
	if data.Type == protocol.Error { //|| data.Type == protocol.Success {
		conn.SendError("Bad request")
		return
	}
	mu.Lock()
	defer mu.Unlock()

	if data.Type == protocol.Other {
		log.Print("data is Ohter")
		switch data.Content {
		case "IsIDExist":
			if !db.IsIDExist(data.Receiver.ID) {
				conn.SendErrorf("ID %s does not exist", data.Receiver.ID)
				return
			}
			conn.SendOther("IsIDExist")
		case "login":
			login(conn, data)
		case "signup":
			signup(conn, data)
		}
		return
	}
	client, ok := clients[data.Receiver.ID]
	if !ok {
		conn.SendErrorf("%s is offline", data.Receiver.ID)
		if data.Type == protocol.Normal {
			db.PreserveMessage(data)
		}
		return
	}
	if data.Type == protocol.Normal {
		data.Time = time.Now().Format("15:04:05")
		client.SendData(data)
	}
}

func HandleConnection(c *conns.Connection) {
	if !c.Connected {
		mu.Lock()
		defer mu.Unlock()
		name, ok := connIDs[c.Number]
		if !ok {
			return
		}
		delete(clients, name)
		delete(connIDs, c.Number)
		log.Print(clients)
	}
}
