package test

import (
	// "../lib/xx"
	"../card"
	// "../lib/xxio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

var card_cls = []string{"diamond", "club", "heart", "spade"}

type MyCards struct {
	C1, C2, C3, C4, C5 string
	Rank               int
}

func (c *MyCards) SetVal(r int, c1, c2, c3, c4, c5 string) {
	c.C1, c.C2, c.C3, c.C4, c.C5 = c1, c2, c3, c4, c5
	c.Rank = r
}

func CardTest() {
	var op int
	fmt.Println("Options:")
	fmt.Println("1 : Write")
	fmt.Println("2 : Run Test")

	fmt.Scanln(&op)

	switch op {
	case 1:
		Write()
	case 2:
		Test()
	}
}

func Write() {
	var idx, r int
	var c1, c2, c3, c4, c5 string
	var mc MyCards
	var cfm string

	f, err := os.OpenFile("./data/cardtest.json", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
	if err != nil {
		fmt.Println("!!!OpenFile Error: ", err)
	}
	fmt.Printf("input index,  -1 to quit: ")
	fmt.Scanln(&idx)
	for i := idx; idx >= 0; i++ {
		fmt.Printf("input 5 cards:\n\t")
		fmt.Scanf("%s%s%s%s%s\n", &c1, &c2, &c3, &c4, &c5)
		fmt.Printf("input expected rank:\t")
		fmt.Scanf("%d\n", &r)

		fmt.Printf("Confirm your data? [Y/n]:  ")
		fmt.Scanln(&cfm)
		if cfm == "y" || cfm == "Y" || cfm == "" {
			mc.SetVal(r, c1, c2, c3, c4, c5)
			// str :=
			pre := []byte("\"" + strconv.Itoa(i) + "\":")
			f.Write(pre)
			buf, _ := json.MarshalIndent(mc, "", "\t")
			f.Write(buf)
			f.Write([]byte(",\n"))
		}

		fmt.Printf("input index,  -1 to quit: ")
		fmt.Scanln(&idx)
	}

	f.Close()

}

func Test() {
	var num int
	fmt.Printf("Input the num of test cases:  ")
	fmt.Scanf("%d\n", &num)
	data, err := Read("user")
	if err != nil {
		fmt.Println(err)
		return
	}
	var allright = true
	for i := 0; i < num; i++ {
		onecase := data[strconv.Itoa(i)].(map[string]interface{})
		cards := card.NewCards()
		for j := 1; j <= 5; j++ {
			cardidx := "C" + strconv.Itoa(j)
			val, _ := strconv.Atoi(onecase[cardidx].(string))
			cls := card_cls[val/100]
			idx := val % 100
			card := card.New(cls, idx)
			cards.Append(card)
		}

		expectedrank := int(onecase["Rank"].(float64))
		testresult := cards.Rank()
		fmt.Println(expectedrank, testresult)
		if expectedrank != testresult {
			fmt.Printf("Case #%d is wrong. \n The expected rank is %d, but the program result is %d\n", i, expectedrank, testresult)
			allright = false
		}
	}

	if allright {
		fmt.Println("All right!")
	}
}

func Read(name string) (map[string]interface{}, error) {
	file, err := ioutil.ReadFile("./data/cardtest.json")
	if err != nil {
		return nil, err
	}
	var data map[string]interface{}
	err = json.Unmarshal(file, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
