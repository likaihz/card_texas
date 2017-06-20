// cli.go
package main

import (
	"./lib/xx"
	// "bytes"
	"encoding/json"
	"flag"
	"fmt"
	// "io"
	"io/ioutil"
	"log"
	"net/http"
	"net/websocket"
	"os"
	"strings"
)

var origin = "http://127.0.0.1:8000/"
var url = "ws://127.0.0.1:8000/receive"

func main() {
	flag.Parse()
	if flag.NArg() > 0 {
		if flag.Arg(0) == "login" {
			buf := "\"test\": \"\"}"
			//jsonstr := []byte(buf)
			client := &http.Client{}
			reqest, err := http.NewRequest("POST", "http://127.0.0.1:8001/test", strings.NewReader(buf))
			_, err = client.Do(reqest)
			if err != nil {
				fmt.Println("Fatal error ", err.Error())
				os.Exit(0)
			}
		}
	} else {
		ws, err := websocket.Dial(url, "", origin)
		if err != nil {
			log.Fatal(err)
		}
		var mapmsg map[string]interface{}
		mapmsg, err = Read()
		// message := []byte("hello, world!你好")
		msg2 := xx.Map2str(mapmsg)
		// _, err = ws.Write(message)
		err = websocket.Message.Send(ws, msg2)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Send: %s\n", msg2)

		var msg = make([]byte, 512)
		m, err := ws.Read(msg)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Receive: %s\n", msg[:m])
		ws.Close() //关闭连接
	}
}

func Read() (map[string]interface{}, error) {
	file, err := ioutil.ReadFile("./data/inst.json")
	if err != nil {
		return nil, err
	}
	var data map[string]interface{}
	err = json.Unmarshal(file, &data)
	fmt.Println(data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
