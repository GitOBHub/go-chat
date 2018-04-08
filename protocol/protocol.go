package protocol

import (
	"strings"
)

type Package struct {
	DataLen uint64
	Data
}

type Data struct {
	Type byte
	Time string
	User
	Message
}

type User struct {
	ID    string
	Name  string
	Sex   string
	Birth string
}

type Message struct {
	Receiver User
	Content  string
}

const (
	Normal byte = iota
	Error
	Other
)

const headLen = 8

/*func (data *User) Read(p []byte) (int, error) {
	s := []byte(data.ID + "\r\n" +
		data.Name + "\r\n" +
		data.Sex + "\r\n" +
		data.Birth)
	for i := range s {
		p[i] = s[i]
	}
	return len(s), io.EOF
}

func Bytes(r io.Reader) []byte {
	var buf bytes.Buffer
	buf.ReadFrom(r)
	return buf.Bytes()
}*/

func DecodeData(d []byte) *Data {
	var data Data
	data.Type = d[0]
	ds := strings.SplitN(string(d[1:]), "\r\n", 10)
	data.Time,
		data.ID, data.Name, data.Sex, data.Birth,
		data.Receiver.ID, data.Receiver.Name, data.Receiver.Sex, data.Receiver.Birth,
		data.Content =
		ds[0], ds[1], ds[2], ds[3], ds[4], ds[5], ds[6], ds[7], ds[8], ds[9]
	return &data
}

func EncodeData(data *Data) []byte {
	members := []string{data.Time,
		data.ID, data.Name, data.Sex, data.Birth,
		data.Receiver.ID, data.Receiver.Name, data.Receiver.Sex, data.Receiver.Birth,
		data.Content}
	d := []byte(strings.Join(members, "\r\n"))
	//	head := make([]byte, headLen)
	//	binary.PutUvarint(head, uint64(len(d)+1))
	//	head = append(head, data.Type)
	return append([]byte{data.Type}, d...)
}
