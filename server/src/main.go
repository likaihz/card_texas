package main

import (
	//"./lib/xx"
	//"./match"
	"./test"
	"flag"
	//"fmt"
	//"log"
	//"net/http"
	//"net/websocket"
	//"os"
)

// var rooms *match.Groups
// var player map[string](chan string)

// var lg *log.Logger
// var file *os.File

func main() {
	flag.Parse()
	var args []string = flag.Args()
	if len(args) == 0 {
		// lg, file = xx.Log2file("main.log")
		// rooms = match.NewGroups(1000)
		// player = map[string](chan string){}
		// http.Handle("/receive", websocket.Handler(receive))
		// err := http.ListenAndServe(":8000", nil)
		// if err != nil {
		// 	panic("ListenAndServe: " + err.Error())
		return
		//}
	} else if args[0] == "cardtest" {
		test.CardTest()
	} else if args[0] == "comparetest" {
		test.CompareTest()
	}
}

// func receive(conn *websocket.Conn) {
// 	defer conn.Close()
// 	ch := make(chan string)
// 	go send(conn, ch)

// 	var e match.Receptor
// 	for {
// 		var req string
// 		err := websocket.Message.Receive(conn, &req)
// 		if err != nil {
// 			fmt.Println(err)
// 			break
// 		}
// 		msg := xx.Str2map(req)
// 		opt := msg["opt"].(string)
// 		switch opt {
// 		case "connect":
// 			if !connect(msg, ch) {
// 				return
// 			}
// 			break
// 		case "create":
// 			e = create(msg, ch)
// 			break
// 		case "leave":
// 			leave(msg, ch, e)
// 			break
// 		case "invite":
// 			invite(msg, e)
// 			break
// 		case "accept":
// 			e = accept(msg, ch)
// 			break
// 		case "access":
// 			e = access(msg, ch)
// 			break
// 		case "dismiss":
// 			// ...
// 			break
// 		default:
// 			e.Receive(msg)
// 		}
// 	}
// }

// func send(conn *websocket.Conn, ch <-chan string) {
// 	for msg := range ch {
// 		err := websocket.Message.Send(conn, msg)
// 		if err != nil {
// 			lg.Println(err)
// 			break
// 		}
// 	}
// }

// func connect(msg map[string]interface{}, ch chan string) bool {
// 	uid := msg["uid"].(string)
// 	open := msg["open"].(string)
// 	if open == "on" {
// 		player[uid] = ch
// 		return true
// 	}
// 	delete(player, uid)
// 	return false
// }

// func create(msg map[string]interface{}, ch chan string) match.Receptor {
// 	uid := msg["uid"].(string)
// 	e := rooms.Create(uid, ch)
// 	e.SetConfig(msg)
// 	return e
// }

// func leave(msg map[string]interface{}, ch chan string, e match.Receptor) {
// 	uid := msg["uid"].(string)
// 	e.Leave(uid, ch)
// 	if e.Occupancy() == 0 {
// 		rooms.Remove(e)
// 	}
// }

// func invite(msg map[string]interface{}, e match.Receptor) {
// 	who := msg["who"].(string)
// 	msg["idx"] = e.Index()
// 	msg["who"] = "xxx"
// 	friend := player[who]
// 	friend <- xx.Map2str(msg)
// }

// func accept(msg map[string]interface{}, ch chan string) match.Receptor {
// 	uid := msg["uid"].(string)
// 	idx := msg["idx"].(float64)
// 	e := rooms.Get(int(idx))
// 	if e != nil {
// 		e.Enter(uid, ch)
// 	}
// 	return e
// }

// func access(msg map[string]interface{}, ch chan string) match.Receptor {
// 	e := rooms.Pop()
// 	if e != nil {
// 		uid := msg["uid"].(string)
// 		e.Enter(uid, ch)
// 	} else {
// 		e = create(msg, ch)
// 	}
// 	return e
// }
