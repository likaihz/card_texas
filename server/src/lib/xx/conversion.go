package xx

import (
	"encoding/json"
	"fmt"
)

func Str2map(s string) map[string]interface{} {
	m := map[string]interface{}{}
	err := json.Unmarshal([]byte(s), &m)
	if err != nil {
		fmt.Println("Str2map() err: ", err)
		return nil
	}
	return m
}

func Map2str(m map[string]interface{}) string {
	buf, err := json.Marshal(m)
	if err != nil {
		fmt.Println("Map2str() err: ", err)
		return ""
	}
	return string(buf)
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
