package server

import (
	"log"

	"go-chat/chat"
	"go-chat/protocol"
)

func login(conn *chat.Connection, data *protocol.Data) {
	log.Print("Enter login()")
	_, dup := clients[data.ID]
	if dup {
		conn.SendErrorf("ID %s is online", data.ID)
		return
	}
	userData := db.UserData(data.ID)
	if userData == nil {
		conn.SendErrorf("ID %s does not exist", data.ID)
		return
	}
	clients[data.ID] = conn
	connIDs[conn.Number] = data.ID
	conn.User = *userData
	conn.SendOther("login")
	log.Printf("ID %s login", conn.ID)
	datas := db.MessagePreserved(data.ID)
	for _, d := range datas {
		conn.SendData(&d)
	}
	return
}
