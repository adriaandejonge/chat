package chat

import (
	"github.com/adriaandejonge/chat/websocket"
)


func CreateHandler() websocket.Handler {
	chat := New()

	return func(ws *websocket.Conn) {
		var rw ReaderWriter
		rw = &(*ws)
		conn := chat.addConnection(&rw)		
		conn.serveConnection(chat.chRead) // REMOVE chat.chRead -> HIDE BEHIND API
		chat.removeConnection(&conn)
		// OF BOVENSTAANDE DRIE DE API IN??
	}
}