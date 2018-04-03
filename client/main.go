package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
	"unicode"
	"unicode/utf8"

	"server/chat/chat"
	"server/chat/print"
	"server/chat/protocol"
)

var (
	muRecv   sync.Mutex
	received bool
	muSent   sync.Mutex
	sent     bool
	muChat   sync.Mutex
	chating  string
	toChat   string

	msgCache = make(map[string]string)
	isSignup = flag.Bool("signup", false, "Sign up")
)

var idExist = make(chan bool)

//var symbolStr = "‘’“”…·①②③④⑤⑥⑦⑧⑨⑩
var symbols = map[rune]bool{
	'‘': true, '’': true, '“': true,
	'”': true, '…': true,
}

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
		fmt.Println()
		//FIXME: c.Close()
	}()

	conn := &chat.Connection{Conn: c}
	if *isSignup {
		signup(conn)
	} else {
		login(conn)
	}
	go handleConnInput(conn)
	handleKeyInput(conn)
}

func handleKeyInput(conn *chat.Connection) {
	input := bufio.NewScanner(os.Stdin)
	for {
		print.PrintPrompt(" Select a freind for chat ")
		if !input.Scan() {
			break
		}
		friend := input.Text()
		if !isIDValid(friend) {
			print.PrintErrorln(" Invalid user ID! input again ")
			continue
		}
		toChat = friend
		conn.SendOtherto("IsIDExist", friend)
		if ok := <-idExist; !ok {
			toChat = ""
			continue
		}
		changeChating(friend)
		print.PrintPrompt(" Input your message \n")
		var eof bool
		for {
			if eof = !input.Scan(); eof {
				break
			}
			fmt.Printf("\033[1A\033[K")

			muRecv.Lock() //使前后打印出的消息间隔合理
			if !received {
				fmt.Println()
			}
			received = false
			muRecv.Unlock()

			msg := input.Text()
			//debug
			fmt.Printf("client: msg len %d", len(msg))
			if len(msg) == 0 {
				break
			}
			//		printMessageBlock(msg, right)

			muSent.Lock()
			sent = true
			muSent.Unlock()
			n, err := conn.SendMessageto(msg, friend)
			log.Printf("SendMessageto return %d\n", n)
			if err != nil {
				print.PrintErrorln("%s", err)
			}
		}
		if eof {
			break
		}
	}
}

func handleConnInput(conn *chat.Connection) {
	for {
		data := conn.ReadData()
		if data == nil {
			break
		}
		if data.Type == protocol.Error {
			print.PrintErrorln(" %s ", data.Content)
			muChat.Lock()
			info := fmt.Sprintf("ID %s does not exist", toChat)
			muChat.Unlock()
			if data.Content == info {
				idExist <- false
			}
			continue
		}
		if data.Type == protocol.Other {
			switch data.Content {
			case "IsIDExist":
				idExist <- true
			}
			continue
		}
		printMessage(data)
	}
	fmt.Println("\nConnection closed by foreign host")
	os.Exit(0)
}

func changeChating(who string) {
	muChat.Lock()
	chating = who
	msgs, ok := msgCache[who]
	if !ok {
		muChat.Unlock()
		return
	}
	delete(msgCache, who)
	muChat.Unlock()
	fmt.Print(msgs)
}

const (
	right int = iota
	left
)

func printMessage(data *protocol.Data) {
	muChat.Lock()
	if chating != data.ID {
		messageRemind(data)
		muChat.Unlock()
		return
	}
	muChat.Unlock()
	muSent.Lock()
	if !sent {
		fmt.Println()
	}
	sent = false
	muSent.Unlock()

	fmt.Printf("%s ", data.Time)
	fmt.Printf("\033[43;30m%s\033[0m ", data.Name)
	//printMessageBlock(data.Content, left)
	muRecv.Lock()
	received = true
	muRecv.Unlock()
}

