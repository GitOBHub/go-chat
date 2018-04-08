package chat

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/GitOBHub/net/conns"
	"go-chat/protocol"
)

type Connection struct {
	conns.Connection
	protocol.User
	Mu sync.Mutex
}

func NewConn(c *conns.Connection) *Connection {
	return &Connection{Connection: *c}
}

func (conn *Connection) ReadData() *protocol.Data {
	data, err := conn.Recv()
	if err != nil {
		log.Print("*Connection.ReadData: Read ", err)
		return nil
	}
	return protocol.DecodeData(data)
}

func (conn *Connection) SendData(data *protocol.Data) (int, error) {
	pack := protocol.EncodeData(data)
	return conn.Send(pack)
}

func (conn *Connection) SendMessageto(msg string, receiver string) (int, error) {
	var data protocol.Data
	data.Type = protocol.Normal
	data.Time = time.Now().Format("15:04:05")
	data.User = conn.User
	data.Receiver.ID = receiver
	data.Content = msg
	return conn.SendData(&data)
}

func (conn *Connection) SendError(err string) (int, error) {
	var data protocol.Data
	data.Type = protocol.Error
	data.Time = time.Now().Format("15:04:05")
	data.Content = err
	return conn.SendData(&data)
}

func (conn *Connection) SendErrorf(format string, args ...interface{}) (int, error) {
	err := fmt.Sprintf(format, args...)
	return conn.SendError(err)
}

func (conn *Connection) SendOther(other string) (int, error) {
	var data protocol.Data
	data.Type = protocol.Other
	data.Time = time.Now().Format("15:04:05")
	data.User = conn.User
	data.Content = other
	return conn.SendData(&data)
}

func (conn *Connection) SendOtherto(other string, receiver string) {
	var data protocol.Data
	data.Type = protocol.Other
	data.Time = time.Now().Format("15:04:05")
	data.User = conn.User
	data.Receiver.ID = receiver
	data.Content = other
	conn.SendData(&data)
}
