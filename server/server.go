package server

import (
	"log"
	"net"
	"sync"
	"time"

	"github.com/GitOBHub/net/conns"
	"github.com/GitOBHub/net/server"
	"go-chat/chat"
	"go-chat/database"
	"go-chat/proto"
)

type ChatServer struct {
	clients     map[string]*chat.ChatConn
	connections map[net.Addr]string
	db          *database.DB
	mu          sync.Mutex
}

func NewChatServer(addr string, d *database.DB) *server.Server {
	cs := new(ChatServer)
	cs.clients = make(map[string]*chat.ChatConn, 10)
	cs.connections = make(map[net.Addr]string, 10)
	cs.db = d
	s := server.NewServer(addr, cs)
	return s
}

func (srv *ChatServer) HandleMessage(c *conns.Conn, b []byte) {
	log.Print("Enter handleMessage")
	conn := &chat.ChatConn{Conn: *c}
	data := proto.DecodeData(b)
	if data.Type == proto.Error || data.Type == proto.Success {
		conn.SendError("", "Bad request")
		return
	}
	srv.mu.Lock()
	defer srv.mu.Unlock()

	if data.Type == proto.Request {
		log.Print("data is Request")
		switch data.Topic {
		case "isIDExist":
			srv.isIDExist(conn, data.Content)
		case "login":
			srv.login(conn, data.Content)
		case "signup":
			srv.signup(conn, data.Content)
		}
		return
	}
	client, ok := srv.clients[data.Receiver]
	if !ok {
		conn.SendErrorf("", "%s is offline", data.Receiver)
		if data.Type == proto.Normal {
			srv.db.PreserveMessage(data)
		}
		return
	}
	if data.Type == proto.Normal {
		data.Time = time.Now().Format("15:04:05")
		client.SendData(data)
	}
}

func (srv *ChatServer) HandleConn(c *conns.Conn) {
	if !c.Connected {
		srv.mu.Lock()
		defer srv.mu.Unlock()
		name, ok := srv.connections[c.RemoteAddr()]
		if !ok {
			return
		}
		delete(srv.clients, name)
		delete(srv.connections, c.RemoteAddr())
		log.Print(srv.clients)
	}
}
