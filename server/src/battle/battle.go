package battle

import (
	"../card"
	// "../lib/user"
	"../lib/ws"
	"../lib/xx"
	"../lib/xxio"
	"../room"
	"fmt"
	// "sort"
	"strconv"
	"time"
)

const (
	LIMIT  = 100
	expire = 1800
)

// class Battle
type Battle struct {
	match              *Match
	kind, roomnum      string
	config             map[string]interface{}
	readying, betting  chan map[string]interface{}
	timer              *time.Timer
	room               *room.Room
	currents           []*room.Player
	cards              []*card.Card
	pot                int          //奖池大小
	board              []*card.Card //公共牌
	roundcnt, roundnum int
	began, ended       bool
	mode               string
	when               string
}

func New(match *Match, roomnum string) *Battle {
	b := &Battle{
		match: match, roomnum: roomnum,
	}
	b.kind = "holdem"
	b.room = room.New(10)
	b.readying = make(chan map[string]interface{})
	b.betting = make(chan map[string]interface{})
	b.pot = 0
	b.board = nil

	dura := time.Duration(expire)
	b.timer = time.NewTimer(dura * time.Second)
	go func() {
		<-b.timer.C
		b.timer.Stop()
		b.match.Remove(b.roomnum)
		print("battle expired...\n")
	}()
	return b
}

// interface
func (b *Battle) Kind() string {
	if b == nil {
		return ""
	}
	return b.kind
}

func (b *Battle) Roomnum() string {
	if b == nil {
		return ""
	}
	return b.roomnum
}

func (b *Battle) Connect(msg map[string]interface{}, conn *ws.Conn) bool {
	// fmt.Println("-- Battle.Connect() --")
	uid := msg["uid"].(string)
	ok := b.room.Check(uid)
	if !ok {
		return true
	}
	ok, val := xx.Getnumber(msg, "round")
	if !ok {
		return true
	}
	round := int(val)
	if round == 0 { // back to room
		b.Sendconfig(conn, "back")
		return true
	}
	// in room n relink
	ok, val = xx.Getnumber(msg, "msgptr")
	if !ok {
		return true
	}
	msgptr := int(val)
	if round != b.roundcnt {
		msgptr = 0
	}
	return b.room.Relink(uid, conn, msgptr)
}

func (b *Battle) Connected() bool {
	return b.room.Connected()
}

func (b *Battle) Receive(msg map[string]interface{}, conn *ws.Conn) bool {
	if b.Ended() {
		return false
	}
	if !b.check(msg) {
		return true
	}
	opt := msg["opt"].(string)
	uid := msg["uid"].(string)
	switch opt {
	case "enter":
		b.room.Enter(uid, conn)
	case "leave":
		nobody := b.room.Leave(uid)
		if nobody {
			b.End()
		}
	case "back":
		b.room.Back(uid, conn)
	case "ready":
		b.readying <- msg
	case "bet":
		b.betting <- msg
	}
	return true
}

func (b *Battle) Enterable() bool {
	if b.began {
		return false
	}
	return b.room.Enterable()
}

func (b *Battle) Getconfig() map[string]interface{} {
	b.config["roundcnt"] = b.roundcnt
	return b.config
}

func (b *Battle) Setconfig(data map[string]interface{}, host string) bool {
	if !checkcfg(data) {
		return false
	}
	b.config = data
	b.config["host"] = host
	b.config["kind"] = b.kind
	b.config["roomnum"] = b.roomnum
	b.mode = data["mode"].(string)
	//b.turning = data["turning"].(string)
	b.roundnum = data["roundnum"].(int)
	return true
}

func (b *Battle) Sendconfig(conn *ws.Conn, act string) {
	msg := map[string]interface{}{
		"opt": "room", "status": "ok",
	}
	switch act {
	case "enter", "back":
		msg["act"] = act
		msg["data"] = b.Getconfig()
	default:
		msg["status"] = act
	}
	conn.Send(msg)
}

func (b *Battle) Msg(currents []*room.Player) map[string]interface{} {
	data := map[string]interface{}{}
	for _, p := range currents {
		j := strconv.Itoa(p.Idx)
		data[j] = p.Cardmsg()
	}
	return data
}

