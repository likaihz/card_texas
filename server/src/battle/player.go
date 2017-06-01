package battle

import (
	"../card"
	"../mm"
	// "fmt"
	"strconv"
	"sync"
)

// class Player
type Player struct {
	sync.Mutex
	room          mm.Room
	idx           int
	uid           string
	name          string
	sending       chan string
	receiving     chan map[string]interface{}
	cards         *card.Cards
	chips, stakes int
	folded        bool
}

func NewPlayer(room mm.Room, idx int) *Player {
	p := &Player{room: room, idx: idx}
	seat := room.Seat()
	data := seat[idx]
	p.uid = data["uid"].(string)
	p.name = data["name"].(string)
	p.fold = false
	p.Init()
	return p
}

func (p *Player) Init() {
	p.cards = card.NewCards()
	p.stakes = 0
	p.folded = false
}

// interface of info
func (p *Player) Idx() int {
	return p.idx
}

func (p *Player) Name() string {
	return p.name
}

func (p *Player) Uid() string {
	return p.uid
}
func (p *Player) Rank() int {
	return p.cards.Rank()
}

func (p *Player) Chips() int {
	return p.chips
}

func (p *Player) Addstakes(s int) int {
	p.stakes += s
	return p.stakes
}

func (p *Player) Stakes() int {
	return p.stakes
}

func (p *Player) Folded() bool {
	return p.folded
}

func (p *Player) Cards() *card.Cards {
	return p.cards
}

// interface of action
func (p *Player) Obtain(c *card.Card) {
	p.cards.Insert(c)
}

func (p *Player) Compare(player *Player) int {
	cards := player.Cards()
	return p.cards.Compare(cards)
}

func (p *Player) Win(player *Player) {
	// ...
}

func (p *Player) Fold() {
	p.fold = true
}

func (p *Player) All_in() int {
	c := p.chips
	p.chips -= c
	p.stakes += c
	return c
}

func (p *Player) Sendcard(c *card.Card) {
	msg := map[string]interface{}{}
	msg["opt"] = "card"
	msg["data"] = c.Msg()
	p.send(msg)
}

func (p *Player) Sendcards() {
	msg := map[string]interface{}{}
	msg["opt"] = "cards"
	msg["data"] = p.cards.Msg()
	p.send(msg)
}

func (p *Player) Sendpermissions(per ...string) {
	msg := map[string]interface{}{}
	msg["opt"] = "permissions"
	pmsg := map[string]interface{}{}
	for i, s := range per {
		j := strconv.Itoa(i)
		pmsg[j] = s
	}
	msg["data"] = pmsg
	p.send(msg)
}

func (p *Player) Stakeover() int {
	c := p.chips
	p.chips = 0
	return c

}

// implementation
func (p *Player) send(msg map[string]interface{}) {
	p.room.Sendone(p.uid, msg)
}
