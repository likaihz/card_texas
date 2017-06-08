package room

import (
	"../card"
	"../lib/ws"
	"../lib/xxio"
	"fmt"
	"strconv"
	"sync"
)

//玩家在进入房间时自动买入筹码的数目
const INITCHIPS = 1000

// class Player
type Player struct {
	sync.Mutex
	conn          *ws.Conn
	Idx           int
	Uid, Name     string
	cards         *card.Cards
	Chips, Stakes int //chips是玩家拥有的筹码，stakes是下的赌注
	Status        string
	Action        string
}

func NewPlayer(idx int, uid string, conn *ws.Conn) *Player {
	p := &Player{Idx: idx, Uid: uid}
	name, err := getname(uid)
	if err != nil {
		return nil
	}
	p.Name = name
	p.Connect(conn, 0)
	p.Init()
	return p
}

func (p *Player) Debug() {
	fmt.Println("---- debug ----")
	p.conn.Receive()
}

func (p *Player) Init() {
	p.Status = "waiting"
	p.conn.Empty(false)
	p.cards = card.NewCards()
	p.Stakes = 0
	//进入房间时系统自动为玩家买入筹码 ...
	p.Chips = INITCHIPS
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

func (p *Player) Folded() bool {
	return p.Action == "fold"
}

func (p *Player) Allin() bool {
	return p.Action == "allin"
}

func (p *Player) Rank() int {
	return p.cards.Rank()
}

func (p *Player) Fold() {
	p.Action = "fold"
}

//下注
func (p *Player) Addstakes(s int) int {
	p.Stakes += s
	return p.Stakes
}

//全下，返回值为全下的金额
func (p *Player) All_in() int {
	c := p.Chips
	p.Chips -= c
	p.Stakes += c
	p.Status = "allin"
	return c
}

func (p *Player) Cards() *card.Cards {
	return p.cards
}

func (p *Player) Obtain(c *card.Card) {
	p.cards.Append(c)
}

func (p *Player) Compare(player *Player) int {
	cards := player.Cards()
	return p.cards.Compare(cards)
}

// 买入筹码
func (p *Player) Addchips(n int) {
	p.Chips += n
}

func (p *Player) Call(n int) {
	p.Stakes += n
}
func (p *Player) Raise() {

}

// interface of net
func (p *Player) Seatmsg() map[string]interface{} {
	data := map[string]interface{}{}
	data["idx"] = p.Idx
	data["name"] = p.Name
	data["chips"] = p.Chips
	data["status"] = p.Status
	return data
}

func (p *Player) Cardmsg() map[string]interface{} {
	//data := map[string]interface{}{}
	cards := p.cards
	return cards.Msg()
}

func (p *Player) Connect(conn *ws.Conn, ptr int) {
	p.conn = conn.Replace(p.conn)
	p.conn.Resend(ptr)
}

func (p *Player) Disconnect() {
	p.conn.Close()
}

func (p *Player) Send(opt string, data map[string]interface{}) {
	if len(data) > 0 && p.Present() {
		msg := map[string]interface{}{
			"opt": opt, "data": data,
		}
		p.conn.Send(msg)
	}
}

func (p *Player) Sendactive(opt string, data map[string]interface{}) {
	if len(data) > 0 && p.Active() {
		msg := map[string]interface{}{
			"opt": opt, "data": data,
		}
		p.conn.Send(msg)
	}
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
		idx := strconv.Itoa(p.Idx)
		s += idx + " : "
		s += p.Name
		s += " " + p.Status
	}
	fmt.Println(s)
}
