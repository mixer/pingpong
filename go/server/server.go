package server

import (
	"crypto/tls"
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
		cfg := &tls.Config{
			MinVersion:               tls.VersionTLS10,
			CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
			PreferServerCipherSuites: true,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,
			},
		}
		srv := &http.Server{
			Addr:         addr,
			Handler:      nil,
			TLSConfig:    cfg,
			TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
		}
		err = srv.ListenAndServeTLS(cert, key)
	}

	panic(err)
}

// Starts a ping-ping cycle with the websocket client. This is basically an
// echo server. The client is allowed to send five packets for 10
// seconds, being disconnected after whichever comes first.
func pingpong(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Strict-Transport-Security", "max-age=31536000")

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
