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

type Room struct {
	sync.Mutex
	receiving chan map[string]interface{}
	sending   chan map[string]interface{}
	clients   map[string]client
	started   bool
	number    int
	seat      []map[string]interface{}
	config    map[string]interface{}
	// testing
	log  *log.Logger
	file *os.File
}

func NewRoom(number int) *Room {
	r := &Room{number: number}
	r.receiving = make(chan map[string]interface{})
	r.sending = make(chan map[string]interface{})
	r.clients = map[string]client{}
	n := r.seatnum()
	r.seat = make([]map[string]interface{}, n)
	r.log, r.file = xx.Log2file("debug.log")
	go r.broadcast()
	return r
}

// interface
func (r *Room) Received() map[string]interface{} {
	return <-r.receiving
}

func (r *Room) Receiving(msg map[string]interface{}) {
	r.receiving <- msg
}

func (r *Room) Send(msg map[string]interface{}) {
	if len(msg) > 1 {
		r.sending <- msg
	}
}

func (r *Room) Sendone(uid string, msg map[string]interface{}) {
	cli := r.clients[uid]
	if len(msg) > 1 && cli != nil {
		cli <- xx.Map2str(msg)
	}
}

func (r *Room) Println(args ...interface{}) {
	r.log.Println(args...)
}

func (r *Room) Number() int {
	return r.number
}

func (r *Room) Occupancy() int {
	cnt := 0
	for _, usr := range r.seat {
		if usr != nil {
			cnt++
		}
	}
	return cnt
}

func (r *Room) Enterable() bool {
	if r == nil {
		return false
	}
	if r.started {
		return false
	}
	occu := r.Occupancy()
	return occu < r.seatnum()
}

func (r *Room) Seat() []map[string]interface{} {
	return r.seat
}

func (r *Room) Config() map[string]interface{} {
	return r.config
}

func (r *Room) Setconfig(data map[string]interface{}) {
	data["roomnum"] = r.number
	r.config = data
}

func (r *Room) Enter(uid string, cli client) {
	r.Lock()
	defer r.Unlock()
	for i, usr := range r.seat {
		if usr == nil {
			usr, err := getusr(uid)
			if err != nil {
				return
			}
			r.seat[i] = usr
			r.clients[uid] = cli
			r.Send(r.seatmsg())
			return
		}
	}
}

func (r *Room) Leave(uid string) {
	r.Lock()
	defer r.Unlock()
	for i, usr := range r.seat {
		if isusr(uid, usr) {
			r.seat[i] = nil
			delete(r.clients, uid)
			r.Send(r.seatmsg())
			return
		}
	}
}

func (r *Room) Ready(uid string) {
	r.Lock()
	defer r.Unlock()
	if r.Occupancy() <= 1 {
		return
	}
	ready := true
	for i, usr := range r.seat {
		if usr == nil {
			continue
		}
		if isusr(uid, usr) {
			r.seat[i]["ready"] = "ok"
			j := strconv.Itoa(i)
			msg := map[string]interface{}{
				"opt": "ready", "idx": j,
			}
			r.Send(msg)
		}
		ok := usr["ready"].(string)
		ready = ready && ok == "ok"
	}
	if ready {
		go r.start()
	}
}

// implementation
func (r *Room) start() {
	msg := map[string]interface{}{
		"opt": "start", "status": "ok",
	}
	r.Send(msg)
	r.started = true
	b := battle.New(r)
	for !b.Ended() {
		b.Round()
	}
	r.started = false
}

func (r *Room) broadcast() {
	for {
		msg := xx.Map2str(<-r.sending)
		r.Lock()
		for _, cli := range r.clients {
			cli <- msg
		}
		r.Unlock()
	}
}

func (r *Room) seatmsg() map[string]interface{} {
	msg := map[string]interface{}{"opt": "seat"}
	for i, usr := range r.seat {
		j := strconv.Itoa(i)
		msg[j] = usr
	}
	return msg
}

func (r *Room) seatnum() int {
	return 5
}

// database manipulation
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
	usr["ready"] = ""
	return usr, nil
}
