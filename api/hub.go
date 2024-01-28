package api

import (
	"fmt"
	"log"

	"github.com/Courtcircuits/mitter-server/types"
	"github.com/gorilla/websocket"
)

type Hub struct {
	pool   map[int]Connection
	lastId int
}

func NewHub() *Hub {
	return &Hub{
		pool:   make(map[int]Connection),
		lastId: 0,
	}
}

func (h *Hub) AddConnection(conn *websocket.Conn) Connection {
	h.lastId++
	h.pool[h.lastId] = Connection{id: h.lastId, Conn: conn, Hub: h, Owner: Owner{
		id:   h.lastId,
		name: fmt.Sprintf("user%d", h.lastId),
	}, authed: false}
	return h.pool[h.lastId]
}

func (h *Hub) RemoveConnection(id int) {
	h.pool[h.lastId].Conn.Close() // close connection before removing in order to avoid memory leaks
	delete(h.pool, id)
}

func (h *Hub) Broadcast(msg types.Message) error {
	log.Printf("from hub : %q\n", msg)
	for _, conn := range h.pool {
		if conn.Owner.name != msg.Name_owner {
			conn.SendMessage(msg)
		}
	}
	return nil
}
