package server

import (
	"log"

	"go-chat/chat"
	"go-chat/proto"
)

func (srv *ChatServerHandler) login(conn *chat.ChatConn, content string) {
	user := proto.DecodeUser(content)
	_, dup := srv.clients[user.ID]
	if dup {
		conn.SendErrorf("login", "ID %s is online", user.ID)
		return
	}
	userSaved := srv.db.UserData(user.ID)
	if userSaved == nil {
		conn.SendErrorf("login", "ID %s does not exist", user.ID)
		return
	}

	log.Printf("ID %s login", user.ID)
	datas := srv.db.RestoreMessage(user.ID)
	for _, d := range datas {
		conn.SendData(&d)
	}

	user = userSaved
	toSend := proto.EncodeUser(user)
	conn.SendSuccess("login", toSend)
	conn.User = *user
	srv.clients[user.ID] = conn
	return
}
