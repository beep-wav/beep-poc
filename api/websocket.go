package api

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/Courtcircuits/mitter-server/types"
	"github.com/gorilla/websocket"
)

type Owner struct {
	id   int
	name string
}

type Connection struct {
	id     int
	Conn   *websocket.Conn
	Hub    *Hub
	Owner  Owner
	authed bool
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var ErrSockReqInvalidFormat = errors.New("the message is not the good format, it must be TextMessage")

func Handler(w http.ResponseWriter, r *http.Request, h *Hub) error {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		return err
	}
	connection := h.AddConnection(conn)

	log.Println("new connection")

	connection.ReceiveMessages()

	h.RemoveConnection(connection.id)
	return nil
}

// send a message to the client
func (c *Connection) SendMessage(msg types.Message) error {
	content := []byte(msg.Content)

	if err := c.Conn.WriteMessage(websocket.TextMessage, content); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// receive a message from the client
func (c *Connection) ReceiveMessages() error {
	for {
		messageType, p, err := c.Conn.ReadMessage() //message is a JSON
		if err != nil {
			log.Println(err)
			return err
		}

		if messageType != websocket.TextMessage {
			return ErrSockReqInvalidFormat
		}

		msg := types.Message{
			ID:         1,
			Content:    string(p),
			Timestamp:  string(rune(time.Now().Unix())),
			Name_owner: c.Owner.name,
		}

		if err != nil {
			log.Println(err)
			return err
		}

		c.Hub.Broadcast(msg)

		if string(p) == "exit" {
			break
		}

	}
	return nil
}
