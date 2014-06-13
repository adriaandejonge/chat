package main

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"github.com/bmizerany/pat"
	"net/http"
	"github.com/adriaandejonge/chat/chat"

	// TEMP:
	_ "net/http/pprof"
	
)

func main() {
	mux := pat.New()
	
	mux.Get("/chat", websocket.Handler(chat.CreateHandler()))
	mux.Get("/", http.FileServer(http.Dir(".")))

	http.Handle("/", mux)

	fmt.Println("Started, serving at 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
