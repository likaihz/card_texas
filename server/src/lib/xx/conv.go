package xx

import (
	"encoding/json"
	"log"
)

func Str2map(s string) map[string]interface{} {
	m := map[string]interface{}{}
	err := json.Unmarshal([]byte(s), &m)
	if err != nil {
		log.Println("Str2map() ", err)
		return nil
	}
	return m
}

func Map2str(m map[string]interface{}) string {
	buf, err := json.Marshal(m)
	if err != nil {
		log.Println("Map2str() ", err)
		return ""
	}
	return string(buf)
}

func Getstring(msg map[string]interface{}, k string) (bool, string) {
	val, ok := msg[k]
	if !ok {
		log.Println("Getstring() unexist key: ", k)
		return false, ""
	}
	str, ok := val.(string)
	if !ok {
		log.Printf("Getstring() value of \"%s\" is not a string!\n", k)
		return false, ""
	}
	return true, str
}

func Getnumber(msg map[string]interface{}, k string) (bool, float64) {
	val, ok := msg[k]
	if !ok {
		log.Println("Getnumber() unexist key: ", k)
		return false, 0
	}
	num, ok := val.(float64)
	if !ok {
		log.Printf("Getnumber() value of \"%s\" is not a number!\n", k)
		return false, 0
	}
	return true, num
}

func Getmap(msg map[string]interface{}, k string) (bool, map[string]interface{}) {
	val, ok := msg[k]
	if !ok {
		log.Println("Getmap() unexist key: ", k)
		return false, nil
	}
	tbl, ok := val.(map[string]interface{})
	if !ok {
		log.Println("Getmap() value of \"%s\" is not a map!", k)
		return false, nil
	}
	return true, tbl
}

func Restrict(num, a, b float64) float64 {
	if num < a {
		return a
	}
	if num > b {
		return b
	}
	return num
}

func Sign(num float64) float64 {
	if num > 0 {
		return 1
	} else if num < 0 {
		return -1
	}
	return 0
}