func printMessageBlock(msg string, whichSide int) {
	var lenPrint int
	msgRunes := []rune(msg)
	var lines [][]rune
	var nLine int
	for _, r := range msgRunes {
		if _, ok := symbols[r]; ok { //2or3 bytes per rune, printed in 1 space
			lenPrint += 1
		} else if utf8.RuneLen(r) == 3 { //3 bytes per rune, printed in 2 space
			lenPrint += 2
		} else {
			lenPrint += 1
		}
		if lenPrint%15 == 0 {
			lines = append(lines, msgRunes[nLine*15:nLine*15+15])
			nLine++
		}
	}
	if lenPrint-nLine*15 < 15 {
		lines = append(lines, msgRunes[nLine*15:])
	}
	placePrint := 0
	if whichSide == right {
		if lenPrint < 15 {
			placePrint = 50 - lenPrint
		} else {
			placePrint = 35
		}
	}
	for i, line := range lines {
		fmt.Printf("\033[%dC", placePrint)
		lineStr := string(line)
		fmt.Printf("\033[42;30m %s ", lineStr)
		if i == 0 {
			if whichSide == right {
				fmt.Printf("\033[0m %s\n", time.Now().Format("15:04:05"))
				continue
			}
		} else if i == len(lines)-1 && len(lines) > 1 {
			if len(line) < 15 {
				for i := 0; i < 15-len(line); i++ {
					fmt.Print(" ")
				}
			}
		}
		fmt.Println("\033[0m")
		if whichSide == left {
			fmt.Printf("\033[14C")
		}
	}
	//fmt.Print("\n")
}

func messageRemind(data *protocol.Data) {
	fmt.Printf("\033[44;30m New message from %s \033[0m\n", data.Name)
	toPrint := fmt.Sprintf("%s \033[43;30m%s\033[0m \033[7m %s \033[0m\n\n", data.Time, data.Name, data.Content)
	msgCache[data.ID] += toPrint
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

func login(conn *chat.Connection) {
	in := bufio.NewScanner(os.Stdin)
	var userData protocol.User
	for {
		print.PrintPrompt(" User ID ")
		if !in.Scan() {
			os.Exit(0)
		}
		id := in.Text()
		if !isIDValid(id) {
			print.PrintErrorln(" Invalid user ID! input again ")
			continue
		}
		userData.ID = id
		sendUserData(conn, &userData, "login")
		resp := conn.ReadData()
		if resp == nil {
			log.Print("\nConnection closed by foreign host")
			os.Exit(0)
		}
		if resp.Type == protocol.Other && resp.Content == "login" {
			print.PrintPrompt(" Login successfully \n")
			conn.User = resp.User
			return
		}
		if resp.Type == protocol.Error {
			print.PrintErrorln(" %s ", resp.Content)
			continue
		}
		print.PrintErrorln(" Unknown data: %s ", resp.Content)
	}

}

func signup(conn *chat.Connection) {
	print.PrintPrompt(" Sign up \n")
	var userData protocol.User
	in := bufio.NewScanner(os.Stdin)
	for {
		print.PrintPrompt(" User ID ")
		if !in.Scan() {
			os.Exit(0)
		}
		id := in.Text()
		if !isIDValid(id) {
			print.PrintErrorln(" Invalid user ID! input again ")
			continue
		}
		userData.ID = id
		//
		print.PrintPrompt(" Username ")
		if !in.Scan() {
			os.Exit(0)
		}
		userData.Name = in.Text()
		//
		print.PrintPrompt(" Sex ")
		if !in.Scan() {
			os.Exit(0)
		}
		userData.Sex = in.Text()
		//
		print.PrintPrompt("Birth")
		if !in.Scan() {
			os.Exit(0)
		}
		userData.Birth = in.Text()
		//
		sendUserData(conn, &userData, "signup")
		resp := conn.ReadData()
		if resp == nil {
			log.Print("\nConnection closed by foreign host")
			os.Exit(0)
		}
		if resp.Type == protocol.Other && resp.Content == "signup" {
			print.PrintPrompt(" Sign up finished \n")
			return
		}
		if resp.Type == protocol.Error {
			print.PrintErrorln(" %s ", resp.Content)
			continue
		}
		print.PrintErrorln(" BUG! unknown data: %s", resp.Content)
	}
}

func sendUserData(conn *chat.Connection, user *protocol.User, event string) {
	conn.User = *user
	conn.SendOther(event)
}
