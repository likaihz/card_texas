package match

import (
	"../battle"
	"../lib/xx"
	"../lib/xxio"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
)

type client chan<- string

type Group struct {
	sync.Mutex
	entering  chan client
	leaving   chan client
	receiving chan map[string]interface{}
	sending   chan string
	fighting  bool
	idx       int
	seat      []map[string]interface{}
	config    map[string]interface{}
	// testing
	log  *log.Logger
	file *os.File
}

func NewGroup(uid string, ch chan string, idx int) *Group {
	g := &Group{idx: idx}
	g.entering = make(chan client)
	g.leaving = make(chan client)
	g.receiving = make(chan map[string]interface{})
	g.sending = make(chan string)
	num := g.seatnum()
	g.seat = make([]map[string]interface{}, num)
	g.log, g.file = xx.Log2file("debug.log")
	// initialize
	go g.broadcast()
	g.Enter(uid, ch)
	return g
}

// interface
func (g *Group) Println(args ...interface{}) {
	g.log.Println(args...)
}

func (g *Group) Index() int {
	return g.idx
}

func (g *Group) Occupancy() int {
	cnt := 0
	for _, usr := range g.seat {
		if usr != nil {
			cnt++
		}
	}
	return cnt
}

func (g *Group) SetConfig(msg map[string]interface{}) {
	g.config = msg["config"].(map[string]interface{})
}

func (g *Group) Receive(msg map[string]interface{}) {
	if g.fighting {
		g.receiving <- msg
	}
}

func (g *Group) Send(msg map[string]interface{}) {
	if len(msg) > 1 {
		g.sending <- xx.Map2str(msg)
	}
}

func (g *Group) Enterable() bool {
	if g == nil {
		return false
	}
	occu := g.Occupancy()
	return occu < g.seatnum()
}

func (g *Group) Enter(uid string, cli client) {
	g.Lock()
	defer g.Unlock()
	for i, usr := range g.seat {
		if usr == nil {
			usr, err := getusr(uid)
			if err != nil {
				return
			}
			usr["sending"] = cli
			usr["receiving"] = g.receiving
			g.seat[i] = usr
			g.entering <- cli
			break
		}
	}
	if g.Occupancy() == g.seatnum() {
		go g.fight()
	}
}

func (g *Group) Leave(uid string, cli client) {
	g.Lock()
	defer g.Unlock()
	for i, usr := range g.seat {
		if isusr(uid, usr) {
			g.seat[i] = nil
			g.leaving <- cli
			return
		}
	}
}

// implementation
func (g *Group) fight() {
	g.fighting = true
	b := battle.New(g.seat)
	b.SetConfig(g.config)
	for !b.Ended() {
		b.Round()
	}
	g.fighting = false
}

func (g *Group) broadcast() {
	var mu sync.Mutex
	clients := map[client]int{}
	go func() {
		for {
			msg := <-g.sending
			mu.Lock()
			for cli := range clients {
				cli <- msg
			}
			mu.Unlock()
		}
	}()

	for {
		select {
		case cli := <-g.entering:
			mu.Lock()
			clients[cli] = 0
			mu.Unlock()
		case cli := <-g.leaving:
			mu.Lock()
			delete(clients, cli)
			mu.Unlock()
		}
		msg := g.getseat()
		g.Send(msg)
	}
}

func (g *Group) getseat() map[string]interface{} {
	g.Lock()
	defer g.Unlock()
	msg := map[string]interface{}{"opt": "seat"}
	for i, usr := range g.seat {
		j := strconv.Itoa(i)
		msg[j] = usr
	}
	return msg
}

func (g *Group) seatnum() int {
	return 5
}

func isusr(uid string, usr map[string]interface{}) bool {
	if usr == nil {
		return false
	}
	x := usr["uid"].(string)
	return x == uid
}

func getusr(uid string) (map[string]interface{}, error) {
	cfg, err := xxio.Read("user")
	if err != nil {
		return nil, err
	}
	val, ok := cfg[uid]
	if !ok {
		return nil, fmt.Errorf("invalid uid!")
	}
	data := val.(map[string]interface{})
	usr := map[string]interface{}{}
	usr["uid"] = data["uid"].(string)
	usr["name"] = data["name"].(string)
	return usr, nil
}
