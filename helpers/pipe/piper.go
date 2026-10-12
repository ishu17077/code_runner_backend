package pipe

import (
	"bytes"

	"github.com/gorilla/websocket"
	"k8s.io/client-go/tools/remotecommand"
)

type TerminalSession struct {
	Ws             *websocket.Conn
	SizeChan       chan remotecommand.TerminalSize
	InitialInput   *bytes.Buffer
	DoneChan       chan struct{}
	initialDataLen int
}

func (t *TerminalSession) Read(p []byte) (int, error) {
	if t.InitialInput != nil && t.InitialInput.Len() > 0 {
		//? also drains the buffer to not execute next time
		t.initialDataLen = t.InitialInput.Len()
		n, _ := t.InitialInput.Read(p)
		return n, nil
	}

	_, message, err := t.Ws.ReadMessage()
	if err != nil {
		select {
		case t.DoneChan <- struct{}{}:
		default:
		}
		return 0, err
	}
	n := copy(p, message)

	return n, nil
}

func (t *TerminalSession) Write(p []byte) (int, error) {
	if t.initialDataLen > 0 {
		n := len(p)
		if n <= t.initialDataLen {
			t.initialDataLen = t.initialDataLen - n
			return n, nil
		}
	}
	err := t.Ws.WriteMessage(websocket.TextMessage, p)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func (t *TerminalSession) Next() *remotecommand.TerminalSize {
	size := <-t.SizeChan
	return &size
}
