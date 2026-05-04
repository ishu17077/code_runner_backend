package pipe

import (
	"github.com/gorilla/websocket"
	"k8s.io/client-go/tools/remotecommand"
)

type TerminalSession struct {
	Ws       *websocket.Conn
	SizeChan chan remotecommand.TerminalSize
	DoneChan chan struct{}
}

func (t *TerminalSession) Read(p []byte) (int, error) {
	_, message, err := t.Ws.ReadMessage()
	if err != nil {
		return 0, err
	}
	n := copy(p, message)

	return n, nil
}

func (t *TerminalSession) Write(p []byte) (int, error) {
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
