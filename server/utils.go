package server

import (
	"go-chat/chat"
)

func isIDExist(conn *chat.ChatConn, id string) {
	if !db.IsIDExist(id) {
		conn.SendErrorf("isIDExist", "ID %s does not exist", id)
		return
	}
	conn.SendSuccess("isIDExist", "")
}
