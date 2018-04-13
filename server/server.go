package server

import (
	"log"
	"sync"
	"time"

	"github.com/gitobhub/net/server"
	"go-chat/chat"
	"go-chat/database"
	"go-chat/proto"
)

type ChatServerHandler struct {
	clients map[string]*chat.ChatConn
	db      *database.DB
	mu      sync.Mutex
}

func NewChatServer(addr string, d *database.DB) *server.Server {
	cs := new(ChatServerHandler)
	cs.clients = make(map[string]*chat.ChatConn, 10)
	cs.db = d
	s := server.NewServer(addr, cs)
	c := new(chat.ChatConn)
	s.SetConnType(c)
	return s
}

func (srv *ChatServerHandler) HandleMessage(c server.ConnInterface, b []byte) {
	conn := c.(*chat.ChatConn)
	data := proto.DecodeData(b)
	if data.Type == proto.Error || data.Type == proto.Success {
		conn.SendError("", "Bad request")
		return
	}
	srv.mu.Lock()
	defer srv.mu.Unlock()

	if data.Type == proto.Request {
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

func (srv *ChatServerHandler) HandleConn(c server.ConnInterface) {
	conn := c.(*chat.ChatConn)
	if !conn.Connected {
		srv.mu.Lock()
		defer srv.mu.Unlock()
		delete(srv.clients, conn.User.ID)
		log.Print(srv.clients)
	}
}
