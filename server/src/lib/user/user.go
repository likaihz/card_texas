package user

import (
	"../xx"
	"../xxdb"
	"../xxio"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"gopkg.in/mgo.v2/bson"
	"math"
	"strconv"
	"strings"
	"time"
)

const (
	DB = "YBH"
)

// create a new user
func New(token map[string]interface{}) (map[string]interface{}, error) {
	data, err := xxio.Read("user")
	if err != nil {
		xxerr("New", "read", err)
		return nil, err
	}
	ok, _ := xx.Getnumber(data, "roomcard")
	if !ok {
		err = xxerr("New", "roomcard", nil)
		return nil, err
	}
	ok, openid := xx.Getstring(token, "openid")
	if !ok {
		err = xxerr("New", "openid", nil)
		return nil, err
	}
	key, err := generateKey(openid)
	if err != nil {
		xxerr("New", "generateKey", nil)
		return nil, err
	}
	data["key"] = key
	uid, err := generateUid()
	if err != nil {
		xxerr("New", "generateUid", nil)
		return nil, err
	}
	data["uid"] = uid
	doc := map[string]interface{}{"uid": uid}
	doc["data"] = data
	now := time.Now().Unix()
	tm := time.Unix(now, 0)
	doc["time"] = tm.Format("2006-01-02 15:04:05")
	err = xxdb.Insert(DB, "user", doc)
	if err != nil {
		xxerr("New", "insert", err)
		return nil, err
	}
	return doc, nil
}

func ResetKey(uid string) error {
	val, err := Search(uid, "token.openid")
	if err != nil {
		xxerr("ResetKey", "search", nil)
		return err
	}
	openid := val.(string)
	key, err := generateKey(openid)
	if err != nil {
		xxerr("ResetKey", "generateKey", nil)
		return err
	}
	err = Upsert(uid, "data.key", key)
	if err != nil {
		xxerr("ResetKey", "upsert", nil)
		return err
	}
	return nil
}

func Checkroomcard(uid string) bool {
	val, err := Search(uid, "data.roomcard")
	if err != nil {
		xxerr("Addroomcard", "search", err)
		return false
	}
	num, ok := val.(float64)
	if !ok || num < 1 {
		xxerr("Addroomcard", "type err", nil)
		return false
	}
	return true
}

func Addroomcard(uid string, n float64) bool {
	if n > 10000 {
		xxerr("Addroomcard", "range err", nil)
		return false
	}
	n = math.Floor(n)
	val, err := Search(uid, "data.roomcard")
	if err != nil {
		xxerr("Addroomcard", "search", err)
		return false
	}
	num, ok := val.(float64)
	if !ok || num+n < 0 {
		xxerr("Addroomcard", "type err", nil)
		return false
	}
	err = Inc(uid, "data.roomcard", float64(n))
	if err != nil {
		xxerr("Addroomcard", "inc", err)
		return false
	}
	return true
}

// find doc in user's data
func Find(uid string) (map[string]interface{}, error) {
	sel := bson.M{}
	doc, err := xxdb.Find(DB, "user", bson.M{"uid": uid}, sel)
	if err != nil {
		return nil, xxerr("Find", "", err)
	}
	return doc, nil
}

func FindBy(pth, val string) (map[string]interface{}, error) {
	sel := bson.M{}
	doc, err := xxdb.Find(DB, "user", bson.M{pth: val}, sel)
	if err != nil {
		return nil, xxerr("Find", "", err)
	}
	return doc, nil
}

// search user's infomation item
func Search(uid, pth string) (interface{}, error) {
	doc, err := Find(uid)
	if err != nil {
		return nil, xxerr("Search", "", err)
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
	err := xxdb.Upsert(DB, "user", bson.M{"uid": uid}, bson.M{pth: val})
	if err != nil {
		return xxerr("Upsert", "", err)
	}
	return nil
}

// increase user's infomation item
func Inc(uid, pth string, num float64) error {
	err := xxdb.Inc(DB, "user", bson.M{"uid": uid}, bson.M{pth: num})
	if err != nil {
		return xxerr("Inc", "", err)
	}
	return nil
}

// unset key in user data
func Unset(uid, pth string) error {
	err := xxdb.Unset(DB, "user", bson.M{"uid": uid}, bson.M{pth: 1})
	if err != nil {
		return xxerr("Unset", "", err)
	}
	return nil
}

// remove user's full data
func Remove(uid string) error {
	err := xxdb.Remove(DB, "user", bson.M{"uid": uid})
	if err != nil {
		return xxerr("Remove", "", err)
	}
	return nil
}

// insert element to array
func ArrayInsert(uid, name, pth string, content interface{}) error {
	var condition bson.M
	if uid != "" {
		condition = bson.M{"uid": uid}
	} else {
		condition = bson.M{"name": name}
	}
	err := xxdb.AddToSet(DB, "user", condition, bson.M{pth: content})
	if err != nil {
		return xxerr("ArrayInsert", "", err)
	}
	return nil
}

// remove element from array
func ArrayRemove(uid, name, pth string, content interface{}) error {
	var condition bson.M
	if uid != "" {
		condition = bson.M{"uid": uid}
	} else {
		condition = bson.M{"name": name}
	}
	err := xxdb.Pull(DB, "user", condition, bson.M{pth: content})
	if err != nil {
		return xxerr("ArrayRemove", "", err)
	}
	return nil
}

// implementation
// generate user's ID
func generateUid() (string, error) {
	num, err := getcount()
	if err != nil {
		xxerr("generateUid", "getcount", err)
		return "", err
	}
	num++
	err = setcount(num)
	if err != nil {
		xxerr("generateUid", "setcount", err)
		return "", err
	}
	num += 100000
	uid := strconv.Itoa(int(num))
	return uid, nil
}

func getcount() (float64, error) {
	doc, err := Find("xx")
	if err != nil {
		xxerr("getcount", "Find", err)
		return 0, err
	}
	if doc == nil {
		doc = map[string]interface{}{
			"uid": "xx", "count": 0.0,
		}
		err = xxdb.Insert("YBH", "user", doc)
		if err != nil {
			xxerr("getcount", "Insert", err)
			return 0, err
		}
	}
	ok, count := xx.Getnumber(doc, "count")
	if !ok {
		err = xxerr("getcount", "Getnumber", err)
		return 0, err
	}
	return count, nil
}

func setcount(num float64) error {
	err := Upsert("xx", "count", num)
	if err != nil {
		xxerr("setcount", "", err)
	}
	return err
}

func generateKey(openid string) (string, error) {
	data := md5.Sum([]byte(openid))
	base := data[:]
	salt := make([]byte, 12)
	_, err := rand.Read(salt)
	if err != nil {
		xxerr("generateKey", "", err)
		return "", err
	}
	s1 := hex.EncodeToString(base)
	s2 := hex.EncodeToString(salt)
	key := s1 + s2
	return key, nil
}

func xxerr(fn, info string, err error) error {
	fn = "user " + fn + "(): "
	return xxio.Error(true, fn, info, err)
}