func (b *Battle) Run() {
	defer b.End()
	var tm []string
	running := false
	for b.roundcnt = 1; b.roundcnt <= b.roundnum; b.roundcnt++ {
		ok := b.ready()
		if !ok {
			break
		}
		if !running {
			running = true
			// if !b.pay() {
			// 	return
			// }
			tm = b.start()
		}
		ok = b.round()
		if !ok {
			break
		}
	}
	b.over(tm)
}

func (b *Battle) End() {
	b.ended = true
	b.match.Remove(b.roomnum)
	b.timer.Stop()
}

func (b *Battle) Ended() bool {
	if b == nil {
		return false
	}
	return b.ended
}

// implementation
func (b *Battle) ready() bool {
	wait := 17000
	d := time.Duration(wait)
	timer := time.NewTimer(d * time.Millisecond)
	b.when = "ready"
	defer func() {
		b.when = ""
		timer.Stop()
	}()
	for {
		select {
		case msg := <-b.readying:
			uid := msg["uid"].(string)
			if ok := b.room.Ready(uid); ok {
				b.began = true
				return true
			}
		case <-timer.C:
			if b.began {
				ok := b.room.Autoready()
				return ok
			}
		}
	}
	return false
}

// func (b *Battle) pay() bool {
// 	inf := "battle pay(): "
// 	ok, host := xx.Getstring(b.config, "host")
// 	if !ok {
// 		fmt.Println(inf + "invalid host!!")
// 		return false
// 	}
// 	var num float64
// 	switch b.roundnum {
// 	case 2, 5, 10:
// 		num = -1
// 	case 20:
// 		num = -2
// 	default:
// 		fmt.Println(inf + "invalid roundnum!!")
// 		return false
// 	}
// 	ok = user.Addroomcard(host, num)
// 	if !ok {
// 		fmt.Println(inf + "addroomcard failed!!")
// 		return false
// 	}
// 	return true
// }

func (b *Battle) round() bool {
	b.cards = card.Init()
	b.currents = b.room.Newround()
	d := b.currents[0]
	b.pot = 0
	b.board = []*card.Card{} //clear the board

	//blinds
	sb := b.currents[1]
	bb := b.currents[2]
	sb.Addstakes(LIMIT / 2)
	bb.Addstakes(LIMIT)

	// start message
	data := map[string]interface{}{
		"button": d.Idx, "round": b.roundcnt,
		"blind": LIMIT,
	}
	seats := map[string]interface{}{}
	for _, p := range b.currents {
		j := strconv.Itoa(p.Idx)
		seats[j] = p.Seatmsg()
	}
	data["seats"] = seats
	b.room.Sendactive("start", data)

	//deal cards
	for _, p := range b.currents {
		p.Init()
		b.deal(p, 2)
		p.Sendactive("handcards", p.Cardmsg())
	}

	//pre-flop
	ended := b.betround("preflop")
	if ended == 0 {
		//flop round
		for i := 0; i < 5; i++ {
			b.cards, b.board[i] = card.Pop(b.cards)
		}
		b.sendboard(3)
		ended = b.betround("flop")
		if ended == 0 {
			//turn
			b.sendboard(4)
			ended = b.betround("turn")
			if ended == 0 {
				//river
				b.sendboard(5)
				ended = b.betround("river")
			}
		}
	}

	switch ended {
	case 1:
		//翻出所有公共牌并比大小 ...
		b.sendboard(5)
		winners := winner(b.currents)
		winnings := b.pot / len(winners)
		data := map[string]interface{}{}
		data["winners"] = map[string]interface{}{}
		for _, p := range winners {
			p.Addchips(winnings)
			j := strconv.Itoa(p.Idx)
			data["winners"].(map[string]interface{})[j] = winnings
		}
		data["board"] = b.brdmsg()
		hc := map[string]interface{}{}
		for _, p := range b.currents {
			if !p.Folded() {
				hc[strconv.Itoa(p.Idx)] = p.Cards().Msg()
			}
		}
		data["handcards"] = hc
		b.room.Sendactive("final", data)
	case 2:
		//奖池中的全部筹码归唯一没有弃牌的玩家 ...
		var winner *room.Player
		for _, p := range b.currents {
			if p.Active() && !p.Folded() {
				winner = p
				break
			}
		}
		data := map[string]interface{}{"1": winner.Idx}
		b.room.Sendactive("final", map[string]interface{}{"winner": data, "winnings": b.pot})
	}
	return true
}

