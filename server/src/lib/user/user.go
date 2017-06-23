package user

import (
	"../mongo"
	"../xx"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"gopkg.in/mgo.v2/bson"
	"log"
	"math"
	"strconv"
	"strings"
	"time"
)

const DB = "YBH"

// create a new user
func New(token map[string]interface{}) map[string]interface{} {
	data, err := xx.Read("user")
	if err != nil {
		log.Println("New() ", err)
		return nil
	}
	ok, _ := xx.Getnumber(data, "roomcard")
	if !ok {
		log.Println("New() ", "user.json error!")
		return nil
	}
	ok, openid := xx.Getstring(token, "openid")
	if !ok {
		log.Println("New() ", "token error!")
		return nil
	}
	key, ok := generateKey(openid)
	if !ok {
		log.Println("New() generate key failed!")
		return nil
	}
	data["key"] = key
	uid, ok := generateUid()
	if !ok {
		log.Println("New() generate uid failed!")
		return nil
	}
	data["uid"] = uid
	doc := map[string]interface{}{"uid": uid}
	doc["data"] = data
	now := time.Now().Unix()
	tm := time.Unix(now, 0)
	doc["time"] = tm.Format("2006-01-02 15:04:05")
	err = mongo.Insert(DB, "user", doc)
	if err != nil {
		log.Println("New() ", err)
		return nil
	}
	return doc
}

func ResetKey(uid string) bool {
	val, err := Search(uid, "token.openid")
	if err != nil {
		log.Println("ResetKey() ", err)
		return false
	}
	openid, ok := val.(string)
	if !ok {
		log.Println("ResetKey() openid type error!")
		return false
	}
	key, ok := generateKey(openid)
	if !ok {
		log.Println("ResetKey() generate key failed!")
		return false
	}
	err = Upsert(uid, "data.key", key)
	if err != nil {
		log.Println("ResetKey() ", err)
		return false
	}
	return true
}

// roundnum: "one", "two"
func Checkroomcard(uid, roundnum string) bool {
	var cost float64
	switch roundnum {
	case "one":
		cost = 1
	case "two":
		cost = 2
	default:
		log.Println("Checkroomcard() invalid roundnum!")
		return false
	}
	val, err := Search(uid, "data.roomcard")
	if err != nil {
		log.Println("Checkroomcard() ", err)
		return false
	}
	num, ok := val.(float64)
	if !ok {
		log.Println("Checkroomcard() roomcard type error!")
		return false
	}
	return num >= cost
}

func Addroomcard(uid string, n float64) bool {
	if n > 10000 {
		log.Println("Addroomcard() range error!")
		return false
	}
	n = math.Floor(n)
	val, err := Search(uid, "data.roomcard")
	if err != nil {
		log.Println("Addroomcard() ", err)
		return false
	}
	num, ok := val.(float64)
	if !ok {
		log.Println("Addroomcard() roomcard type error!")
		return false
	}
	if num+n < 0 {
		log.Println("Addroomcard() roomcard can not be negative!")
		return false
	}
	err = Inc(uid, "data.roomcard", float64(n))
	if err != nil {
		log.Println("Addroomcard() ", err)
		return false
	}
	return true
}

// find doc in user's data
func Find(uid string) (map[string]interface{}, error) {
	sel := bson.M{}
	doc, err := mongo.Find(DB, "user", bson.M{"uid": uid}, sel)
	if err != nil {
		log.Println("Find() ", err)
		return nil, err
	}
	return doc, nil
}

func FindBy(pth, val string) (map[string]interface{}, error) {
	sel := bson.M{}
	doc, err := mongo.Find(DB, "user", bson.M{pth: val}, sel)
	if err != nil {
		log.Println("FindBy() ", err)
		return nil, err
	}
	return doc, nil
}

// search user's infomation item
func Search(uid, pth string) (interface{}, error) {
	doc, err := Find(uid)
	if err != nil {
		log.Println("Search() ", err)
		return nil, err
	}
	if doc == nil {
		return nil, nil
	}
	if pth == "" {
		return doc, nil
	}
	keys := strings.Split(pth, ".")
	k := keys[0]
	for i := 0; i < len(keys)-1; i++ {
		ok, val := xx.Getmap(doc, k)
		if !ok {
			return nil, nil
		}
		k = keys[i+1]
		doc = val
	}
	return doc[k], nil
}

// update user's infomation item
func Upsert(uid, pth string, val interface{}) error {
	err := mongo.Upsert(DB, "user", bson.M{"uid": uid}, bson.M{pth: val})
	if err != nil {
		log.Println("Upsert() ", err)
	}
	return err
}

// increase user's infomation item
func Inc(uid, pth string, num float64) error {
	err := mongo.Inc(DB, "user", bson.M{"uid": uid}, bson.M{pth: num})
	if err != nil {
		log.Println("Inc() ", err)
	}
	return err
}

// unset key in user data
func Unset(uid, pth string) error {
	err := mongo.Unset(DB, "user", bson.M{"uid": uid}, bson.M{pth: 1})
	if err != nil {
		log.Println("Unset() ", err)
	}
	return err
}

// remove user's full data
func Remove(uid string) error {
	err := mongo.Remove(DB, "user", bson.M{"uid": uid})
	if err != nil {
		log.Println("Remove() ", err)
	}
	return err
}

// insert element to array
func ArrayInsert(uid, name, pth string, content interface{}) error {
	var condition bson.M
	if uid != "" {
		condition = bson.M{"uid": uid}
	} else {
		condition = bson.M{"name": name}
	}
	err := mongo.AddToSet(DB, "user", condition, bson.M{pth: content})
	if err != nil {
		log.Println("ArrayInsert() ", err)
	}
	return err
}

// remove element from array
func ArrayRemove(uid, name, pth string, content interface{}) error {
	var condition bson.M
	if uid != "" {
		condition = bson.M{"uid": uid}
	} else {
		condition = bson.M{"name": name}
	}
	err := mongo.Pull(DB, "user", condition, bson.M{pth: content})
	if err != nil {
		log.Println("ArrayRemove() ", err)
	}
	return err
}

// implementation
// generate user's ID
func generateUid() (string, bool) {
	num, ok := getcount()
	if !ok {
		log.Println("generateUid() getcount failed!")
		return "", false
	}
	num++
	ok = setcount(num)
	if !ok {
		log.Println("generateUid() setcount failed!")
		return "", false
	}
	num += 100000
	uid := strconv.Itoa(int(num))
	return uid, true
}

func getcount() (float64, bool) {
	val, err := Search("xx", "count")
	if err != nil {
		log.Println("getcount() ", err)
		return 0, false
	}
	if val == nil {
		ok := setcount(0)
		return 0, ok
	}
	count, ok := val.(float64)
	if !ok {
		log.Println("getcount() count type error!")
		return 0, false
	}
	return count, true
}

func setcount(num float64) bool {
	err := Upsert("xx", "count", num)
	if err != nil {
		log.Println("setcount() ", "", err)
		return false
	}
	return true
}

func generateKey(openid string) (string, bool) {
	data := md5.Sum([]byte(openid))
	base := data[:]
	salt := make([]byte, 12)
	_, err := rand.Read(salt)
	if err != nil {
		log.Println("generateKey() ", err)
		return "", false
	}
	s1 := hex.EncodeToString(base)
	s2 := hex.EncodeToString(salt)
	key := s1 + s2
	return key, true
}
