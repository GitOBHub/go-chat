package server

import (
	"log"

	"server/chat/chat"
	"server/chat/protocol"
)

func (s *ChatServer) login(conn *chat.Connection, data *protocol.Data) {
	log.Print("Enter login()")
	_, dup := s.clients[data.ID]
	if dup {
		conn.SendErrorf("ID %s is online", data.ID)
		return
	}
	userData := s.db.UserData(data.ID)
	if userData == nil {
		conn.SendErrorf("ID %s does not exist", data.ID)
		return
	}
	s.clients[data.ID] = conn
	conn.User = *userData
	conn.SendOther("login")
	log.Printf("ID %s login", conn.ID)
	return
}

func (s *ChatServer) signup(conn *chat.Connection, data *protocol.Data) {
	conn.User = data.User
	s.mu.Lock()
	defer s.mu.Unlock()
	_, dup := s.clients[conn.ID]
	if dup {
		conn.SendErrorf("ID %s is already taken", conn.ID)
		return
	}
	err := s.db.Register(&conn.User)
	if err != nil {
		conn.SendErrorf("ID %s is already taken", conn.ID)
		return
	}
	s.clients[conn.ID] = conn
	conn.SendOther("signup")
	log.Printf("ID %s sign up", conn.ID)
	return
}
