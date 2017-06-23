package mongo

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

const URL = "127.0.0.1:27017"

// dial to db
func Dial(db, collection string) (*mgo.Session, *mgo.Collection, error) {
	var ses *mgo.Session
	var col *mgo.Collection
	var err error
	if ses, err = mgo.Dial(URL); err != nil {
		log.Println("Dial(): ", err)
		return nil, nil, err
	}
	ses.SetMode(mgo.Monotonic, true)
	col = ses.DB(db).C(collection)
	return ses, col, nil
}

// find a document in db
func Find(db, collection string, condition, selection bson.M) (map[string]interface{}, error) {
	ses, col, err := Dial(db, collection)
	defer ses.Close()
	if err != nil {
		log.Println("Find(): ", err)
		return nil, err
	}
	doc := map[string]interface{}{}
	if len(selection) == 0 {
		err = col.Find(condition).One(&doc)
	} else {
		err = col.Find(condition).Select(selection).One(&doc)
	}
	if err != nil {
		return nil, nil
	}
	return doc, nil
}

// insert a document to db
func Insert(db, collection string, doc map[string]interface{}) error {
	ses, col, err := Dial(db, collection)
	defer ses.Close()
	if err != nil {
		log.Println("Insert(): ", err)
		return err
	}
	err = col.Insert(doc)
	if err != nil {
		log.Println("Insert(): ", err)
		return err
	}
	return nil
}

// update content of a document
// content: {path("xx.xx.xx..."): val}
func Upsert(db, collection string, condition, content bson.M) error {
	ses, col, err := Dial(db, collection)
	defer ses.Close()
	if err != nil {
		log.Println("Upsert(): ", err)
		return err
	}
	_, err = col.Upsert(condition, bson.M{"$set": content})
	if err != nil {
		log.Println("Upsert(): ", err)
		return err
	}
	return nil
}

func UpdateAll(db, collection string, condition, content bson.M) error {
	ses, col, err := Dial(db, collection)
	defer ses.Close()
	if err != nil {
		log.Println("UpdateAll(): ", err)
		return err
	}
	_, err = col.UpdateAll(condition, content)
	if err != nil {
		log.Println("UpdateAll(): ", err)
		return err
	}
	return nil
}

// content: {path("xx.xx.xx..."): val}
func Inc(db, collection string, condition, content bson.M) error {
	ses, col, err := Dial(db, collection)
	defer ses.Close()
	if err != nil {
		log.Println("Inc(): ", err)
		return err
	}
	_, err = col.Upsert(condition, bson.M{"$inc": content})
	if err != nil {
		log.Println("Inc(): ", err)
		return err
	}
	return nil
}

// unset a key match the selector
// content: {path("xx.xx.xx..."): 1}
func Unset(db, collection string, condition, content bson.M) error {
	ses, col, err := Dial(db, collection)
	defer ses.Close()
	if err != nil {
		log.Println("Unset(): ", err)
		return err
	}
	_, err = col.Upsert(condition, bson.M{"$unset": content})
	if err != nil {
		log.Println("Unset(): ", err)
		return err
	}
	return nil
}

// remove a document match the selector in db
func Remove(db, collection string, selector bson.M) error {
	ses, col, err := Dial(db, collection)
	defer ses.Close()
	if err != nil {
		log.Println("Remove(): ", err)
		return err
	}
	err = col.Remove(selector)
	if err != nil {
		log.Println("Remove(): ", err)
		return err
	}
	return nil
}

// remove all documents match the selector in db
func RemoveAll(db, collection string, selector bson.M) error {
	ses, col, err := Dial(db, collection)
	defer ses.Close()
	if err != nil {
		log.Println("RemoveAll(): ", err)
		return err
	}
	_, err = col.RemoveAll(selector)
	if err != nil {
		log.Println("RemoveAll(): ", err)
		return err
	}
	return nil
}

// 对数组在末尾加入一个元素，如果这个元素和数组中元素重复，则不添加
// 如果数组不存在，则新建数组
func AddToSet(db, collection string, condition, content bson.M) error {
	ses, col, err := Dial(db, collection)
	defer ses.Close()
	if err != nil {
		log.Println("AddToSet(): ", err)
		return err
	}
	_, err = col.Upsert(condition, bson.M{"$addToSet": content})
	if err != nil {
		log.Println("AddToSet(): ", err)
		return err
	}
	return nil
}

// 对数组在末尾加入一个元素，如果这个元素和数组中元素重复，依然添加
// 如果数组不存在，则新建数组
func Push(db, collection string, condition, content bson.M) error {
	ses, col, err := Dial(db, collection)
	defer ses.Close()
	if err != nil {
		log.Println("Push(): ", err)
		return err
	}
	_, err = col.Upsert(condition, bson.M{"$push": content})
	if err != nil {
		log.Println("Push(): ", err)
		return err
	}
	return nil
}

//从数组中删除一个元素，如果有重复的元素，全部删除
func Pull(db, collection string, condition, content bson.M) error {
	ses, col, err := Dial(db, collection)
	defer ses.Close()
	if err != nil {
		log.Println("Pull(): ", err)
		return err
	}
	_, err = col.Upsert(condition, bson.M{"$pull": content})
	if err != nil {
		log.Println("Pull(): ", err)
		return err
	}
	return nil
}
