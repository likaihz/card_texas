 package room

import (
	"../battle"
	"../card"
	"../lib/xx"
	"../lib/xxio"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
)

const SEATS = 10

type client chan<- string

type Room struct {
	sync.Mutex
	players []*Player
}

func New() *Room {
	r := &Room{}
	r.players = make([]*Player, SEATS)
	return r
}

// interface
func (r *Room) Enter(uid string, cli chan string) {
	r.Lock()
	defer r.Unlock()
	for i, p := range r.players {
		if p == nil {
			r.players[i] = NewPlayer(i, uid, cli)
			r.sendseat()
			return
		}
	}
}

func (r *Room) Leave(uid string) bool {
	r.Lock()
	defer r.Unlock()
	for i, p := range r.players {
		if p.Is(uid) {
			r.players[i] = nil
			r.sendseat()
			break
		}
	}
	return r.occupancy() == 0
}

func (r *Room) Ready(uid string) bool {
	r.Lock()
	defer r.Unlock()
	ready := true
	for i, p := range r.players {
		if p == nil {
			continue
		}
		if p.Is(uid) {
			p.SetReady(true)
			r.sendready(i)
		}
		ready = ready && p.Ready()
	}
	return ready && r.occupancy() > 1
}

func (r *Room) Newround(round int, dealer *Player) []*Player {
	r.Lock()
	defer r.Unlock()
	currents := r.currents(dealer)
	idx := dealer.Idx()
	r.sendstart(idx, round)
	return currents
}

func (r *Room) Player(idx int) *Player {
	if idx < 0 || idx >= SEATS {
		return nil
	}
	return r.players[idx]
}

//对应徳扑中的下注动作，需要修改 ...
func (r *Room) Setraise(player *Player, msg map[string]interface{}) {
	idx := msg["idx"].(float64)
	val := msg["raise"].(float64)
	num := int(val)
	player.SetRaise(num)
	r.sendraise(int(idx), num)
}

func (r *Room) Enterable() bool {
	return r.occupancy() < SEATS
}

func (r *Room) Next(player *Player) *Player {
	idx := -1
	if player != nil {
		idx = player.Idx()
	}
	for i := 0; i < SEATS; i++ {
		idx++
		if idx >= SEATS {
			idx = 0
		}
		p := r.players[idx]
		if p.Ready() {
			return p
		}
	}
	return nil
}

// implementation
func (r *Room) currents(dealer *Player) []*Player {
	arr := make([]*Player, 0)
	arr = append(arr, dealer)
	dealer.Reset()
	for p := r.Next(dealer); p != dealer; p = r.Next(p) {
		arr = append(arr, p)
	}
	return arr
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

// interface of send
func (r *Room) Msg() map[string]interface{} {
	msg := map[string]interface{}{}
	for i, p := range r.players {
		if p != nil {
			j := strconv.Itoa(i)
			msg[j] = p.Msg()
		}
	}
	return msg
}

func (r *Room) Send(msg map[string]interface{}) {
	for _, p := range r.players {
		if p != nil {
			p.Send(msg)
		}
	}
}

func (r *Room) Sendresult() {
	msg := map[string]interface{}{
		"opt": "result",
	}
	msg["data"] = r.Msg()
	r.Send(msg)
}

//广播所有人的下注 ...
func (r *Room) Sendstakes() {
	msg := map[string]interface{
		"opt": "stakes",
	}

	smsg := map[string]interface{}{}
	for i, p:=range r.players {
		if p != nil {
			j := strconv.Itoa(i)
			msg[j] = p.stakes
		}
	}
	msg["data"] = smsg
	r.Send(msg)
}

//broadcast the board to all players
func (r *Room) Sendboard(board []*card.Card) {
	msg := map[string]interface{}{}
	msg["opt"] = "board"
	cmsg := map[string]interface{}{}
	for i, c := range board {
		cmsg[strconv.Itoa(i)] = c.Msg()
	}
	msg["data"] = cmsg
	r.Send(msg, "all")
}

func (r *Room) Sendactions(uid string, act string, data int) {
	msg := map[string]interface{}{}
	msg["opt"] = "actions"
	amsg := map[string]interface{}{}
	amsg["uid"] = uid
	amsg["act"] = act
	amsg["data"] = strconv.Itoa(data)
	msg["data"] = amsg

	r.Send(msg, "all")
}

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
