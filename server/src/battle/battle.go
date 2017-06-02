package battle

import (
	"../card"
	"../lib/ws"
	"../lib/xx"
	"../room"
	"fmt"
	"strconv"
)

var limit = 100 //注限

// class Battle
type Battle struct {
	match     *Match
	number    int
	config    map[string]interface{}
	receiving chan map[string]interface{}
	room      *room.Room
	cards     []*card.Card
	pot       int 	//奖池大小
	board     []*card.Card 	//公共牌
	roundnum  int
	roundcnt  int
	
	mode      string
	ended     bool
	
}

func New(match *Match, number int) *Battle {
	b := &Battle{
		match: match, number: number,
	}
	b.room = room.New()
	b.receiving = make(chan map[string]interface{})
	b.pot = 0
	b.board = nil
	return b
}

func (b *Battle) Number() int {
	if b == nil {
		return -1
	}
	return b.number
}

func (b *Battle) Connect(uid string, conn *ws.Conn) {
	b.room.Connect(uid, conn)
}

func (b *Battle) Received(msg map[string]interface{}, conn *ws.Conn) (bool, *Battle) {
	if b == nil || b.ended {
		return true, nil
	}
	ok, opt := xx.Getstring(msg, "opt")
	if !ok {
		b.end()
		return false, nil
	}
	uid := msg["uid"].(string)
	switch opt {
	case "enter":
		b.room.Enter(uid, cli)
	case "leave":
		nobody := b.room.Leave(uid)
		if nobody {
			b.end()
		}
	case "ready":
		ok := b.room.Ready(uid)
		if ok {
			go b.round()
		}
	default:
		b.receiving <- msg
	}
	return true, b
}

func (b *Battle) receive(option string) map[string]interface{} {
	for {
		msg := <-b.receiving
		ok, opt := xx.Getstring(msg, "opt")
		if !ok {
			continue
		}
		if opt == option {
			return msg
		}
	}
	return nil
}

func (b *Battle) Enterable() bool {
	return b.room.Enterable()
}

func (b *Battle) Getconfig() map[string]interface{} {
	b.config["roundcnt"] = b.roundcnt
	return b.config
}

func (b *Battle) Setconfig(data map[string]interface{}) bool {
	if !checkcfg(data) {
		return false
	}
	data["roomnum"] = b.number
	b.config = data
	b.mode = data["mode"].(string)
	// b.turning = data["turning"].(string)
	b.roundnum = data["roundnum"].(int)
	return true
}

func (b *Battle) Ended() boo l {
	if b == nil {
		return false
	}
	return b.ended
}

func (b *Battle) Msg(currents []*room.Player) map[string]interface{} {
	data := map[string]interface{}{}
	for _, p := range currents {
		j := strconv.Itoa(p.Idx)
		data[j] = p.Cardmsg()
	}
	return data
}

// implementation
func (b *Battle) round() {
	b.roundcnt++
	b.cards = card.Init()
	currents := b.room.Newround()
	d := currents[0]
	b.pot = 0
	b.board = []*card.Card{} //clear the board

	data := map[string]interface{}{
		"dealer": d.Idx, "round": b.roundcnt,
	}
	b.room.Sendactive("start", data)

	//blinds
	sb := currents[1]
	bb := currents[2]
	sb.Addstakes(limit / 2)
	bb.Addstakes(limit)

	data = map[string]interface{}{
		"bigblind": map[string]interface{}{"idx": bb.Idx, "num": limit},
		"smallblind": map[string]interface{}{"idx": sb.Idx, "num": limit/2}
	}
	b.room.Sendactive("blinds", data)
	
	//deal cards
	for _, p := range b.players {
		p.Init()
		b.deal(p, 2)
		p.Sendactive("cards", p.Cardmsg())
	}

	//pre-flop
	ended := b.betround("preflop", currents, limit)
	if ended == 0 {
		//flop round
		for i := 0; i < 3; i++ {
			board[i] = b.popcard()
		}
		b.sendboard()
		ended = b.betround("flop", currents, 0)
		if ended == 0 {
			//turn
			board[3] = b.popcard()
			b.sendboard()
			ended = b.betround("turn", currents, 0)
			if ended == 0 {
				//river
				board[4] = b.popcard()
				b.sendboard()
				ended = b.betround("river", currents, 0)
			}
		}
	}

	switch ended {
	case 1:
		//翻出所有公共牌并比大小 ...
	case 2:
		//奖池中的全部筹码归唯一没有弃牌的玩家 ...
	}

	if b.round >= b.roundnum {
		b.End()
	}
}

