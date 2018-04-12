package server

import (
	"go-chat/chat"
)

func (srv *ChatServerHandler) isIDExist(conn *chat.ChatConn, id string) {
	if !srv.db.IsIDExist(id) {
		conn.SendErrorf("isIDExist", "ID %s does not exist", id)
		return
	}
	conn.SendSuccess("isIDExist", "")
}
