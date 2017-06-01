package room

import (
	// "../card"
	"../lib/ws"
	"fmt"
	"strconv"
	"sync"
)

const SEATNUM = 10

// type client chan<- string

type Room struct {
	sync.Mutex
	players []*Player
	playing bool
	current int
}

func New() *Room {
	r := &Room{}
	r.players = make([]*Player, SEATNUM)
	return r
}

// interface
func (r *Room) Connect(uid string, conn *ws.Conn) {
	for _, p := range r.players {
		if p.Is(uid) && p.Status == "active" {
			p.Connect(conn)
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
			if r.playing && p.Active() {
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
	ok := true
	if r.playing {
		r.playing = false
		r.roundover()
		r.Send("over", r.Msg())
	}
	for i, p := range r.players {
		if p == nil {
			continue
		}
		if p.Is(uid) {
			p.Status = "active"
			r.Send("ready", map[string]interface{}{"idx": p.Idx})
		}
		ok = ok && p.Active()
	}
	ok = ok && r.occupancy() > 1
	if ok {
		r.playing = true
	}
	return ok
}

func (r *Room) Newround() []*Player {
	r.Lock()
	defer r.Unlock()
	r.current = r.Next(r.current)
	currents := r.currents()
	return currents
}

// //对应徳扑中的下注动作，需要修改 ...
// func (r *Room) Setraise(player *Player, msg map[string]interface{}) {
// 	idx := msg["idx"].(float64)
// 	val := msg["raise"].(float64)
// 	num := int(val)
// 	player.SetRaise(num)
// 	r.sendraise(int(idx), num)
// }

func (r *Room) Enterable() bool {
	return r.occupancy() < SEATNUM
}

func (r *Room) currents() []*Player {
	d := r.current
	arr := make([]*Player, 0)
	arr = append(arr, d)
	for p := r.Next(d.Idx); p != d; p = r.Next(p.Idx) {
		arr = append(arr, p)
	}
	return arr
}

func (r *Room) Next(idx int) *Player {
	for i := 0; i < SEATNUM; i++ {
		idx++
		if idx >= SEATNUM {
			idx = 0
		}
		p := r.players[idx]
		if p.Active() {
			return p
		}
	}
	return nil
}

// implementation

// interface of send
// func (r *Room) Msg() map[string]interface{} {
// 	msg := map[string]interface{}{}
// 	for i, p := range r.players {
// 		if p != nil {
// 			j := strconv.Itoa(i)
// 			msg[j] = p.Msg()
// 		}
// 	}
// 	return msg
// }

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
	if r.playing {
		data["playing"] = "ok"
	}
	data["dealer"] = r.current
	for i, p := range r.players {
		if p != nil {
			j := strconv.Itoa(i)
			data[j] = p.Seatmsg()
		}
	}
	return data
}

// func (r *Room) Sendresult() {
// 	msg := map[string]interface{}{
// 		"opt": "result",
// 	}
// 	msg["data"] = r.Msg()
// 	r.Send(msg)
// }

//广播所有人的下注 ...
// func (r *Room) Sendstakes() {
// 	msg := map[string]interface{
// 		"opt": "stakes",
// 	}

// 	smsg := map[string]interface{}{}
// 	for i, p:=range r.players {
// 		if p != nil {
// 			j := strconv.Itoa(i)
// 			msg[j] = p.stakes
// 		}
// 	}
// 	msg["data"] = smsg
// 	r.Send(msg)
// }

//broadcast the board to all players
// func (r *Room) Sendboard(board []*card.Card) {
// 	msg := map[string]interface{}{}
// 	msg["opt"] = "board"
// 	cmsg := map[string]interface{}{}
// 	for i, c := range board {
// 		cmsg[strconv.Itoa(i)] = c.Msg()
// 	}
// 	msg["data"] = cmsg
// 	r.Send(msg, "all")
// }

// func (r *Room) Sendactions(uid string, act string, data int) {
// 	msg := map[string]interface{}{}
// 	msg["opt"] = "actions"
// 	amsg := map[string]interface{}{}
// 	amsg["uid"] = uid
// 	amsg["act"] = act
// 	amsg["data"] = strconv.Itoa(data)
// 	msg["data"] = amsg

// 	r.Send(msg, "all")
// }

func (r *Room) sendseat() {
	msg := map[string]interface{}{
		"opt": "seat",
	}
	data := r.Msg()
	if len(data) == 0 {
		return
	}
	msg["data"] = data
	r.Send(msg)
}

func (r *Room) sendready(i int) {
	msg := map[string]interface{}{
		"opt": "ready",
	}
	msg["idx"] = strconv.Itoa(i)
	r.Send(msg)
}

func (r *Room) sendstart(dealer, round int) {
	msg := map[string]interface{}{
		"opt": "start",
	}
	msg["dealer"] = dealer
	msg["round"] = round
	r.Send(msg)
}

//待修改 ...
func (r *Room) sendraise(idx, raise int) {
	msg := map[string]interface{}{
		"opt": "raise",
	}
	msg["idx"] = idx
	msg["raise"] = raise
	r.Send(msg)
}

// implementation
func (r *Room) roundover() {
	for i, p := range r.players {
		if p != nil {
			if p.Status == "escape" || !p.Connected() {
				r.players[i] = nil
			} else {
				p.Init()
			}
		}
	}
}

func (r *Room) occupancy() int {
	cnt := 0
	for _, p := range r.players {
		if p != nil {
			cnt++
		}
	}
	return cnt
}
