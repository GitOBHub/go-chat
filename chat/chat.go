package chat

import (
	"encoding/binary"
	"fmt"
	//"io"
	"net"
	"sync"
	"time"

	"server/chat/protocol"
)

type Connection struct {
	net.Conn
	Number int
	protocol.User
	Mu sync.Mutex
}

func (conn *Connection) ReadData() *protocol.Data {
	var dataLen uint64
	if err := binary.Read(conn.Conn, binary.LittleEndian, &dataLen); err != nil {
		return nil
	}
	data := make([]byte, dataLen)
	fmt.Printf("ReadData: slice len: %d\n", len(data))
	//	if _, err := io.ReadFull(conn.Conn, data); err != nil {
	if _, err := conn.Conn.(*net.TCPConn).Read(data); err != nil {
		return nil
	}
	fmt.Println("ReadData: finish")
	return protocol.DecodeData(data)
}

func (conn *Connection) SendData(data *protocol.Data) (int, error) {
	pack := protocol.EncodeData(data)
	return conn.Conn.(*net.TCPConn).Write(pack)
}

func (conn *Connection) SendMessageto(msg string, receiver string) (int, error) {
	fmt.Printf("SendMessageto: len of msg: %d\n", len(msg))
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
