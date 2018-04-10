package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"unicode"

	"go-chat/chat"
	"go-chat/client/session"
	"go-chat/color"
	"go-chat/proto"
)

var (
	isSignup = flag.Bool("signup", false, "Sign up")

	mu       sync.Mutex
	sessions = make(map[string]*session.Session)

	loginDone  = make(chan *proto.Data)
	signupDone = make(chan *proto.Data)
	idExist    = make(chan *proto.Data)
)

func main() {
	/*	if len(os.Args) == 1 {
		fmt.Printf("usage: %s host:port\n", os.Args[0])
		return
	}*/
	flag.Parse()
	c, err := net.Dial("tcp", "localhost:5000")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		//fmt.Println()
		//FIXME: c.Close()
	}()

	conn := chat.NewChatConn(c)
	go handleConnInput(conn)
	if *isSignup {
		signup(conn)
	} else {
		login(conn)
	}
	startChat(conn)
}

func startChat(conn *chat.ChatConn) {
	input := bufio.NewScanner(os.Stdin)
	for {
		if len(sessions) > 0 {
			color.PrintBlueln(" Unread message from ")
		}
		for _, sess := range sessions {
			color.PrintBlue(" %s ", sess.ID)
			fmt.Print("\t")
		}
		if len(sessions) > 0 {
			fmt.Println()
		}

		color.PrintPrompt(" Select a freind for chat \n")
		if !input.Scan() {
			break
		}
		friend := input.Text()
		if !isIDValid(friend) {
			color.PrintErrorln(" Invalid user ID! input again ")
			continue
		}
		conn.SendRequest("isIDExist", friend)
		resp := <-idExist
		if resp.Type == proto.Error {
			color.PrintErrorln(" %s ", resp.Content)
			continue
		}
		sess, ok := sessions[friend]
		if !ok {
			sess = session.NewSession(friend)
			sessions[friend] = sess
		}
		sess.Run()
		for msg := range sess.ToSend {
			_, err := conn.SendMessageto(msg, friend)
			if err != nil {
				color.PrintErrorln("%s", err)
			}
		}
		delete(sessions, sess.ID)
	}
}

func handleConnInput(conn *chat.ChatConn) {
	for {
		data := conn.ReadData()
		if data == nil {
			break
		}

		if data.Type != proto.Normal {
			switch data.Topic {
			case "isIDExist":
				idExist <- data
			case "login":
				loginDone <- data
			case "signup":
				signupDone <- data
			default:
				if data.Type == proto.Error {
					color.PrintErrorln(" %s ", data.Content)
				} else {
					color.PrintPrompt(" %s \n", data.Content)
				}
			}
			continue
		}
		prepareMessage(data)
	}
	fmt.Println("Connection closed by foreign host")
	os.Exit(0)
}

func isIDValid(id string) bool {
	if len(id) == 0 {
		return false
	}
	rs := []rune(id)
	for _, r := range rs {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	/*	if strings.Contains(name, " ") {
		fmt.Println("username cannot contain \" \"")
		return false
	}*/
	return true
}
