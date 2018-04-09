package main

import (
	"bufio"
	"os"

	"go-chat/chat"
	"go-chat/color"
	"go-chat/proto"
)

func signup(conn *chat.ChatConn) {
	color.PrintPrompt(" Sign up \n")
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
		//
		color.PrintPrompt(" Username ")
		if !in.Scan() {
			os.Exit(0)
		}
		user.Name = in.Text()
		//
		color.PrintPrompt(" Sex ")
		if !in.Scan() {
			os.Exit(0)
		}
		user.Sex = in.Text()
		//
		color.PrintPrompt("Birth")
		if !in.Scan() {
			os.Exit(0)
		}
		user.Birth = in.Text()
		//
		toSend := proto.EncodeUser(user)
		conn.SendRequest("signup", toSend)
		resp := <-signupDone
		if resp.Type == proto.Success && resp.Topic == "signup" {
			color.PrintPrompt(" Sign up finished \n")
			conn.User = *user
			return
		}
		if resp.Type == proto.Error {
			color.PrintErrorln(" %s ", resp.Content)
			continue
		}
		color.PrintErrorln(" Unknown data: %s", resp.Content)
	}
}
