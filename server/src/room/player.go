package room

import (
	"../card"
	// "../lib/user"
	"../lib/ws"
	"../lib/xx"
	// "fmt"
	"log"
	"strconv"
	"sync"
)

//玩家在进入房间时自动买入筹码的数目
const INITSCORE = 1000

// class Player
type Player struct {
	sync.Mutex
	conn                   *ws.Conn
	cards                  *card.Cards
	Idx                    int
	roundscore             int
	Uid, Name, Avt, Status string
	Chip, Score            int //chips是玩家拥有的筹码，stakes是下的赌注
	Action                 string
	Statistic              map[string]int
}

func NewPlayer(idx int, uid string, conn *ws.Conn) *Player {
	p := &Player{Idx: idx, Uid: uid}
	if !p.getinfo() {
		return nil
	}
	p.Statistic = map[string]int{}
	conn.Empty(false)
	p.Connect(conn, 0)
	p.Init()
	return p
}

func (p *Player) Debug() {
	log.Println("---- debug ----")
	p.conn.Receive()
}

func (p *Player) Init() {
	p.Status = "waiting"
	p.conn.Empty(false)
	p.cards = card.NewCards()
	//p.Stakes = 0
	//进入房间时系统自动为玩家买入筹码 ...
	p.Score = INITSCORE
}

func (p *Player) Connected() bool {
	return p.conn.Connected()
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
func (p *Player) Addchip(s int) int {
	p.Score += s
	return p.Score
}

//全下，返回值为全下的金额
func (p *Player) All_in() int {
	s := p.Score
	p.Score -= s
	p.Chip += s
	p.Status = "allin"
	return s
}

func (p *Player) Count() {
	if p.roundscore > 0 {
		p.Statistic["winnum"]++
	}
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
func (p *Player) Buyscore(n int) {
	p.Score += n
}

func (p *Player) Addscore(n int) {
	p.Score += n
}

func (p *Player) Call(n int) {
	p.Chip += n
}
func (p *Player) Raise(b, n int) {
	p.Chip += (b + n)
}

// interface of net
func (p *Player) Seatmsg() map[string]interface{} {
	data := map[string]interface{}{}
	data["idx"] = p.Idx
	data["uid"] = p.Uid
	data["name"] = p.Name
	data["chip"] = p.Chip
	data["score"] = p.Score
	data["status"] = p.Status
	return data
}

func (p *Player) Cardmsg() map[string]interface{} {
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

func (p *Player) Record(kind string, record map[string]interface{}) {
	// pth := "data.record." + kind
	// val, err := user.Search(p.Uid, pth)
	// if err != nil {
	// 	fmt.Println("Player Record(): search ", err)
	// 	return
	// }
	// arr, ok := val.([]interface{})
	// if !ok {
	// 	fmt.Println("Player Record(): reord is not an array!!")
	// 	return
	// }
	// arr = append(arr, record)
	// if len(arr) > 20 {
	// 	arr = arr[1:]
	// }
	// err = user.Upsert(p.Uid, pth, arr)
	// if err != nil {
	// 	fmt.Println("Player Record(): upsert ", err)
	// 	return
	// }
}

func (p *Player) Save(pth string, val interface{}) {
	// err := user.Upsert(p.Uid, pth, val)
	// if err != nil {
	// 	fmt.Println("Player Upsert(): failed!")
	// }
}

// implementation
func (p *Player) getinfo() bool {
	uid := p.Uid
	// val, err := user.Search(uid, "info")
	// if err != nil {
	// 	fmt.Println("getinfo err: ", err)
	// 	return false
	// }
	// info, ok := val.(map[string]interface{})
	// if !ok {
	// 	return false
	// }
	// ok, name := xx.Getstring(info, "nickname")
	// if !ok {
	// 	return false
	// }
	// ok, avt := xx.Getstring(info, "headimgurl")
	// if !ok {
	// 	return false
	// }
	// p.Name, p.Avt = name, avt
	name, ok := getname(uid)
	if ok {
		p.Name = name
	}
	return true
}

// test
func getname(uid string) (string, bool) {
	cfg, err := xx.Read("user")
	if err != nil {
		log.Println("getname() ", err)
		return "", false
	}
	ok, data := xx.Getmap(cfg, uid)
	if !ok {
		log.Println("getname() invalid data!")
		return "", false
	}
	ok, name := xx.Getstring(data, "name")
	if !ok {
		log.Println("getname() invalid name!")
		return "", false
	}
	return name, true
}

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
	log.Println(s)
}
