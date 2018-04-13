package chat

import (
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gitobhub/net/server"
	"go-chat/proto"
)

type ChatConn struct {
	server.Conn
	proto.User
	Mu sync.Mutex
}

func NewChatConn(c net.Conn) *ChatConn {
	conn := server.NewConn(c)
	return &ChatConn{Conn: *conn}
}

func (chatConn *ChatConn) New(c net.Conn) server.ConnInterface {
	return NewChatConn(c)
}

func (conn *ChatConn) Recv() ([]byte, error) {
	dataLenStr, err := conn.Reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	dataLenStr = strings.TrimSuffix(dataLenStr, "\r\n")
	dataLen, err := strconv.Atoi(dataLenStr)
	if err != nil {
		return nil, err
	}
	data := make([]byte, dataLen)
	conn.Reader.Read(data)
	return data, nil
}

func (conn *ChatConn) ReadData() *proto.Data {
	data, err := conn.Recv()
	if err != nil {
		return nil
	}
	return proto.DecodeData(data)
}

func (conn *ChatConn) SendData(data *proto.Data) (int, error) {
	toSend := string(proto.EncodeData(data))
	toSend = strconv.Itoa(len(toSend)) + "\r\n" + toSend
	return io.WriteString(conn, toSend)
}

//Called by server
func (conn *ChatConn) sendResponse(respType byte, topic, content string) (int, error) {
	data := new(proto.Data)
	data.Type = respType
	data.Time = time.Now().Format("15:04:05")
	data.Topic = topic
	data.Content = content
	return conn.SendData(data)
}

func (conn *ChatConn) SendSuccess(topic, content string) (int, error) {
	return conn.sendResponse(proto.Success, topic, content)
}

func (conn *ChatConn) SendSuccessf(topic, format string, args ...interface{}) (int, error) {
	content := fmt.Sprintf(format, args...)
	return conn.SendSuccess(topic, content)
}

func (conn *ChatConn) SendError(topic, content string) (int, error) {
	return conn.sendResponse(proto.Error, topic, content)
}

func (conn *ChatConn) SendErrorf(topic, format string, args ...interface{}) (int, error) {
	content := fmt.Sprintf(format, args...)
	return conn.SendError(topic, content)
}

//Called by client
func (conn *ChatConn) SendMessageto(msg string, receiver string) (int, error) {
	data := new(proto.Data)
	data.Type = proto.Normal
	data.Sender = conn.User.ID
	data.Receiver = receiver
	data.Time = time.Now().Format("15:04:05")
	data.Content = msg
	return conn.SendData(data)
}

func (conn *ChatConn) SendRequest(topic, content string) (int, error) {
	data := new(proto.Data)
	data.Type = proto.Request
	data.Sender = conn.User.ID
	data.Time = time.Now().Format("15:04:05")
	data.Topic = topic
	data.Content = content
	return conn.SendData(data)
}
