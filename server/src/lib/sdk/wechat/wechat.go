package wechat

import (
	"../../xx"
	"log"
	"net/http"
)

const (
	DOMAIN = "https://api.weixin.qq.com/sns/"
	APPID  = "appid=" + "wx343538f7838ec29a"
	SECRET = "secret=" + "715bbff77206aab95eb101fd9d288603"
)

func Access(val interface{}) map[string]interface{} {
	code, ok := val.(string)
	if !ok {
		log.Println("Access() invalid code!")
		return nil
	}
	url := geturl("code", code, "")
	return fetch(url)
}

func Refresh(token map[string]interface{}) map[string]interface{} {
	ok, refresh := xx.Getstring(token, "refresh_token")
	if !ok {
		log.Println("Refresh() invalid refresh_token")
		return nil
	}
	url := geturl("refresh", refresh, "")
	return fetch(url)
}

func Getuser(token map[string]interface{}) map[string]interface{} {
	ok, access := xx.Getstring(token, "access_token")
	if !ok {
		log.Println("Refresh() invalid access_token")
		return nil
	}
	ok, openid := xx.Getstring(token, "openid")
	if !ok {
		log.Println("Refresh() invalid openid")
		return nil
	}
	url := geturl("access", access, openid)
	data := fetch(url)
	if data == nil {
		log.Println("Refresh() fetch failed!")
		return nil
	}
	delete(data, "openid")
	delete(data, "unionid")
	delete(data, "privilege")
	return data
}

// implementation
// key: "code", "access", "refresh"
func geturl(key, v1, v2 string) string {
	url := DOMAIN
	switch key {
	case "code":
		url += "oauth2/access_token?"
		url += APPID + "&" + SECRET
		url += "&code=" + v1
		url += "&grant_type=authorization_code"
	case "access":
		url += "userinfo?"
		url += "access_token=" + v1
		url += "&openid=" + v2
	case "refresh":
		url += "oauth2/refresh_token?" + APPID
		url += "&grant_type=refresh_token"
		url += "&refresh_token=" + v1
	}
	return url
}

func fetch(url string) map[string]interface{} {
	res, err := http.Get(url)
	if err != nil {
		log.Println("fetch() ", err)
		return nil
	}
	data, err := xx.Decode(res.Body)
	res.Body.Close()
	if err != nil {
		log.Println("fetch() ", err)
		return nil
	}
	msg, ok := data["errmsg"]
	if ok {
		log.Println("fetch() return error: ", msg)
		return nil
	}
	return data
}
