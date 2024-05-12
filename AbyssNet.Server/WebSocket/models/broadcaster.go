package models

import "github.com/lxzan/gws"

func Broadcast(conns []*gws.Conn, opcode gws.Opcode, payload []byte) {
	var b = gws.NewBroadcaster(opcode, payload)
	defer b.Close()
	for _, item := range conns {
		_ = b.Broadcast(item)
	}
}
