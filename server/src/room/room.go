package room

import (
	"../lib/ws"
	// "fmt"
	"strconv"
	"sync"
)

// const SEATNUM = 10

type Room struct {
	sync.Mutex
	players []*Player
	seatnum int
	current int
}

func New(num int) *Room {
	r := &Room{}
	r.seatnum = num
	r.players = make([]*Player, num)
	return r
}

// interface
func (r *Room) Relink(uid string, conn *ws.Conn, ptr int) bool {
	for _, p := range r.players {
		if p.Is(uid) {
			p.Connect(conn, ptr)
			return true
		}
	}
	return false
}

func (r *Room) Connected() bool {
	for _, p := range r.players {
		if p != nil {
			if p.Connected() {
				return true
			}
		}
	}
	return false
}

func (r *Room) Check(uid string) bool {
	for _, p := range r.players {
		if p.Is(uid) {
			return true
		}
	}
	return false
}

func (r *Room) Enter(uid string, conn *ws.Conn) bool {
	if r.Back(uid, conn) {
		return true
	}
	r.Lock()
	defer r.Unlock()
	for i, p := range r.players {
		if p == nil {
			p = NewPlayer(i, uid, conn)
			r.Send("enter", p.Seatmsg())
			r.players[p.Idx] = p
			p.Send("seat", r.Msg())
			return true
		}
	}
	return false
}

func (r *Room) Back(uid string, conn *ws.Conn) bool {
	r.Lock()
	defer r.Unlock()
	for _, p := range r.players {
		if p.Is(uid) {
			msg := map[string]interface{}{
				"opt": "seat", "data": r.Msg(),
			}
			conn.Send(msg)
			conn.Empty(false)
			p.Connect(conn, 0)
			return true
		}
	}
	return false
}

func (r *Room) Leave(uid string) bool {
	r.Lock()
	defer r.Unlock()
	for i, p := range r.players {
		if p.Is(uid) {
			r.players[i] = nil
			r.Send("leave", p.Seatmsg())
			break
		}
	}
	return r.occupancy() == 0
}

func (r *Room) Escape(uid string) bool {
	r.Lock()
	defer r.Unlock()
	for _, p := range r.players {
		if p.Is(uid) {
			p.Status = "escape"
			r.Sendactive("escape", p.Seatmsg())
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

func (r *Room) Getrecord() map[string]interface{} {
	tbl := map[string]interface{}{}
	for i, p := range r.players {
		if p != nil {
			s := strconv.Itoa(i)
			tbl[s] = map[string]interface{}{
				"name": p.Name,
			}
		}
	}
	return tbl
}

func (r *Room) Record(kind string, record map[string]interface{}) {
	for _, p := range r.players {
		if p != nil {
			p.Record(kind, record)
		}
	}
}

func (r *Room) Save(pth string, val interface{}) {
	for _, p := range r.players {
		if p != nil {
			p.Save(pth, val)
		}
	}
}

func (r *Room) Enterable() bool {
	return r.occupancy() < r.seatnum
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
	for i := 0; i < r.seatnum; i++ {
		idx++
		if idx >= r.seatnum {
			idx = 0
		}
		p := r.players[idx]
		if p.Active() {
			return idx
		}
	}
	return -1
}

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
