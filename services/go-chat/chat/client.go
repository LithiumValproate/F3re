package chat

import (
	"github.com/gorilla/websocket"
	"go-chat/participant"
	"log"
	"time"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

// Client 是一个将 participant (业务逻辑) 和 websocket 连接 (网络) 绑定的适配器
type Client struct {
	participant participant.Participant
	room        *Room
	conn        *websocket.Conn
	send        chan []byte // 带缓冲的 channel，用于发送消息
}

// NewClient 创建一个新的 Client 实例
func NewClient(p participant.Participant, room *Room, conn *websocket.Conn) *Client {
	return &Client{
		participant: p,
		room:        room,
		conn:        conn,
		send:        make(chan []byte, 256), // 256 是缓冲大小
	}
}

type clientMessage struct {
	client *Client
	data   []byte
}

// readPump 从 WebSocket 连接中读取消息，并将其交给 Room 处理
func (c *Client) readPump() {
	defer func() {
		c.room.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, rawMessage, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		// 将原始消息和发送者 client 包装起来，发送给 room 进行集中处理
		msg := clientMessage{c, rawMessage}
		c.room.incoming <- msg
	}
}

// writePump 将消息从 client.send channel 推送到 WebSocket 连接
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The room closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// 将 channel 中缓冲的其他消息也一并写入，提高效率
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