func (b *Battle) deal(player *room.Player, num int) *card.Card {
	var c *card.Card
	for i := 0; i < num; i++ {
		b.cards, c = card.Pop(b.cards)
		player.Obtain(c)
	}
	return c
}

func (b *Battle) betround(round string) int {
	//用返回值表示下注回合是否结束
	//head参数并不是真正的最先行动的人，要首先确定其没有弃牌

	//首先找到真正的head
	b.when = "bet"
	var head, crt, min int
	if round == "preflop" {
		head = 2
		crt = next(head, b.currents)
		min = LIMIT
	} else {
		crt = next(0, b.currents)
		head = -1
		min = 0
	}
	for ; ; crt = next(crt, b.currents) {
		//首先确定当前玩家能够进行的动作
		s := b.currents[crt].Stakes
		c := b.currents[crt].Chips
		var checkable bool
		var per []string
		var data map[string]interface{}
		if crt == head {
			//一圈下注已经结束，回到第一个下(加)注的人，那么该轮下注结束
			break
		} else if s >= min {
			//当前玩家能够选择看牌、加注、弃牌、全下
			per = []string{"allin", "check", "fold", "raise"}
			checkable = true
		} else if min-s >= c {
			//可以弃牌、全下
			per = []string{"allin", "fold"}
			checkable = false
		} else {
			//可以跟注、加注、弃牌、全下
			per = []string{"allin", "call", "fold", "raise"}
			checkable = false
		}
		data = permsg(per)
		b.currents[crt].Sendactive("permissions", data)
		//接收玩家的身份信息uid操作信息act, 阻塞并等待
		var act, uid string
		var num int
		func() {
			//接收下注消息的时候要考虑一件事：如果收到该玩家上一轮的下注消息怎么办？...
			//考虑是不是要在消息中添加一个字段用来标记是哪一轮下注
			wait := 14000
			d := time.Duration(wait)
			timer := time.NewTimer(d * time.Millisecond)
			defer timer.Stop()
			for {
				select {
				case msg := <-b.betting:
					uid = msg["uid"].(string)
					data := msg["data"].(map[string]interface{})
					act = data["act"].(string)
					num = int(msg["num"].(float64))
					legal := false
					for _, s := range per {
						if s == act {
							legal = true
							break
						}
					}
					if uid == b.currents[crt].Uid && legal {
						// 还没考虑num是否合法 ...
						return
					}
				case <-timer.C:
					uid = b.currents[crt].Uid
					if checkable {
						act = "check"
					} else {
						act = "fold"
					}
					return
				}
			}
		}()
		switch act {
		case "call":
			//跟注
			b.currents[crt].Call(min)
			num = min
		case "raise":
			//加注，更新head
			b.currents[crt].Raise(num, min)
			min += num
			head = crt
		case "fold":
			//弃牌
			b.currents[crt].Fold()
			b.Addpot(s)
		case "allin":
			//全下
			b.currents[crt].All_in()
			num = c
		case "check":
			//看牌,好像没什么做的？
		}
		//将当前玩家的动作广播给所有玩家
		b.sendactions(b.currents[crt].Idx, act, num)
	}
	//将所有人的下注加到奖池中，并广播该轮下注结束后奖池数量以及每个玩家剩下的筹码
	chips := map[string]interface{}{}
	for _, p := range b.currents {
		if p.Active() && !p.Folded() {
			b.Addpot(p.Stakes)
			p.Stakes = 0
			chips[strconv.Itoa(p.Idx)] = p.Chips
		}
	}

	//判断之后是否还需要继续下一轮下注
	//结束的情况有三种：所有人都all in、只剩下一个人没有弃牌、河牌圈已经结束
	all_allin := true
	not_fold := 0
	for _, p := range b.currents {
		if p.Active() {
			if all_allin && !p.Allin() {
				all_allin = false
			}
			if !p.Folded() {
				not_fold++
			}
		}
	}

	b.when = ""
	if all_allin || round == "river" {
		return 1
	}
	if not_fold == 1 {
		return 2
	}
	data := map[string]interface{}{"pot": b.pot, "chips": chips}
	b.room.Sendactive("over", data)
	return 0
}

