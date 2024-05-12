package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Gandalf-Le-Dev/abyssnet.server.websocket/models"
	"github.com/lxzan/gws"
)

var (
	port = flag.Int("port", 6666, "port to listen")
)

func main() {
	flag.Parse()
	log.Printf("Starting AbyssNet WebSocket server on port %d ...\n", *port)

	var wsh = models.NewWebSocketHandler()
	upgrader := gws.NewUpgrader(wsh, &gws.ServerOption{
		ParallelEnabled:   true,                                  // Parallel message processing
		Recovery:          gws.Recovery,                          // Exception recovery
		PermessageDeflate: gws.PermessageDeflate{Enabled: false}, // Enable compression
	})
	
	http.HandleFunc("/connect", func(writer http.ResponseWriter, request *http.Request) {
		log.Printf("Connection request from: %s\n", request.RemoteAddr)

		socket, err := upgrader.Upgrade(writer, request)
		if err != nil {
			return
		}

		go func() {
			defer socket.NetConn().Close()
			socket.ReadLoop() // Blocking prevents the context from being GC.
		}()

		go func() {
			for {
				// Broadcast message to all clients
				models.Broadcast(wsh.GetClients(), gws.OpcodeBinary, []byte("Hello, clients!"))

				// Wait for some time before broadcasting again
				time.Sleep(5 * time.Second)
			}
		}()

	})
	log.Panic(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
