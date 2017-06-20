package xxio

import (
	"../xxtea"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// create error infomation
func Error(out bool, fn, info string, err error) error {
	txt := fn
	if info != "" {
		txt += info + ": "
	}
	err = fmt.Errorf(txt, err)
	if out {
		fmt.Println(err)
	}
	return err
}

func Decode(r io.Reader) (map[string]interface{}, error) {
	body, err := ioutil.ReadAll(r)
	if err != nil {
		fmt.Println("reading err: ", err)
		return nil, err
	}
	data := map[string]interface{}{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("decoding err: ", err)
		return nil, err
	}
	return data, nil
}

func Request(r *http.Request, crypt bool) (map[string]interface{}, error) {
	body, err := ioutil.ReadAll(r.Body)
	fmt.Println(string(body))
	defer r.Body.Close()
	if err != nil {
		return nil, xxerr("Request", "read", err)
	}
	if crypt {
		body = xxtea.Decrypt(body, []byte("玄襄科技"))
	}
	data := map[string]interface{}{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, xxerr("Request", "unmarshal", err)
	}
	return data, nil
}

// response net infomation
func Response(w http.ResponseWriter, res *map[string]interface{}, crypt bool) error {
	data, err := json.Marshal(*res)
	if err != nil {
		return xxerr("Response", "marshal", err)
	}
	if crypt {
		data = xxtea.Encrypt(data, []byte("玄襄科技"))
	}
	w.Write(data)
	return nil
}

// read config file in server
func Read(name string) (map[string]interface{}, error) {
	info := "file " + name
	file, err := ioutil.ReadFile(getpath(name))
	if err != nil {
		return nil, xxerr("Read", "read "+info, err)
	}
	var data map[string]interface{}
	err = json.Unmarshal(file, &data)
	if err != nil {
		return nil, xxerr("Read", "unmarshal "+info, err)
	}
	return data, nil
}

// update config file to client
func Update(client map[string]interface{}, res *map[string]interface{}) error {
	server, err := Read("version")
	if err != nil {
		return xxerr("Update", "", err)
	}
	tbl := *res
	var doc map[string]interface{}
	vsn := map[string]string{}
	for name, val := range server {
		svr := val.(string)
		clt := client[name]
		if clt == nil || clt.(string) != svr {
			doc, err = Read(name)
			if err == nil {
				tbl[name] = doc
				vsn[name] = svr
			}
		}
	}
	tbl["version"] = vsn
	return nil
}

// -- implementation --
func getpath(name string) string {
	return "./data/" + name + ".json"
}

func xxerr(fn, info string, err error) error {
	fn = "xxio." + fn
	return Error(true, fn, info, err)
}