func (b *Battle) compare(p1, p2 *Player) {
	if p1.Compare(p2) {
		b.addscore(p1, p2)
	} else {
		b.addscore(p2, p1)
	}
}

func (b *Battle) popcard() *card.Card {
	return card.Pop(b.cards)
}

func (b *Battle) deal(player *Player, num int) {
	for i := 0; i < num; i++ {
		c := b.popcard()
		player.Obtain(c)
	}
	player.Sendcards()
}

func (b *Battle) betround(round string, currents []*room.Player) {
	//用返回值表示下注回合是否结束
	//head参数并不是真正的最先行动的人，要首先确定其没有弃牌

	//首先找到真正的head
	var head, crt, min int
	if round == "preflop" {
		head = 2
		crt = next(head, currents)
		min = limit
	}
	else {
		crt = next(0, currents)
		head = -1
		min = 0
	}
	for ; ; crt = next(crt, currents) {
		//首先确定当前玩家能够进行的动作
		s := currents[crt].Stakes
		c := currents[crt].Chips
		var data map[string]interface{}
		if crt == head {
			//一圈下注已经结束，回到第一个下(加)注的人，那么该轮下注结束
			break
		} else if s >= min {
			//当前玩家能够选择看牌、加注、弃牌、全下
			data = permsg("check", "raise", "allin", "fold")
		} else if min-s >= c {
			//可以弃牌、全下
			data = permsg("allin", "fold")
		} else {
			//可以跟注、加注、弃牌、全下
			data = permsg("call", "fold", "raise", "allin")
		}
		currents[crt].Sendactive("permisssions", data)
		//接收玩家的身份信息uid操作信息act
		msg := b.receive("actions")
		uid := msg["uid"].(string)
		for uid != currents[crt].Uid {
			msg = b.receive("actions")
			uid = msg["uid"].(string)
		}
		act := msg["act"].(string)
		fnum := msg["num"].(float64)
		num := int(fdata)
		switch act {
		case "call":
			//跟注
			currents[crt].Call(min)
			num = min
		case "raise":
			//加注，更新head
			currents[crt].Raise(sdata, min)
			sdata += min
			head = crt
		case "fold":
			//弃牌
			currents[crt].Fold()
			b.Addpot(s)
		case "allin":
			//全下
			currents[crt].All_in()
			data = c
		case "check":
			//看牌,好像没什么做的？
		}
		//将当前玩家的动作广播给所有玩家
		b.sendactions(uid, act, data)
	}
	//将所有人的下注加到奖池中，并广播该轮下注结束后奖池数量以及每个玩家剩下的筹码
	chips := map[string]interface{}{}
	for _, p := range currents {
		if p.Active() && !p.Folded() {
			b.Addpot(p.Stakes)
			p.Stakes = 0
			chips[p.Uid] = p.Chips
		}
	}
	data = map[string]interface{}{"pot": pot, "chips": chips}
	b.room.Sendactive("pot", data)
	//判断之后是否还需要继续下一轮下注
	//结束的情况有两种：所有人都all in或者只剩下一个人没有弃牌
	all_allin := true
	not_fold := 0
	for _, p := range currents {
		if p.Active() {
			if !p.Allin() {
				all_allin = false
			}
			if !p.Folded() {
				not_fold ++
			}
		}
	}
	
	if all_allin || round == "river"{
		return 1
	}
	if not_fold == 1 {
		return 2
	}
	return 0
}


func (b *Battle) Addpot(n int) {
	b.pot += n
}

func (b *Battle) addstakes(p *room.Player, n int) {

}

func (b *Battle) sendactions(uid string, act string, num int) {
	data := map[string]interface{}{}
	data["uid"] = uid
	data["action"] = act
	data["num"] = strconv.Itoa(num)
	b.room.Sendactive("actions", data)
}

func (b *Battle) sendboard() {
	data := map[string]interface{}{}
	if len(board) != 0{
		for i, e := range board {
			s := strconv.Itoa(i)
			data[s] = e.Msg()
		}
	}
	b.room.Sendactive("board", data)
}

func next(c int, currents []*room.Player) int {
	num = len(currents)
	next := -1
	for i := c + 1; next < 0; i++ {
		i = i % num
		if currents[i].Status == "active" && !currents[i].Folded() {
			next = i
		}
	}
	return next
}

func permsg(p ...string) map[string]interface{} {
	data := map[string]interface{}{}
	for i, s := range p {
		j := strconv.Itoa(i)
		data[j] = s
	}
	return data
}