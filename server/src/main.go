package main

import (
	"./battle"
	"./lib/user"
	"./lib/ws"
	"./lib/xx"
	// "fmt"
	"net/websocket"
)

var match *battle.Match
var players map[string]*ws.Conn

func main() {
	match = battle.NewMatch(2000)
	players = map[string]*ws.Conn{}
	// battle.Test()
	ws.Listen("8000", "/receive", receive)
}

func receive(c *websocket.Conn) {
	conn := ws.New(c)
	defer conn.Close()
	ok := true
	var b *battle.Battle
	for ok {
		msg, opt := conn.Receive()
		if msg == nil {
			break
		}
		switch opt {
		case "heartbeat":
			break
		case "connect":
			ok, b = connect(msg, conn)
		case "create":
			ok, b = create(msg, conn)
		case "access":
			ok, b = access(msg, conn)
		case "invite":
			ok = invite(msg, b.Roomnum())
		default:
			ok = b.Received(msg, conn)
		}
	}
}

func connect(msg map[string]interface{}, conn *ws.Conn) (bool, *battle.Battle) {
	ok, cnn := xx.Getstring(msg, "connect")
	if !ok {
		return false, nil
	}
	uid := msg["uid"].(string)
	if cnn == "off" {
		delete(players, uid)
		return false, nil
	}
	ok, b := checkroom(msg, conn)
	if !ok {
		return false, nil
	}
	players[uid] = conn
	return true, b
}

func checkroom(msg map[string]interface{}, conn *ws.Conn) (bool, *battle.Battle) {
	inf := "main checkroom(): "
	uid := msg["uid"].(string)
	val, err := user.Search(uid, "data.roomnum")
	if err != nil {
		println(inf + "user.Search failed!")
		return false, nil
	}
	roomnum, ok := val.(string)
	if !ok {
		println(inf + "roomnum type err in db!")
		return false, nil
	}
	b := match.Get(roomnum)
	if b == nil {
		res := map[string]interface{}{
			"opt": "connect", "status": "ok",
		}
		conn.Send(res)
		return true, nil
	}
	ok = b.Connect(msg, conn)
	if !ok {
		println(inf + "battle.Connect failed!")
		return false, nil
	}
	return true, b
}

func create(msg map[string]interface{}, conn *ws.Conn) (bool, *battle.Battle) {
	var b *battle.Battle
	uid := msg["uid"].(string)
	ok := user.Checkroomcard(uid)
	if !ok {
		b.Sendconfig(conn, "lack")
		return false, nil
	}
	b = match.Create()
	if b == nil {
		b.Sendconfig(conn, "fail")
		return true, nil
	}
	ok, data := xx.Getmap(msg, "config")
	if !ok {
		return false, nil
	}
	ok = b.Setconfig(data, uid)
	if !ok {
		b.Sendconfig(conn, "fail")
		return false, nil
	}
	b.Sendconfig(conn, "enter")
	go b.Run()
	return true, b
}

func access(msg map[string]interface{}, conn *ws.Conn) (bool, *battle.Battle) {
	var b *battle.Battle
	ok, num := xx.Getstring(msg, "roomnum")
	if !ok {
		b.Sendconfig(conn, "unexist")
		return false, nil
	}
	b = match.Get(num)
	if b == nil {
		b.Sendconfig(conn, "unexist")
		return true, nil
	}
	if !b.Enterable() {
		b.Sendconfig(conn, "full")
		return true, nil
	}
	b.Sendconfig(conn, "enter")
	return true, b
}

func invite(msg map[string]interface{}, roomnum string) bool {
	if roomnum == "" {
		return true
	}
	ok, who := xx.Getstring(msg, "who")
	if !ok {
		return false
	}
	uid := msg["uid"].(string)
	delete(msg, "uid")
	msg["who"] = uid
	msg["roomnum"] = roomnum
	conn := players[who]
	conn.Send(msg)
	return true
}
