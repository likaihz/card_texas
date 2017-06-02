package card

// import (
// 	"fmt"
// 	"strconv"
// )

var card_cls = map[string]int{
	"heart": 3, "spade": 4, "diamond": 1, "club": 2,
}

// class Card
type Card struct {
	cls string
	idx int
}

func New(cls string, idx int) *Card {
	return &Card{cls: cls, idx: idx}
}

// interface of Card
func (c *Card) Msg() map[string]interface{} {
	msg := map[string]interface{}{}
	msg["cls"] = c.cls
	msg["idx"] = c.idx
	return msg
}

func (c *Card) Cls() string {
	return c.cls
}

func (c *Card) Idx() int {
	return c.idx
}

func (c *Card) Val() int {
	val := c.idx
	if val == 1 {
		val = 14
	}
	return val
}

func (c *Card) Ncls() int {
	return card_cls[c.cls]
}

func (c *Card) Same(c1 *Card) bool {
	same := c.cls == c1.Cls()
	same = same && c.idx == c1.Idx()
	return same
}

func (c *Card) Setidx(n int) {
	c.idx = n
}

// func (c *Card) CompareCls(c1 *Card) int {
// 	cls := c.Ncls()
// 	if cls < c1.Ncls() {
// 		return -1
// 	} else if cls > c1.Ncls() {
// 		return 1
// 	}
// 	if c.idx < c1.Idx() {
// 		return -1
// 	} else if c.idx > c1.Idx() {
// 		return 1
// 	}
// 	return 0
// }

// func (c *Card) CompareIdx(c1 *Card) int {
// 	if c.idx < c1.Idx() {
// 		return -1
// 	} else if c.idx > c1.Idx() {
// 		return 1
// 	}
// 	cls := c.Ncls()
// 	if cls < c1.Ncls() {
// 		return -1
// 	} else if cls > c1.Ncls() {
// 		return 1
// 	}
// 	return 0
// }

func (c *Card) Compare(c1 *Card) int {
	idx := c.idx
	idx1 := c1.Idx()
	if idx < idx1 {
		return -1
	} else if idx == idx1 {
		return 0
	} else {
		return 1
	}
}

func Init() []*Card {
	arr := make([]*Card, 0)
	num := 13
	for cls := range card_cls {
		for i := 1; i <= num; i++ {
			e := New(cls, i)
			arr = append(arr, e)
		}
	}
	return arr
}
