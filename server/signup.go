package server

import (
	"log"

	"go-chat/chat"
	"go-chat/protocol"
)

func signup(conn *chat.Connection, data *protocol.Data) {
	conn.User = data.User
	mu.Lock()
	defer mu.Unlock()
	_, dup := clients[conn.ID]
	if dup {
		conn.SendErrorf("ID %s is already taken", conn.ID)
		return
	}
	err := db.Register(&conn.User)
	if err != nil {
		conn.SendErrorf("ID %s is already taken", conn.ID)
		return
	}
	clients[conn.ID] = conn
	connIDs[conn.Number] = data.ID
	conn.SendOther("signup")
	log.Printf("ID %s sign up", conn.ID)
	return
}
