package session

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"time"
	"unicode/utf8"

	"go-chat/color"
	"go-chat/proto"
)

type Session struct {
	ID             string
	ToSend         chan string
	unreadMessages []*proto.Message
	mu             sync.Mutex
	newMsg         chan *proto.Message
	closeDone      chan struct{}
	lastPrintSide  int
	lastRemindTime time.Time
}

var (
	mu      sync.Mutex
	running string
)

//var symbolStr = "‘’“”…·①②③④⑤⑥⑦⑧⑨⑩
var symbols = map[rune]bool{
	'‘': true, '’': true, '“': true,
	'”': true, '…': true,
}

const (
	pending int = iota
	right
	left
)

func NewSession(id string) *Session {
	s := &Session{ID: id}
	s.ToSend = make(chan string)
	s.newMsg = make(chan *proto.Message)
	s.closeDone = make(chan struct{})
	return s
}

func (s *Session) Run() {
	mu.Lock()
	defer mu.Unlock()
	if running != "" {
		mu.Unlock()
		return
	}
	running = s.ID

	color.PrintPrompt(" Input \":close\" to close current session \n")
	go s.printPeer()
	go s.interact()
}

func (s *Session) Close() {
	mu.Lock()
	defer mu.Unlock()
	running = ""
	close(s.ToSend)
}

func (s *Session) PutMessage(msg *proto.Message) {
	if running != s.ID {
		s.unreadMessages = append(s.unreadMessages, msg)
		if time.Since(s.lastRemindTime).Seconds() > 2 {
			color.PrintBlueln(" New message from %s ", s.ID)
			s.lastRemindTime = time.Now()
		}
		return
	}
	s.newMsg <- msg
}

func (s *Session) interact() {
	input := bufio.NewScanner(os.Stdin)
in:
	for input.Scan() {
		text := input.Text()
		fmt.Print("\033[1A\033[K")

		switch text {
		case ":close":
			break in
		case "":
			color.PrintErrorln(" Blank input! ")
			continue
		}
		msg := new(proto.Message)
		msg.Time = time.Now().Format("15:04:05")
		msg.Content = text
		s.printMessageBlock(msg, right)
		s.ToSend <- msg.Content
	}
	s.Close()
	s.closeDone <- struct{}{}
}

func (s *Session) printPeer() {
	//FIXME:
	for _, m := range s.unreadMessages {
		s.printMessageBlock(m, left)
	}
	for {
		select {
		case msg := <-s.newMsg:
			s.printMessageBlock(msg, left)
		case <-s.closeDone:
			return
		}
	}
}

func (s *Session) printMessageBlock(msg *proto.Message, whichSide int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.lastPrintSide == whichSide {
		fmt.Println()
	}

	var lenPrint int
	msgRunes := []rune(msg.Content)
	var lines [][]rune
	var eol int
	for i, r := range msgRunes {
		_, ok := symbols[r]
		if utf8.RuneLen(r) == 3 && !ok { //3 bytes per rune, printed in 2 space
			lenPrint += 2
			if lenPrint == 16 {
				print(lenPrint, "\n")
				tempRunes := append(msgRunes[eol:i], rune(' '))
				lines = append(lines, tempRunes)
				eol = i
				lenPrint = 2
			} else if lenPrint == 15 {
				print(lenPrint, "\n")
				lines = append(lines, msgRunes[eol:i+1])
				eol = i + 1
				lenPrint = 0
			}
		} else { //1or2or3 bytes per rune, printed in 1 space
			lenPrint += 1
			if lenPrint == 15 {
				print(lenPrint, "\n")
				lines = append(lines, msgRunes[eol:i+1])
				eol = i + 1
				lenPrint = 0
			}
		}
	}
	lines = append(lines, msgRunes[eol:])

	var placePrint int
	if whichSide == right {
		if len(lines) == 1 && lenPrint < 15 {
			placePrint = 50 - lenPrint
		} else {
			placePrint = 35
		}
	}

	for i, line := range lines {
		if whichSide == right {
			fmt.Printf("\033[%dC", placePrint)
			fmt.Printf("\033[42;30m %s \033[0m", string(line))
			if i == 0 {
				fmt.Printf(" %s\n", time.Now().Format("15:04:05"))
				continue
			}
		} else if whichSide == left {
			if i == 0 {
				fmt.Printf("%s \033[43;30m%s\033[0m ", msg.Time, s.ID)
				color.PrintPrompt(" %s \n", string(line))
				placePrint = len(msg.Time) + len(s.ID) + 2
				continue
			}
			fmt.Printf("\033[%dC", placePrint)
			color.PrintPrompt(" %s ", string(line))
		}
		if i == len(lines)-1 && len(lines) > 1 && lenPrint < 15 {
			for i := 0; i < 15-lenPrint; i++ {
				if whichSide == right {
					fmt.Print("\033[42;30m \033[0m")
				} else if whichSide == left {
					color.PrintPrompt(" ")
				}
			}
		}
		fmt.Println()
	}
	s.lastPrintSide = whichSide
}
