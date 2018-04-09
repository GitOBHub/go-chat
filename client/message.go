package main

import (
	"go-chat/client/session"
	"go-chat/proto"
)

func prepareMessage(data *proto.Data) {
	mu.Lock()
	sess, ok := sessions[data.Sender]
	if !ok {
		sess = session.NewSession(data.Sender)
		sessions[data.Sender] = sess
	}
	mu.Unlock()
	sess.PutMessage(&data.Message)
}
