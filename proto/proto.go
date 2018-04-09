package proto

import (
	"strings"
)

type Data struct {
	Type     byte
	Sender   string
	Receiver string
	Message
}

type Message struct {
	Time    string
	Topic   string
	Content string
}

type User struct {
	ID    string
	Name  string
	Sex   string
	Birth string
}

const (
	Normal byte = iota
	Request
	Success
	Error
	Other
)

func DecodeData(b []byte) *Data {
	ds := strings.SplitN(string(b[1:]), "\r\n", 5)
	msg := Message{ds[2], ds[3], ds[4]}
	d := &Data{b[0], ds[0], ds[1], msg}
	return d
}

func EncodeData(d *Data) []byte {
	members := []string{d.Sender, d.Receiver,
		d.Message.Time, d.Message.Topic, d.Message.Content}
	b := []byte(strings.Join(members, "\r\n"))
	b = append([]byte{d.Type}, b...)
	return b
}

func EncodeUser(u *User) string {
	members := []string{u.ID, u.Name, u.Sex, u.Birth}
	s := strings.Join(members, "\r\n")
	return s
}

func DecodeUser(s string) *User {
	strs := strings.SplitN(s, "\r\n", 4)
	u := &User{strs[0], strs[1], strs[2], strs[3]}
	return u
}
