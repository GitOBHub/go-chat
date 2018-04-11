package main

import (
	"bufio"
	"os"

	"go-chat/chat"
	"go-chat/client/color"
	"go-chat/proto"
)

func login(conn *chat.ChatConn) {
	in := bufio.NewScanner(os.Stdin)
	for {
		user := new(proto.User)
		color.PrintPrompt(" User ID \n")
		if !in.Scan() {
			os.Exit(0)
		}
		id := in.Text()
		if !isIDValid(id) {
			color.PrintErrorln(" Invalid user ID! input again ")
			continue
		}
		user.ID = id
		//user.Passwd =

		toSend := proto.EncodeUser(user)
		conn.SendRequest("login", toSend)
		resp := <-loginDone
		if resp.Type == proto.Success && resp.Topic == "login" {
			color.PrintPrompt(" Login successfully \n")
			conn.User = *proto.DecodeUser(resp.Content)
			return
		}
		if resp.Type == proto.Error {
			color.PrintErrorln(" %s ", resp.Content)
			continue
		}
		color.PrintErrorln(" Unknown data: %s ", resp.Content)
	}

}
