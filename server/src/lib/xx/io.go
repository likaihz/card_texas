package xx

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

func Decode(r io.Reader) (map[string]interface{}, error) {
	body, err := ioutil.ReadAll(r)
	if err != nil {
		log.Println("Decode() ", err)
		return nil, err
	}
	data := map[string]interface{}{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Println("Decode() ", err)
		return nil, err
	}
	return data, nil
}

func Request(r *http.Request, crypt bool) (map[string]interface{}, error) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Println("Request() ", err)
		return nil, err
	}
	if crypt {
		body = Decrypt(body, []byte("玄襄科技"))
	}
	data := map[string]interface{}{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Println("Request() ", err)
		return nil, err
	}
	return data, nil
}

// response http request
func Response(w http.ResponseWriter, res *map[string]interface{}, crypt bool) error {
	data, err := json.Marshal(*res)
	if err != nil {
		log.Println("Response() ", err)
		return err
	}
	if crypt {
		data = Encrypt(data, []byte("玄襄科技"))
	}
	w.Write(data)
	return nil
}

// read file in directory "data"
func Read(name string) (map[string]interface{}, error) {
	pth := "./data/" + name + ".json"
	file, err := ioutil.ReadFile(pth)
	if err != nil {
		log.Printf("Read() %s err: %v\n", name, err)
		return nil, err
	}
	data := map[string]interface{}{}
	err = json.Unmarshal(file, &data)
	if err != nil {
		log.Printf("Read() %s err: %v\n", name, err)
		return nil, err
	}
	return data, nil
}
