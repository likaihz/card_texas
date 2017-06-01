package battle

import (
	"../card"
	"../lib/xx"
	"../room"
	"fmt"
)

var limit = 100 //注限

// class Battle
type Battle struct {
	match     *Match
	number    int
	config    map[string]interface{}
	receiving chan map[string]interface{}
	room      *room.Room
	current   *room.Player
	cards     []*card.Card
	roundnum  int
	roundcnt  int
	turning   string
	mode      string
	ended     bool
	pot       int
	board     []*card.Card
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

func (b *Battle) Received(msg map[string]interface{}, cli chan string) *Battle {
	if b == nil || b.ended {
		return nil
	}
	opt := msg["opt"].(string)
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
	return b
}

func (b *Battle) receive(option string) map[string]interface{} {
	for {
		msg := <-b.receiving
		val, ok := msg["opt"]
		if !ok {
			continue
		}
		opt, ok := val.(string)
		if ok && opt == option {
			return msg
		}
	}
	return nil
}

func (b *Battle) Enterable() bool {
	return b.room.Enterable()
}

func (b *Battle) Getconfig() map[string]interface{} {
	return b.config
}

func (b *Battle) Setconfig(data map[string]interface{}) bool {
	if !checkcfg(data) {
		return false
	}
	data["roomnum"] = b.number
	b.config = data
	b.mode = data["mode"].(string)
	b.turning = data["turning"].(string)
	b.roundnum = data["roundnum"].(int)
	return true
}

func (b *Battle) Ended() boo l {
	if b == nil {
		return false
	}
	return b.ended
}

// implementation
func (b *Battle) round() {
	b.roundcnt++
	b.cards = card.Init()
	b.current = b.dealer() 
	currents := b.room.Newround(b.roundcnt, b.current)
	b.pot = 0
	b.board = []*card.Card{} //clear the board

	//blinds
	sb := currents[1]
	bb := currents[2]
	sb.Addstakes(limit / 2)
	bb.Addstakes(limit)

	b.room.Sendstakes()
	//deal cards
	for _, p := range b.players {
		p.Init()
		b.deal(p, 2)
		p.Sendcards()
	}

	b.betround(currents, 2, limit)

	//flop round
	for i := 0; i < 3; i++ {
		board[i] = b.popcard()
	}
	b.room.Sendboard(board)
	b.betround(currents, 1, 0)

	//turn
	board[3] = b.popcard()
	b.room.Sendboard(board)
	b.betround(currents, 1, 0)

	//river
	board[4] = b.popcard()
	b.sendboard()
	b.betround(currents, 1, 0)

	b.room.Received()
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

func (b *Battle) addscore(p1, p2 *Player) {
	r := p1.Raise() * p2.Raise()
	cards := p1.Cards()
	s := cards.Price(b.mode) * r
	p1.Addscore(s)
	p2.Addscore(-s)
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

func (b *Battle) next(c int) int {
	b.room.Received()
	num := len(b.room.Seat())
	next := -1
	for i := c + 1; next < 0; i++ {
		i = i % num
		if b.players[i] != nil {
			next = i
		}
	}
	return next
}

func (b *Battle) dealer() *room.Player {
	old := b.current
	return b.room.Next(old)
}

func (b *Battle) betround(currents []*room.Player, head int, min int) int {
	//用返回值表示下注回合是否结束，结束的情况有两种：所有人都all in或者只剩下一个人没有弃牌
	//head参数并不是真正的最先行动的人，要首先确定其没有弃牌
	//...

	//首先找到真正的head
	for currents[head].Folded() {
		head++
	}
	crt := head
	//处理pre-flop round中的特殊情况,直接跳过大盲注位
	for currents[crt].fold() || currents[head].Stakes() == min && min != 0{
		crt++
	}
	for {
		//首先确定当前玩家能够进行的动作
		s := currents[crt].Stakes()
		c := currents[crt].Chips()
		if crt == head && min != 0 {
			//一圈下注已经结束，回到第一个下注的人
			//当前玩家可以选择看牌、加注、全下（假设，实际的规则还要复杂得多）...
			currents[crt].Sendpermissions("check", "raise", "allin")
		} else if s >= min {
			//当前玩家能够选择看牌、加注、弃牌、全下
			currents[crt].Sendpermissions("check", "raise", "allin", "fold")
		} else if min-s >= c {
			//可以弃牌、全下
			currents[crt].Sendpermissions("allin", "fold")
		} else {
			//可以跟注、加注、弃牌、全下
			currents[crt].Sendpermissions("call", "fold", "raise", "allin")
		}
		//接收玩家的身份信息uid操作信息act
		msg:=b.receive("actions")
		uid := msg["uid"].(string)
		for uid != currents[crt].Uid() {
			msg = b.receive("actions")
			uid = msg["uid"].(string)
		}
		act := msg["act"].(string)
		fdata := msg["data"].(float64)
		data := int(fdata)
		switch act {
		case "call":
			//跟注
			currents[crt].Call(min)
			data = min
		case "raise":
			//加注，更新head
			data = min+data
			currents[crt].Raise(data)
			head = crt
		case "fold":
			//弃牌，如果是head弃牌也要更新head
			currents[crt].Fold()
			b.Addpot(s)
			if crt == head {
				for currents[head].Fold() {
					head ++
				}
			}
		case "allin":
			//全下
			currents[crt].All_in()
			data = c
		case "check":
			//看牌,好像没什么做的？
		}

		//将当前玩家的动作广播给所有玩家
		b.room.Sendactions(uid, act, data)
		//判断该轮下注是否结束
		if act == "check" && crt == head {
			break
		}

		crt = b.next(crt)
	}
	//判断之后是否还需要继续下一轮下注 ...

	//将所有人的下注加到奖池中，并广播 ...
}


func (b *Battle) Addpot(n int) {
	b.pot += n
}
