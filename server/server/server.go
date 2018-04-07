package server

import (
	"log"
	"net"
	"sync"
	"time"

	"server/chat/chat"
	"server/chat/database"
	"server/chat/protocol"
)

type Server struct {
	listener net.Listener
	Address  string
	NumConn  int //uint64
	mu       sync.Mutex
	db       *database.DB
	clients  map[string]*chat.Connection
}

func NewServer(addr string, db *database.DB) (*Server, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	cs := make(map[string]*chat.Connection, 10)
	srv := &Server{listener: ln, Address: addr, clients: cs, db: db}
	return srv, nil
}

func (s *Server) Serve() {
	for {
		c, err := s.listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		s.mu.Lock()
		s.NumConn++
		s.mu.Unlock()
		conn := &chat.Connection{Conn: c, Number: s.NumConn}
		log.Printf("connection#%d is up", conn.Number)
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn *chat.Connection) {
	defer func() {
		conn.Close()
		log.Printf("connection#%d is down", conn.Number)
		log.Print(s.clients)
	}()
	for {
		data := conn.ReadData()
		if data == nil {
			break
		}
		s.handleMessage(conn, data)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.clients, conn.ID)
}

func (s *Server) handleMessage(conn *chat.Connection, data *protocol.Data) {
	if data.Type == protocol.Error { //|| data.Type == protocol.Success {
		conn.SendError("Bad request")
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	if data.Type == protocol.Other {
		log.Print("data is Ohter")
		switch data.Content {
		case "IsIDExist":
			if !s.db.IsIDExist(data.Receiver.ID) {
				conn.SendErrorf("ID %s does not exist", data.Receiver.ID)
				return
			}
			conn.SendOther("IsIDExist")
		case "login":
			s.login(conn, data)
		case "signup":
			s.signup(conn, data)
		}
		return
	}
	client, ok := s.clients[data.Receiver.ID]
	if !ok {
		conn.SendErrorf("%s is offline", data.Receiver.ID)
		if data.Type == protocol.Normal {
			s.db.PreserveMessage(data)
		}
		return
	}
	if data.Type == protocol.Normal {
		data.Time = time.Now().Format("15:04:05")
		client.SendData(data)
	}
}
