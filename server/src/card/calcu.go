package card

// implementation of cards
func (c *Cards) calrank() int {
	straight := c.straight()
	flush := c.flush()
	if straight && flush {
		if c.Max().Idx() == 14 {
			return 10
		} else {
			return 9
		}
		return c.rank
	}
	_, most := c.Most()
	if most == 4 { //四条
		return 8
	}

	if c.fullhouse() {
		return 7
	}
	if flush {
		return 6
	}
	if straight {
		return 5
	}
	if most == 3 {
		return 4
	}
	pairnum := c.pairs()
	if pairnum == 2 {
		return 3
	}
	if pairnum == 1 {
		return 2

	}
	return 1
}

// 一对，两对可以一起处理，该函数返回的是一手牌中对子的数目
func (c *Cards) pairs() int {
	tbl := map[int]int{}
	cnt := 0
	for _, crd := range c.arr {
		j := crd.Idx()
		tbl[j]++
	}
	for _, v := range tbl {
		if v == 2 {
			cnt++
		}
	}
	return cnt

}

// 判断是否顺子
func (c *Cards) straight() bool {
	count := map[int]int{}
	max := 0
	for _, e := range c.arr {
		idx := e.Idx()
		count[idx]++
		if max < idx {
			max = idx
		}
	}
	num := len(c.arr)
	if max < num {
		return false
	}
	//处理特殊情况{A, 2, 3, 4, 5}
	if max == 14 && count[13] == 0 {
		count[1], count[14] = count[14], 0
		max = 5

	}
	for i := max; i > max-num; i-- {
		if count[i] != 1 {
			return false
		}
	}

	//这里还要再处理一次特殊情况
	//有一点点小问题，直接改idx会不会对后面产生影响，如果有问题的话可以考虑直接把顺子的最大值返回给调用者
	if max == 5 {
		for i, _ := range c.arr {
			if c.arr[i].Idx() == 14 {
				c.arr[i].Setidx(1)
				break
			}
		}
	}
	return true
}

// 判断是否同花
func (c *Cards) flush() bool {
	cls := c.arr[0].Cls()
	for _, e := range c.arr {
		if cls != e.Cls() {
			return false
		}
	}
	return true
}

// 判断是否葫芦
func (c *Cards) fullhouse() bool {
	tbl := map[int]int{}
	for _, e := range c.arr {
		tbl[e.Idx()]++
	}
	for _, v := range tbl {
		if v != 3 && v != 2 {
			return false
		}
	}
	return true
}

// 判断是否炸弹
// func (c *Cards) bomb() bool {
// 	_, num := c.Most()
// 	return num >= 4
// }

// 从多张牌中找出最大的五张牌组合
func CombinationTraversal(c []*Card) *Cards {
	n := len(c)
	if n < 5 {
		return nil
	}
	max := NewCards()
	flag := []bool{}
	for i := 0; i < n; i++ {
		flag = append(flag, false)
	}
	for i := 0; i < 5; i++ {
		max.Append(c[i])
		flag[i] = true
	}
	for i := 0; i < n-1; i++ {
		if flag[i] && !flag[i+1] {
			flag[i], flag[i+1] = false, true
			tmp := NewCards()
			for j, v := range flag {
				if v {
					tmp.Append(c[j])
				}
			}
			if max.Compare(tmp) < 0 {
				max = tmp
			}
		}
	}
	return max

}
