package battle

import (
	"sync"
)

type Match struct {
	sync.Mutex
	arr []*Battle
}

func NewMatch(num int) *Match {
	m := new(Match)
	m.arr = make([]*Battle, num)
	return m
}

//创建组，返回组编号
func (m *Match) Create() *Battle {
	m.Lock()
	defer m.Unlock()
	for i, b := range m.arr {
		if b == nil {
			b = New(m, number(i)) // param
			m.arr[i] = b
			return b
		}
	}
	return nil
}

// 索引编号为 n 的组
func (m *Match) Get(n int) *Battle {
	i := index(n)
	if i < 0 || i >= len(m.arr) {
		return nil
	}
	return m.arr[i]
}

// 删除编号为 n 的组
func (m *Match) Remove(n int) {
	i := index(n)
	if i < 0 || i >= len(m.arr) {
		return
	}
	m.arr[i] = nil
}

// 随机匹配可进入的组
func (m *Match) Pop() *Battle {
	m.Lock()
	defer m.Unlock()
	for _, b := range m.arr {
		if b.Enterable() {
			return b
		}
	}
	return nil
}

// implementation
func index(num int) int {
	return num - 100000
}

func number(idx int) int {
	return idx + 100000
}
