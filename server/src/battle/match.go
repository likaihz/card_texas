package battle

import (
	"math/rand"
	"sync"
)

type Match struct {
	sync.Mutex
	max     int
	battles map[int]*Battle
}

func NewMatch(num int) *Match {
	m := new(Match)
	m.max = num
	m.battles = map[int]*Battle{}
	return m
}

//创建组，返回组编号
func (m *Match) Create() *Battle {
	m.Lock()
	defer m.Unlock()
	if len(m.battles) > m.max {
		return nil
	}
	n := m.number()
	b := New(m, n)
	m.battles[n] = b
	return b
}

// 索引编号为 n 的组
func (m *Match) Get(n int) *Battle {
	return m.battles[n]
}

// 删除编号为 n 的组
func (m *Match) Remove(n int) {
	delete(m.battles, n)
}

// 随机匹配可进入的组
func (m *Match) Pop() *Battle {
	m.Lock()
	defer m.Unlock()
	for _, b := range m.battles {
		if b.Enterable() {
			return b
		}
	}
	return nil
}

func (m *Match) number() int {
	for {
		n := rand.Intn(900000) + 100000
		_, ok := m.battles[n]
		if !ok {
			return n
		}
	}
	return -1
}
