package server

import (
	"log"

	"go-chat/chat"
	"go-chat/proto"
)

func login(conn *chat.ChatConn, content string) {
	log.Print("Enter login()")
	user := proto.DecodeUser(content)
	_, dup := clients[user.ID]
	if dup {
		conn.SendErrorf("login", "ID %s is online", user.ID)
		return
	}
	userSaved := db.UserData(user.ID)
	if userSaved == nil {
		conn.SendErrorf("login", "ID %s does not exist", user.ID)
		return
	}

	log.Printf("ID %s login", user.ID)
	datas := db.RestoreMessage(user.ID)
	for _, d := range datas {
		conn.SendData(&d)
	}

	user = userSaved
	toSend := proto.EncodeUser(user)
	conn.SendSuccess("login", toSend)
	conn.User = *user
	clients[user.ID] = conn
	//	connIDs[conn.Number] = user.ID
	return
}
