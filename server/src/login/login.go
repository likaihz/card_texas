package main

import (
	"../lib/xxio"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/test", test)
	err := http.ListenAndServe(":8001", nil)
	if err != nil {
		fmt.Println("login server start failed!!")
	}
}

const (
	GNUM = 2
)

var gcount = 0

// fade client user login
func test(w http.ResponseWriter, r *http.Request) {
	fmt.Println("server: login...")
	res := map[string]interface{}{
		"status": "fail",
	}
	defer xxio.Response(w, &res)

	req, err := xxio.Request(r)
	if err != nil {
		return
	}
	cfg, err := getconfig()
	if err != nil {
		return
	}
	res["cfg"] = cfg
	data, err := popuser()
	if err != nil {
		return
	}
	if data == nil {
		// all users are online!
		res["status"] = "ok"
		return
	}
	res["user"] = data
	// room control
	mode := req["test"].(string)
	gcount++
	generate(res, mode, gcount)
	if gcount == GNUM {
		gcount = 0
	}
	res["status"] = "ok"
}

func getconfig() (map[string]interface{}, error) {
	configs := map[string]interface{}{
		"tank": 0,
	}
	data := map[string]interface{}{}
	for name := range configs {
		file, err := xxio.Read(name)
		if err != nil {
			return nil, err
		}
		data[name] = file
	}
	return data, nil
}

func generate(res map[string]interface{}, mode string, gcount int) {
	switch mode {
	case "room", "team":
		if gcount == 1 {
			res["opt"] = "create"
		} else {
			res["opt"] = "join"
		}
		res["cls"] = mode
	case "match":
		if gcount == 1 {
			res["cls"] = "room"
		} else {
			res["cls"] = "team"
		}
		res["opt"] = "create"
	}
}

// 跳过数据库操作
var USER map[string]interface{}
var READ bool

func popuser() (map[string]interface{}, error) {
	if !READ {
		var err error
		USER, err = xxio.Read("user")
		if err != nil {
			return nil, err
		}
		READ = true
	}
	if len(USER) <= 0 {
		return nil, nil
	}
	var k string
	for k = range USER {
		break
	}
	data := USER[k].(map[string]interface{})
	delete(USER, k)
	return data, nil
}
