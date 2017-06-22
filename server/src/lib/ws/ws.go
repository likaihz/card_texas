package ws

import (
	"../xx"
	"fmt"
	"net/http"
	"net/websocket"
	"sync"
)

type Conn struct {
	sync.Mutex
	conn        *websocket.Conn
	ch          chan string
	connected   bool
	messages    []map[string]interface{}
	msgptr, idx int
	onclose     func()
}

func Listen(port, pth string, handle func(*websocket.Conn)) {
	http.Handle(pth, websocket.Handler(handle))
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

var COUNT int

func New(conn *websocket.Conn) *Conn {
	COUNT++
	c := &Conn{}
	c.idx = COUNT
	c.messages = []map[string]interface{}{}
	c.Open(conn)
	return c
}

// interface
func (c *Conn) Msg() []map[string]interface{} {
	return c.messages
}

func (c *Conn) Connected() bool {
	return c.connected
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
}

func (c *Conn) Onclose(f func()) {
	c.onclose = f
}

func (c *Conn) Empty(send bool) {
	if c.connected && send {
		c.Send(nil)
	} else {
		c.clqueue()
	}
}

func (c *Conn) Receive() (map[string]interface{}, string) {
	var req string
	err := websocket.Message.Receive(c.conn, &req)
	if err != nil {
		fmt.Println("ws receive error: ", err)
		return nil, ""
	}
	ok, msg, opt := parse(req)
	if !ok {
		return nil, ""
	}
	if opt == "heartbeat" {
		c.ch <- xx.Map2str(msg)
	}
	return msg, opt
}

func (c *Conn) Send(msg map[string]interface{}) {

	if msg != nil {
		c.inqueue(msg)
		fmt.Println(msg)
	}
	if !c.connected {
		return
	}
	for msg := c.dequeue(); msg != nil; msg = c.dequeue() {
		c.ch <- xx.Map2str(msg)
	}
}

func (c *Conn) Resend(ptr int) {
	c.rsqueue(ptr)
	c.Send(nil)
}

func (c *Conn) Replace(c1 *Conn) *Conn {
	c.Lock()
	defer c.Unlock()
	if c1 != nil {
		messages := c1.Msg()
		c.messages = append(messages, c.messages...)
		c1.Close()
	}
	return c
}

// implementation
func parse(req string) (bool, map[string]interface{}, string) {
	msg := xx.Str2map(req)
	if msg == nil {
		return false, nil, ""
	}
	ok, opt := xx.Getstring(msg, "opt")
	if !ok {
		return false, nil, ""
	}
	ok, _ = xx.Getstring(msg, "uid")
	if !ok {
		return false, nil, ""
	}
	return true, msg, opt
}

func (c *Conn) sending() {
	defer func() {
		c.conn.Close()
		c.connected = false
		// c.onclose()
		fmt.Println("connection closed...")
	}()
	for s := range c.ch {
		if s == "disconnect" {
			break
		}
		err := websocket.Message.Send(c.conn, s)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func (c *Conn) inqueue(msg map[string]interface{}) {
	c.Lock()
	defer c.Unlock()
	c.messages = append(c.messages, msg)
}

func (c *Conn) dequeue() map[string]interface{} {
	c.Lock()
	defer c.Unlock()
	l := len(c.messages)
	if l == 0 || l <= c.msgptr {
		return nil
	}
	msg := c.messages[c.msgptr]
	if c.idx >= 3 {
		opt := msg["opt"].(string)
		fmt.Println("msg ", c.msgptr, opt)
	}
	c.msgptr++
	return msg
}

func (c *Conn) clqueue() {
	c.Lock()
	defer c.Unlock()
	c.msgptr = 0
	c.messages = []map[string]interface{}{}
}

func (c *Conn) rsqueue(i int) {
	c.Lock()
	defer c.Unlock()
	c.msgptr = i
}
