package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  128,
	WriteBufferSize: 128,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

const timeout = 10 * time.Second

// Starts a ping-pong server listening on the provided address. If cert and
// key are not passed, TLS will not be enabled.
func Listen(address string, port int, cert, key string) {
	http.HandleFunc("/", pingpong)

	var err error
	addr := fmt.Sprintf("%s:%d", address, port)
	if cert == "" {
		err = http.ListenAndServe(addr, nil)
	} else {
		err = http.ListenAndServeTLS(addr, cert, key, nil)
	}

	panic(err)
}

// Starts a ping-ping cycle with the websocket client. This is basically an
// echo server. The client is allowed to send five packets for 10
// seconds, being disconnected after whichever comes first.
func pingpong(w http.ResponseWriter, req *http.Request) {
	c, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		return
	}
	defer c.Close()

	for i := 0; i < 5; i++ {
		mt, message, err := c.ReadMessage()
		if err != nil {
			return
		}

		if err := c.WriteMessage(mt, message); err != nil {
			return
		}
	}
}
