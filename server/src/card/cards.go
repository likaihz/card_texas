package card

import (
	// "fmt"
	"math/rand"
	"strconv"
)

// class Cards
type Cards struct {
	arr  []*Card
	rank int
}

func NewCards() *Cards {
	c := &Cards{}
	c.arr = make([]*Card, 0)
	c.rank = -1
	return c
}

// interface of Cards
func (c *Cards) Msg() map[string]interface{} {
	if len(c.arr) == 0 {
		return nil
	}
	msg := map[string]interface{}{}
	for i, e := range c.arr {
		s := strconv.Itoa(i)
		msg[s] = e.Msg()
	}
	return msg
}

func (c *Cards) Card(i int) *Card {
	return c.arr[i]
}

func (c *Cards) Insert(crd *Card) {
	j := len(c.arr)
	for i, e := range c.arr {
		if crd.Compare(e) <= 0 {
			j = i
			break
		}
	}
	c.arr = Insert(c.arr, j, crd)
}

func (c *Cards) Append(crd *Card) {
	c.arr = append(c.arr, crd)
}

func (c *Cards) Remove(crd *Card) *Card {
	j := -1
	for i, e := range c.arr {
		if crd.Same(e) {
			j = i
		}
	}
	if j < 0 {
		return nil
	}
	e := c.arr[j]
	c.arr = Remove(c.arr, j)
	return e
}

func (c *Cards) Rank() int {
	if len(c.arr) == 0 {
		return -1
	}
	if c.rank >= 0 {
		return c.rank
	}
	c.rank = c.calrank()
	return c.rank
}

func (c *Cards) Compare(cards *Cards) int {
	r := c.Rank()
	rank := cards.Rank()

	switch {
	case r > rank:
		return 1
	case r < rank:
		return -1
	case r == 10: //皇家同花顺
		return 0
	}

	// r == rank
	switch rank {
	case 5, 9: // 顺子，同花顺
		m := c.Max()
		max := cards.Max()
		return m.Compare(max)
	case 2, 3, 4, 7, 8: // 一对，两对，三条，葫芦，四条
		return c.Mergecompare(cards)
	}
	//同花，高牌
	return c.Rudecompare(cards)
}

// 从多张牌中找出最大的五张牌组合
func (c *Cards) CombinationTraversal() *Cards {
	n := len(c.arr)
	if n < 5 {
		return nil
	}
	max := NewCards()
	flag := []bool{}
	for i := 0; i < n; i++ {
		flag[i] = false
	}
	for i := 0; i < 5; i++ {
		max.Append(c.arr[i])
		flag[i] = true
	}
	for i := 0; i < n-1; i++ {
		if flag[i] && !flag[i+1] {
			flag[i], flag[i+1] = false, true
			tmp := NewCards()
			for j, v := range flag {
				if v {
					tmp.Append(c.arr[j])
				}
			}
			if max.Compare(tmp) < 0 {
				max = tmp
			}
		}
	}
	return max

}

// func (c *Cards) Price(mode string) int {
// 	r := c.Rank()
// 	if r == 0 {
// 		return 1
// 	}
// 	if mode == "doubling" {
// 		return r
// 	}
// 	switch {
// 	case r <= 6:
// 		return 1
// 	case r > 6 && r < 10:
// 		return 2
// 	case r == 10:
// 		return 3
// 	default:
// 		return 5
// 	}
// }

// interface of function
func Insert(arr []*Card, i int, c *Card) []*Card {
	l := len(arr)
	arr = append(arr, c)
	if i == l {
		return arr
	}
	copy(arr[i+1:], arr[i:l])
	arr[i] = c
	return arr
}

func Remove(arr []*Card, i int) []*Card {
	copy(arr[i:], arr[i+1:])
	return arr[:len(arr)-1]
}

func Pop(arr []*Card) ([]*Card, *Card) {
	num := len(arr)
	if num <= 0 {
		return arr, nil
	}
	i := rand.Intn(num)
	c := arr[i]
	arr = Remove(arr, i)
	return arr, c
}

// implementation
func (c *Cards) Rudecompare(cards *Cards) int {
	t1 := NewCards()
	t2 := NewCards()
	num := len(c.arr)
	for i := 0; i < num; i++ {
		t1.Insert(c.arr[i])
		crd := cards.Card(i)
		t2.Insert(crd)
	}
	for i := num - 1; i >= 0; i-- {
		c1, c2 := t1.Card(i), t2.Card(i)
		idx1, idx2 := c1.Idx(), c2.Idx()
		if idx1 < idx2 {
			return -1
		} else if idx1 > idx2 {
			return 1
		}
	}
	return 0
}

func (c *Cards) Mergecompare(cards *Cards) int {
	c1 := NewCards()
	c2 := NewCards()

	tbl := map[int]int{}
	for _, crd := range c.arr {
		j := crd.Idx()
		tbl[j]++
	}
	for k, v := range tbl {
		tmp := New("", k+(v-1)*100)
		c1.Append(tmp)
	}

	tbl = map[int]int{}
	for _, crd := range cards.arr {
		j := crd.Idx()
		tbl[j]++
	}
	for k, v := range tbl {
		tmp := New("", k+(v-1)*100)
		c2.Append(tmp)
	}
	return c1.Rudecompare(c2)
}

func (c *Cards) Max() *Card {
	num := len(c.arr)
	max := c.arr[0]
	for i := 1; i < num; i++ {
		tmp := c.arr[i]
		if max.Compare(tmp) < 0 {
			max = tmp
		}
	}
	return max
}

func (c *Cards) Most() (int, int) {
	tbl := map[int]int{}
	idx, cnt := 0, 0
	for _, crd := range c.arr {
		j := crd.Idx()
		tbl[j]++
		if cnt < tbl[j] {
			idx, cnt = j, tbl[j]
		}
	}
	return idx, cnt
}

func (c *Cards) Sum() int {
	sum := 0
	for _, e := range c.arr {
		sum += e.Val()
	}
	return sum
}

func (c *Cards) Amount(crd *Card) int {
	cnt := 0
	for _, e := range c.arr {
		if crd.Same(e) {
			cnt++
		}
	}
	return cnt
}

// implementation
