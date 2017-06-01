package room

import (
	"../card"
	"../lib/ws"
	"../lib/xx"
	"../lib/xxio"
	"fmt"
	"strconv"
	"sync"
)

// class Player
type Player struct {
	sync.Mutex
	// ready         bool
	conn          *ws.Conn
	Idx           int
	Uid, Name     string
	cards         *card.Cards
	chips, stakes int //chips是玩家拥有的筹码，stakes是下的赌注
	Status        string
}

func NewPlayer(idx int, uid string, conn *ws.Conn) *Player {
	p := &Player{Idx: idx, Uid: uid}
	name, err := getname(uid)
	if err != nil {
		return nil
	}
	p.Name = name
	// p.fold = false
	p.Connect(conn)
	p.Init()
	return p
}

func (p *Player) Debug() {
	fmt.Println("---- debug ----")
	p.conn.Receive()
}

func (p *Player) Init() {
	p.Status = "waiting"
	p.cards = card.NewCards()
	p.stakes = 0
}

func (p *Player) Connected() bool {
	return p.conn.Connected()
}

func (p *Player) Revive(conn *ws.Conn) {
	p.conn = conn
	p.Status = "revive"
}

// interface of game
func (p *Player) Is(uid string) bool {
	if p == nil {
		return false
	}
	return p.Uid == uid
}

func (p *Player) Present() bool {
	if p == nil {
		return false
	}
	return p.Status != "escape"
}

func (p *Player) Active() bool {
	if p == nil {
		return false
	}
	return p.Status == "active"
}

func (p *Player) Rank() int {
	return p.cards.Rank()
}

func (p *Player) Fold() {
	p.Status = "fold"

}

func (p *Player) Chips() int {
	return p.chips
}

func (p *Player) Stakes() int {
	return p.stakes
}

func (p *Player) Addstakes(s int) int {
	p.stakes += s
	return p.stakes
}

func (p *Player) All_in() int {
	c := p.chips
	p.chips -= c
	p.stakes += c
	return c
}

func (p *Player) Cards() *card.Cards {
	return p.cards
}

func (p *Player) Msg() map[string]interface{} {
	msg := map[string]interface{}{}
	msg["name"] = p.name
	msg["rank"] = p.cards.Rank()
	msg["cards"] = p.cards.Msg()
	msg["score"] = p.score
	return msg
}

func (p *Player) Send(msg map[string]interface{}) {
	if p == nil {
		return
	}
	str := xx.Map2str(msg)
	p.cli <- str
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

func (p *Player) Sendcard(c *card.Card) {
	msg := map[string]interface{}{}
	msg["opt"] = "card"
	msg["data"] = c.Msg()
	p.Send(msg)
}

func (p *Player) Sendcards() {
	msg := map[string]interface{}{}
	msg["opt"] = "cards"
	msg["data"] = p.cards.Msg()
	p.Send(msg)
}

func (p *Player) Stakeover() int {
	c := p.chips
	p.chips = 0
	return c

}

func (p *Player) Sendpermissions(act ...string) {
	//...
}

// implementation
func getname(uid string) (string, error) {
	cfg, err := xxio.Read("user")
	if err != nil {
		return "", err
	}
	val, ok := cfg[uid]
	if !ok {
		return "", fmt.Errorf("invalid uid!")
	}
	data := val.(map[string]interface{})
	name := data["name"].(string)
	return name, nil
}

// test
func (p *Player) Print() {
	s := "player "
	if p == nil {
		s += "none"
	} else {
		idx := strconv.Itoa(p.idx)
		s += idx + " : "
		s += p.name
	}
	fmt.Println(s)
}
