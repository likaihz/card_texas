package main

import (
	"./battle"
	"./lib/user"
	"./lib/ws"
	"./lib/xx"
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
			ok = invite(msg, b.Number())
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
	ok, inroom := xx.Getstring(msg, "inroom")
	if !ok {
		return false, nil
	}
	if inroom == "no" {
		return true, nil
	}
	ok, val := xx.Getnumber(msg, "roomnum")
	if !ok {
		return false, nil
	}
	num := int(val)
	b := match.Get(num)
	if b == nil {
		conn.Send(map[string]interface{}{
			"opt": "out",
		})
	} else {
		ok := b.Connect(msg, conn)
		if !ok {
			return false, nil
		}
	}
	return true, b
}

func create(msg map[string]interface{}, conn *ws.Conn) (bool, *battle.Battle) {
	uid := msg["uid"].(string)
	ok := user.Checkroomcard(uid)
	if !ok {
		senderror(conn, "lack")
		return false, nil
	}
	b := match.Create()
	if b == nil {
		senderror(conn, "fail")
		return true, nil
	}
	ok, data := xx.Getmap(msg, "config")
	if !ok {
		return false, nil
	}
	ok = b.Setconfig(data, uid)
	if !ok {
		senderror(conn, "fail")
		return false, nil
	}
	data = b.Getconfig()
	sendconfig(conn, data)
	go b.Run()
	return true, b
}

func access(msg map[string]interface{}, conn *ws.Conn) (bool, *battle.Battle) {
	ok, val := xx.Getnumber(msg, "roomnum")
	if !ok {
		senderror(conn, "unexist")
		return false, nil
	}
	num := int(val)
	b := match.Get(num)
	if b == nil {
		senderror(conn, "unexist")
		return true, nil
	}
	if !b.Enterable() {
		senderror(conn, "full")
		return true, nil
	}
	data := b.Getconfig()
	sendconfig(conn, data)
	return true, b
}

func invite(msg map[string]interface{}, roomnum int) bool {
	if roomnum < 0 {
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

// implementation
func senderror(conn *ws.Conn, err string) {
	msg := map[string]interface{}{
		"opt": "front", "status": err,
	}
	conn.Send(msg)
}

func sendconfig(conn *ws.Conn, data map[string]interface{}) {
	msg := map[string]interface{}{
		"opt": "front", "status": "ok",
	}
	msg["data"] = data
	conn.Send(msg)
}
