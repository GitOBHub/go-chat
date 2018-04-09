package server

import (
	"log"
	"sync"
	"time"

	"github.com/GitOBHub/net/conns"
	"github.com/GitOBHub/net/server"
	"go-chat/chat"
	"go-chat/database"
	"go-chat/proto"
)

var (
	clients map[string]*chat.ChatConn
	connIDs map[int]string
	db      *database.DB
	mu      sync.Mutex
)

func NewChatServer(addr string, d *database.DB) *server.Server {
	s := server.NewServer(addr)
	clients = make(map[string]*chat.ChatConn, 10)
	connIDs = make(map[int]string, 10)
	db = d
	s.MessageHandleFunc(HandleMessage)
	s.ConnectionHandleFunc(HandleConnection)
	return s
}

func HandleMessage(c *conns.Conn, b []byte) {
	log.Print("Enter handleMessage")
	conn := &chat.ChatConn{Conn: *c}
	data := proto.DecodeData(b)
	if data.Type == proto.Error || data.Type == proto.Success {
		conn.SendError("", "Bad request")
		return
	}
	mu.Lock()
	defer mu.Unlock()

	if data.Type == proto.Request {
		log.Print("data is Request")
		switch data.Topic {
		case "isIDExist":
			isIDExist(conn, data.Content)
		case "login":
			login(conn, data.Content)
		case "signup":
			signup(conn, data.Content)
		}
		return
	}
	client, ok := clients[data.Receiver]
	if !ok {
		conn.SendErrorf("", "%s is offline", data.Receiver)
		if data.Type == proto.Normal {
			db.PreserveMessage(data)
		}
		return
	}
	if data.Type == proto.Normal {
		data.Time = time.Now().Format("15:04:05")
		client.SendData(data)
	}
}

func HandleConnection(c *conns.Conn) {
	/*if !c.Connected {
		mu.Lock()
		defer mu.Unlock()
		name, ok := connIDs[c.Number]
		if !ok {
			return
		}
		delete(clients, name)
		delete(connIDs, c.Number)
		log.Print(clients)
	}*/
}
