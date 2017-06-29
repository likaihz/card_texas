package main

import (
	"./battle"
	// "./lib/user"
	"./lib/ws"
	"./lib/xx"
	"log"
	"net/websocket"
)

var match *battle.Match

func main() {
	xx.Openlog("main")
	defer xx.Closelog()
	match = battle.NewMatch(2000)
	ws.Listen("8000", "/receive", receive)
}

func receive(c *websocket.Conn) {

	conn := ws.New(c)
	defer conn.Close()
	log.Println("in receive")
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
			log.Println(msg)
			ok, b = connect(msg, conn)
		case "create":
			log.Println(msg)
			ok, b = create(msg, conn)
		case "access":
			log.Println(msg)
			ok, b = access(msg, conn)
		default:
			log.Println("receive default", msg)
			ok = b.Receive(msg, conn)
		}
	}
}

func connect(msg map[string]interface{}, conn *ws.Conn) (bool, *battle.Battle) {
	var b *battle.Battle
	ok, cnn := xx.Getstring(msg, "connect")
	if !ok {
		log.Println("connect() invalid msg!")
		return false, nil
	}
	if cnn == "off" {
		return false, nil
	}
	// uid := msg["uid"].(string)
	// val, err := user.Search(uid, "data.roomnum")
	// if err != nil {
	// 	log.Println("connect() ", err)
	// 	return false, nil
	// }
	// roomnum, ok := val.(string)
	// if !ok {
	// 	log.Println("connect() roomnum type error in db!")
	// 	return false, nil
	// }
	res := map[string]interface{}{
		"opt": "connect", "status": "ok",
	}
	// b := match.Get(roomnum)
	if b == nil {
		conn.Send(res)
		return true, nil
	}
	ok = b.Connect(msg, conn)
	if !ok {
		conn.Send(res)
		return true, nil
	}
	return true, b
}

func create(msg map[string]interface{}, conn *ws.Conn) (bool, *battle.Battle) {
	ok, cfg := xx.Getmap(msg, "config")
	if !ok {
		log.Println("create() invalid key: config!")
		return false, nil
	}
	// fmt.Println(cfg)
	// ok, roundnum := xx.Getstring(cfg, "roundnum")
	ok, _ = xx.Getstring(cfg, "roundnum")
	if !ok {
		log.Println("create() invalid key: roundnum!")
		return false, nil
	}
	uid := msg["uid"].(string)
	// ok = user.Checkroomcard(uid, roundnum)
	// if !ok {
	// 	b.Sendconfig(conn, "lack")
	// 	return true, nil
	// }
	b := match.Create()
	if b == nil {
		b.Sendconfig(conn, "fail")
		return true, nil
	}
	ok = b.Setconfig(cfg, uid)
	if !ok {
		b.Sendconfig(conn, "fail")
		return true, nil
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
		// senderror(conn, "unexist")
		return true, nil
	}
	if !b.Enterable() {
		b.Sendconfig(conn, "full")
		return true, nil
	}
	b.Sendconfig(conn, "enter")
	return true, b
}
