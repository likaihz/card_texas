package room

import (
	"../lib/ws"
	// "fmt"
	"strconv"
	"sync"
)

const SEATNUM = 10

type Room struct {
	sync.Mutex
	players []*Player
	// playing bool
	current int
}

func New() *Room {
	r := &Room{}
	r.players = make([]*Player, SEATNUM)
	return r
}

// interface
func (r *Room) Connect(uid string, conn *ws.Conn, ptr int) {
	for _, p := range r.players {
		if p.Is(uid) && p.Status == "active" {
			p.Connect(conn, ptr)
			return
		}
	}
}

func (r *Room) Enter(uid string, conn *ws.Conn) {
	r.Lock()
	defer r.Unlock()
	for i, p := range r.players {
		if p.Is(uid) {
			p.Revive(conn)
			r.players[i] = nil
			r.Send("revive", p.Seatmsg())
			r.players[i] = p
			p.Send("seat", r.Msg())
			return
		}
	}
	for i, p := range r.players {
		if p == nil {
			p = NewPlayer(i, uid, conn)
			r.Send("enter", p.Seatmsg())
			r.players[p.Idx] = p
			p.Send("seat", r.Msg())
			return
		}
	}
}

func (r *Room) Leave(uid string) bool {
	r.Lock()
	defer r.Unlock()
	for i, p := range r.players {
		if p.Is(uid) {
			var opt string
			if p.Active() {
				p.Status = "escape"
				opt = "escape"
			} else {
				r.players[i] = nil
				opt = "leave"
			}
			r.Send(opt, p.Seatmsg())
			break
		}
	}
	return r.occupancy() == 0
}

func (r *Room) Ready(uid string) bool {
	r.Lock()
	defer r.Unlock()
	var ok bool
	for _, p := range r.players {
		if p.Is(uid) {
			p.Print()
			if p.Status == "active" {
				return false
			}
			p.Status = "active"
			r.Send("ready", map[string]interface{}{
				"idx": p.Idx,
			})
			ok = true
		}
	}
	if !ok || r.occupancy() <= 1 {
		return false
	}
	for _, p := range r.players {
		if p != nil {
			ok = ok && p.Active()
		}
	}
	return ok
}

func (r *Room) Autoready() bool {
	for _, p := range r.players {
		if p != nil {
			if p.Status == "active" {
				continue
			}
			p.Status = "active"
			r.Send("ready", map[string]interface{}{
				"idx": p.Idx,
			})
		}
	}
	return r.occupancy() > 1
}

func (r *Room) Newround() []*Player {
	r.Lock()
	defer r.Unlock()
	r.current = r.Next(r.current)
	currents := r.currents()
	return currents
}

func (r *Room) Init() {
	for _, p := range r.players {
		if p != nil {
			p.Init()
		}
	}
}

func (r *Room) Enterable() bool {
	return r.occupancy() < SEATNUM
}

func (r *Room) currents() []*Player {
	d := r.current
	arr := make([]*Player, 0)
	arr = append(arr, r.players[d])
	for p := r.Next(r.players[d].Idx); p != d; p = r.Next(r.players[p].Idx) {
		arr = append(arr, r.players[p])
	}
	return arr
}

func (r *Room) Next(idx int) int {
	for i := 0; i < SEATNUM; i++ {
		idx++
		if idx >= SEATNUM {
			idx = 0
		}
		p := r.players[idx]
		if p.Active() {
			return idx
		}
	}
	return -1
}

// implementation
func (r *Room) Send(opt string, data map[string]interface{}) {
	for _, p := range r.players {
		p.Send(opt, data)
	}
}

func (r *Room) Sendactive(opt string, data map[string]interface{}) {
	for _, p := range r.players {
		p.Sendactive(opt, data)
	}
}

func (r *Room) Msg() map[string]interface{} {
	data := map[string]interface{}{}
	data["dealer"] = r.current
	for i, p := range r.players {
		if p != nil {
			j := strconv.Itoa(i)
			data[j] = p.Seatmsg()
		}
	}
	return data
}

// //待修改 ...
// func (r *Room) sendraise(idx, raise int) {
// 	msg := map[string]interface{}{
// 		"opt": "raise",
// 	}
// 	msg["idx"] = idx
// 	msg["raise"] = raise
// 	r.Send(msg)
// }

// implementation
func (r *Room) occupancy() int {
	cnt := 0
	for _, p := range r.players {
		if p != nil {
			cnt++
		}
	}
	return cnt
}