func (b *Battle) Addpot(n int) {
	b.pot += n
}

func (b *Battle) sendactions(idx int, act string, num int) {
	data := map[string]interface{}{}
	data["idx"] = idx
	data["action"] = act
	data["num"] = num
	b.room.Sendactive("bet", data)
}

func (b *Battle) sendboard(n int) {
	data := map[string]interface{}{}
	tmp := b.board[:n]
	if len(tmp) != 0 {
		for i, e := range tmp {
			s := strconv.Itoa(i + 1)
			data[s] = e.Msg()
		}
	}
	for _, p := range b.currents {
		c := b.board[:n]
		c = append(c, p.Cards().Card(0))
		c = append(c, p.Cards().Card(1))
		max := card.CombinationTraversal(c)
		data["rank"] = max.Rank()
		p.Sendactive("board", data)
	}
}

func next(c int, currents []*room.Player) int {
	num := len(currents)
	next := -1
	for i := c + 1; next < 0; i++ {
		i = i % num
		if currents[i].Status == "active" && !currents[i].Folded() {
			next = i
		}
	}
	return next
}

func permsg(p []string) map[string]interface{} {
	data := map[string]interface{}{}
	for i, s := range p {
		j := strconv.Itoa(i + 1)
		data[j] = s
	}
	return data
}

func winner(currents []*room.Player) []*room.Player {
	//winner := currents[0]
	winner := []*room.Player{currents[0]}
	for _, p := range currents[1:] {
		rst := p.Compare(winner[0])
		if rst == 0 {
			winner = append(winner, p)
		} else if rst > 0 {
			winner = []*room.Player{p}
		}
	}
	return winner
}

// record of batlle
func (b *Battle) start() []string {
	b.room.Save("data.roomnum", b.roomnum)
	t := []string{}
	now := time.Now().Unix()
	local := time.Unix(now, 0)
	t = append(t, local.Format("2006-01-02"))
	t = append(t, local.Format("15:04"))
	return t
}

func (b *Battle) time() []string {
	arr := []string{}
	now := time.Now().Unix()
	t := time.Unix(now, 0)
	arr = append(arr, t.Format("2006-01-02"))
	arr = append(arr, t.Format("15:04"))
	return arr
}

func (b *Battle) over(tm []string) {
	tbl := map[string]interface{}{"time": tm}
	tbl["roomnum"] = b.roomnum
	tbl["content"] = b.room.Getrecord()
	b.room.Record(b.kind, tbl)
	b.room.Save("data.roomnum", "")
	b.room.Save("data.roomnum", "")
	//b.room.Send("over", b.statistic())
	b.room.Init()
}

// check messages from clients
func (b *Battle) check(msg map[string]interface{}) bool {
	inf := "battle.check(): invalid msg %v: %v!!\n"
	opt := msg["opt"].(string)
	switch opt {
	case "enter", "leave":
		if b.began {
			fmt.Printf(inf, "opt", opt)
			fmt.Println("battle has began...")
			return false
		}
	case "back":
		if !b.began {
			fmt.Printf(inf, "opt", opt)
			fmt.Println("battle has not began...")
			return false
		}
	//这里要改 ...
	case "bet":
		if b.when != "bet" {
			return false
		}
	case "ready":
		if b.when != "ready" {
			return false
		}
	default:
		fmt.Println(inf, "opt", opt)
		return false
	}
	return true
}

func (b *Battle) brdmsg() map[string]interface{} {
	m := map[string]interface{}{}
	for i, c := range b.board {
		j := strconv.Itoa(i + 1)
		m[j] = c.Msg()
	}
	return m
}
func checkcfg(data map[string]interface{}) bool {
	cfg, err := xxio.Read("niuniu")
	if err != nil {
		return false
	}
	for k, v := range data {
		ok, tbl := xx.Getmap(cfg, k)
		if !ok {
			return false
		}
		key, ok := v.(string)
		if !ok {
			return false
		}
		ok, num := xx.Getnumber(tbl, key)
		if !ok {
			return false
		}
		if num != 0 {
			data[k] = int(num)
		}
	}
	return true
}
