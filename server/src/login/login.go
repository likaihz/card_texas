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
	defer xxio.Response(w, &res, true)

	_, err := xxio.Request(r, true)
	if err != nil {
		return
	}
	data, err := popuser()
	fmt.Println(data)
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
	res["status"] = "ok"
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
