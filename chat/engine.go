package chat

import (
	"io"
	//TEMP:
	"runtime"
	"fmt"
	"runtime/debug"
	"encoding/json"
)

type ReaderWriter interface {
	io.Reader
	io.Writer
}

type Message struct {
	content []byte
}

type Connection struct {
	chWrite chan *Message
	connection   *ReaderWriter
}

func (c *Connection) handleWrites(stopIt chan bool) {
	dst := *c.connection

loop:
	for {
		select {

		case msg := <-c.chWrite:
			_, _ = dst.Write(msg.content)

		case <-stopIt:
			//close(c.chWrite) // CONSIDER USING THIS TO REPLACE STOPIT WITH RANGE
			//close(stopIt)
			break loop

		}
	}
}

func (c *Connection) serveConnection(chRead chan *Message) {
	stopIt := make(chan bool)
	go c.handleWrites(stopIt)
	defer func() { stopIt <- true }()

	msg := make([]byte, 5 * 1024 * 1024)//32)
	src := *c.connection

	for {
		nr, er := src.Read(msg)

		if nr > 0 {
			chRead <- &Message{msg[0:nr]}
		}
		if er == io.EOF {
			break
		}
	}


	//close(chRead)
}

type Chat struct {
	connections []Connection
	chRead      chan *Message
}

func New() *Chat {
	chat := Chat{make([]Connection, 0, 5), make(chan *Message)}
	go chat.RunChat()
	return &chat
}

func (chat *Chat) RunChat() {
	for {
		msg := <-chat.chRead

		for _, con := range chat.connections {
			con.chWrite <- msg
		}
	}
}

func (chat *Chat) addConnection(connection *ReaderWriter) Connection {
	conn := Connection{make(chan *Message), connection}
	chat.connections = append(chat.connections, conn)
	return conn
}

func (chat *Chat) removeConnection(toRemove *Connection) {
	for i, el := range chat.connections {
		if *toRemove == el {
			//chat.connections = append(chat.connections[:i], chat.connections[i+1:]...)

			// THE FOLLOWING IS UNNECESSARILY COMPLEX AND ONLY TO RULE OUT A MemLeak
			copy(chat.connections[i:], chat.connections[i+1:])
			chat.connections[len(chat.connections)-1] = Connection{}//nil // or the zero value of T
			chat.connections = chat.connections[:len(chat.connections)-1]

			// THE FOLLOWING IS UNNECESSARILY COMPLEX AND ONLY TO RULE OUT A MemLeak
			if cap(chat.connections) > 2*len(chat.connections) {
				fmt.Println("SHRINK")
				shrink := make([]Connection, len(chat.connections))
	    		copy(shrink, chat.connections)
	    		chat.connections = shrink
	    		debug.FreeOSMemory()
    		}

			break
		}
	}
	// DEBUG::
	fmt.Println("Num GoRoutines", runtime.NumGoroutine())
	fmt.Println("Len(connections)", len(chat.connections))
	fmt.Println("Cap(connections)", cap(chat.connections))
	GoRuntimeStats()
}

func GoRuntimeStats() {
	m := &runtime.MemStats{}

	fmt.Println("# goroutines: ", runtime.NumGoroutine())
	runtime.ReadMemStats(m)
	//fmt.Println("Memory Acquired: ", m.Sys)
	//fmt.Println("Memory Used    : ", m.Alloc)
	js, _ := json.Marshal(m)
	fmt.Println("JSON Memory = ", string(js))
	//runtime.GC()
}
