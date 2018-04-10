package server

import (
	"log"

	"go-chat/chat"
	"go-chat/proto"
)

func signup(conn *chat.ChatConn, content string) {
	user := proto.DecodeUser(content)
	_, dup := clients[user.ID]
	if dup {
		conn.SendErrorf("signup", "ID %s is already taken", user.ID)
		return
	}
	if err := db.Register(user); err != nil {
		conn.SendErrorf("signup", "ID %s is already taken", user.ID)
		return
	}

	conn.SendSuccess("signup", "")
	conn.User = *user
	clients[user.ID] = conn
	connections[conn.RemoteAddr()] = user.ID
	log.Printf("ID %s sign up", conn.ID)
	return
}
