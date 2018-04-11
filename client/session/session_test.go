package session

import (
	"testing"
	"time"

	"go-chat/proto"
)

func TestPrintMessageBlcok(t *testing.T) {
	var tests = []struct {
		msg string
	}{
		{"hello"},
		{"123456789012345"},
		{"123456789012345123456789012345"},
		{"你好"},
		{"一二三四五六七"},
		{"一二三四五六七八九十一二三四"},
		{"d[sfl]w[e]f][lew],/.,#@@%#^%^(&)*(_([lf]"},
		{`^*）&（#*……#（@&——&）！@你&￥#*。——#@￥@！，#￥#@”}：“|*^Q($&@"`},
	}

	sess := NewSession("test")
	for _, t := range tests {
		tm := time.Now().Format("15:04:05")
		msg := &proto.Message{Time: tm, Content: t.msg}
		sess.printMessageBlock(msg, left)
		sess.printMessageBlock(msg, right)
	}
}
