package ws

import (
	"../xx"
	"fmt"
	"net/websocket"
)

type Conn struct {
	conn      *websocket.Conn
	ch        chan string
	connected bool
	messages  []map[string]interface{}
}

func New(conn *websocket.Conn) *Conn {
	c := &Conn{}
	c.messages = []map[string]interface{}{}
	c.Open(conn)
	return c
}

// interface
func (c *Conn) Msg() []map[string]interface{} {
	return c.messages
}

func (c *Conn) Open(conn *websocket.Conn) {
	c.connected = true
	c.conn = conn
	c.ch = make(chan string)
	go c.sending()
}

func (c *Conn) Close() {
	if c.connected {
		c.ch <- "disconnect"
	}
	c.conn.Close()
}

func (c *Conn) Empty(send bool) {
	if c.connected && send {
		c.Send(nil)
		return
	}
	c.messages = []map[string]interface{}{}
}

func (c *Conn) Receive() (string, error) {
	var req string
	err := websocket.Message.Receive(c.conn, &req)
	if err != nil {
		fmt.Println(err)
	}
	return req, err
}

func (c *Conn) Send(msg map[string]interface{}) {
	if msg != nil {
		c.inqueue(msg)
	}
	if c.connected {
		msg = c.dequeue()
		for msg != nil {
			c.ch <- xx.Map2str(msg)
			msg = c.dequeue()
		}
	}
}

func (c *Conn) Reload(c1 *Conn) *Conn {
	if c1 != nil {
		messages := c1.Msg()
		c.messages = append(c.messages, messages...)
	}
	return c
}

// implementation
func (c *Conn) sending() {
	for str := range c.ch {
		if str == "disconnect" {
			break
		}
		err := websocket.Message.Send(c.conn, str)
		if err != nil {
			fmt.Println("connection error:")
			fmt.Println(err)
			break
		}
	}
	fmt.Println("connection closed...")
	c.connected = false
}

func (c *Conn) inqueue(msg map[string]interface{}) {
	c.messages = append(c.messages, msg)
}

func (c *Conn) dequeue() map[string]interface{} {
	if len(c.messages) == 0 {
		return nil
	}
	msg := c.messages[0]
	c.messages = c.messages[1:]
	return msg
}
