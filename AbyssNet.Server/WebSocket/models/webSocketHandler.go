package models

import (
	"log"
	"sync"
	"time"

	"github.com/Gandalf-Le-Dev/abyssnet.server.websocket/errors"
	"github.com/lxzan/gws"
)

const (
	PING_INTERVAL = 10 * time.Second
	PING_WAIT     = 10 * time.Second
)

type WebSocketHandler struct {
	clients *sync.Map
}

func NewWebSocketHandler() *WebSocketHandler {
	return &WebSocketHandler{
		clients: new(sync.Map),
	}
}

func (wsh *WebSocketHandler) AddClient(conn *gws.Conn) {
	wsh.clients.Store(conn.RemoteAddr().String(), conn)
}

func (wsh *WebSocketHandler) RemoveClient(conn *gws.Conn) {
	wsh.clients.Delete(conn.RemoteAddr().String())
}

func (wsh *WebSocketHandler) GetClient(addr string) (*gws.Conn, bool) {
	if conn, ok := wsh.clients.Load(addr); ok {
		return conn.(*gws.Conn), true
	}
	return nil, false
}

func (wsh *WebSocketHandler) GetClientCount() int {
	var count int
	wsh.clients.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	return count
}

func (wsh *WebSocketHandler) GetClients() []*gws.Conn {
	clients := make([]*gws.Conn, 0)
	wsh.clients.Range(func(key, value interface{}) bool {
		clients = append(clients, value.(*gws.Conn))
		return true
	})
	return clients
}

func (wsh *WebSocketHandler) OnOpen(socket *gws.Conn) {
	if conn, ok := wsh.GetClient(socket.RemoteAddr().String()); ok {
		conn.WriteClose(errors.NORMAL_CLOSURE, []byte("connection replaced"))
	}
	_ = socket.SetDeadline(time.Now().Add(PING_INTERVAL + PING_WAIT))
	wsh.clients.Store(socket.RemoteAddr().String(), socket)

	log.Printf("Connection established: %s\n", socket.RemoteAddr())
	Broadcast(wsh.GetClients(), gws.OpcodeBinary, []byte("new client connected!"))
}

func (wsh *WebSocketHandler) OnClose(socket *gws.Conn, err error) {
	if conn, ok := wsh.GetClient(socket.RemoteAddr().String()); ok {
		wsh.clients.Delete(conn.RemoteAddr().String())
	}
	log.Printf("Connection closed: %s. With error: %s\n", socket.RemoteAddr(), err)
}

func (wsh *WebSocketHandler) OnPing(socket *gws.Conn, payload []byte) {
	log.Printf("Received ping from %s\n", socket.RemoteAddr())
	_ = socket.SetDeadline(time.Now().Add(PING_INTERVAL + PING_WAIT))
	_ = socket.WritePong(nil)
}

func (wsh *WebSocketHandler) OnPong(socket *gws.Conn, payload []byte) {
	log.Printf("Received pong from %s\n", socket.RemoteAddr())
	_ = socket.SetDeadline(time.Now().Add(PING_INTERVAL + PING_WAIT))
	_ = socket.WriteMessage(gws.OpcodeBinary, []byte("heartbeat from server"))
}

func (wsh *WebSocketHandler) OnMessage(socket *gws.Conn, message *gws.Message) {
	switch message.Opcode {
	case gws.OpcodeBinary:
		log.Printf("Received binary message from %s: %s\n", socket.RemoteAddr(), message.Bytes())
		// echo
		defer message.Close()
		socket.WriteMessage(message.Opcode, message.Bytes())
	case gws.OpcodeText:
		log.Printf("Received text message from %s: %s\n", socket.RemoteAddr(), message.Bytes())
		// echo
		defer message.Close()
		socket.WriteMessage(message.Opcode, message.Bytes())
	}
}
