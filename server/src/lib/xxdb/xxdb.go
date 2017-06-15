package xxdb

import (
	"../xxio"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const URL = "127.0.0.1:27017"

// dial to db
func Dial(db, collection string) (*mgo.Session, *mgo.Collection, error) {
	var ses *mgo.Session
	var col *mgo.Collection
	var err error
	if ses, err = mgo.Dial(URL); err != nil {
		return nil, nil, xxerr("Dial", "", err)
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
		return nil, xxerr("Find", "", err)
	}
	var doc map[string]interface{}
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
		return xxerr("Insert", "", err)
	}
	err = col.Insert(doc)
	if err != nil {
		return xxerr("Insert", "", err)
	}
	return nil
}

// update content of a document
// content: {path("xx.xx.xx..."): val}
func Upsert(db, collection string, condition, content bson.M) error {
	ses, col, err := Dial(db, collection)
	defer ses.Close()
	if err != nil {
		return xxerr("Upsert", "", err)
	}
	_, err = col.Upsert(condition, bson.M{"$set": content})
	if err != nil {
		return xxerr("Upsert", "", err)
	}
	return nil
}

func UpdateAll(db, collection string, condition, content bson.M) error {
	ses, col, err := Dial(db, collection)
	defer ses.Close()
	if err != nil {
		return xxerr("Update", "", err)
	}
	_, err = col.UpdateAll(condition, content)
	if err != nil {
		return xxerr("Update", "", err)
	}
	return nil
}

// content: {path("xx.xx.xx..."): val}
func Inc(db, collection string, condition, content bson.M) error {
	ses, col, err := Dial(db, collection)
	defer ses.Close()
	if err != nil {
		return xxerr("Inc", "", err)
	}
	_, err = col.Upsert(condition, bson.M{"$inc": content})
	if err != nil {
		return xxerr("Inc", "", err)
	}
	return nil
}

// unset a key match the selector
// content: {path("xx.xx.xx..."): 1}
func Unset(db, collection string, condition, content bson.M) error {
	ses, col, err := Dial(db, collection)
	defer ses.Close()
	if err != nil {
		return xxerr("Unset", "", err)
	}
	_, err = col.Upsert(condition, bson.M{"$unset": content})
	if err != nil {
		return xxerr("Unset", "", err)
	}
	return nil
}

// remove a document match the selector in db
func Remove(db, collection string, selector bson.M) error {
	ses, col, err := Dial(db, collection)
	defer ses.Close()
	if err != nil {
		return xxerr("Remove", "", err)
	}
	err = col.Remove(selector)
	if err != nil {
		return xxerr("Remove", "", err)
	}
	return nil
}

// remove all documents match the selector in db
func RemoveAll(db, collection string, selector bson.M) error {
	ses, col, err := Dial(db, collection)
	defer ses.Close()
	if err != nil {
		return xxerr("RemoveAll", "", err)
	}
	_, err = col.RemoveAll(selector)
	if err != nil {
		return xxerr("RemoveAll", "", err)
	}
	return err
}

// read in n save local file
func Save(name, collection string) error {
	info := "config flie " + name
	ses, col, err := Dial("sylm", collection)
	defer ses.Close()

	if err != nil {
		return xxerr("Save", info, err)
	}
	// clean the db
	err = col.DropCollection()
	if err != nil {
		if err.Error() != "ns not found" {
			fmt.Println("not found")
			return xxerr("Save", info, err)
		}
	}
	data, err := xxio.Read(name)
	if err != nil {
		return xxerr("Save", info, err)
	}
	err = Insert("sylm", collection, data)
	if err != nil {
		return xxerr("Save", info, err)
	}
	return nil
}

//对数组在末尾加入一个元素，如果这个元素和数组中元素重复，则不添加。如果数组不存在，则新建数组
func AddToSet(db, collection string, condition, content bson.M) error {
	ses, col, err := Dial(db, collection)
	defer ses.Close()
	if err != nil {
		return xxerr("AddToSet", "", err)
	}
	_, err = col.Upsert(condition, bson.M{"$addToSet": content})
	if err != nil {
		return xxerr("AddToSet", "", err)
	}
	return nil
}

//对数组在末尾加入一个元素，如果这个元素和数组中元素重复，依然添加。如果数组不存在，则新建数组
func Push(db, collection string, condition, content bson.M) error {
	ses, col, err := Dial(db, collection)
	defer ses.Close()
	if err != nil {
		return xxerr("Push", "", err)
	}
	_, err = col.Upsert(condition, bson.M{"$push": content})
	if err != nil {
		return xxerr("Push", "", err)
	}
	return nil
}

//从数组中删除一个元素，如果有重复的元素，全部删除
func Pull(db, collection string, condition, content bson.M) error {
	ses, col, err := Dial(db, collection)
	defer ses.Close()
	if err != nil {
		return xxerr("Pull", "", err)
	}
	_, err = col.Upsert(condition, bson.M{"$pull": content})
	if err != nil {
		return xxerr("Pull", "", err)
	}
	return nil
}

//根据查找条件以切片形式返回多个mongodb中的文档
func FindMulti(db, collection string, condition, selection bson.M, limit int, sort string) ([]map[string]interface{}, error) {
	ses, col, err := Dial(db, collection)
	defer ses.Close()
	if err != nil {
		return nil, xxerr("FindMulti", "", err)
	}
	var iter *mgo.Iter
	if sort != "" {
		iter = col.Find(condition).Select(selection).Sort(sort).Limit(limit).Iter()
	} else {
		iter = col.Find(condition).Select(selection).Limit(limit).Iter()
	}

	// var result []interface{}
	var result []map[string]interface{}
	err = iter.All(&result)
	if err != nil {
		return nil, xxerr("FindMulti", "", err)
	}
	return result, nil
}

func xxerr(fn, info string, err error) error {
	fn = "xxdb." + fn
	return xxio.Error(false, fn, info, err)
}
