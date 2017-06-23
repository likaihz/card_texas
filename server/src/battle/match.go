package battle

import (
	"log"
	"math/rand"
	"strconv"
	"sync"
)

type Match struct {
	sync.Mutex
	max     int
	battles map[string]*Battle
}

func NewMatch(num int) *Match {
	m := new(Match)
	m.max = num
	m.battles = map[string]*Battle{}
	return m
}

//创建组，返回组编号
func (m *Match) Create() *Battle {
	m.Lock()
	defer m.Unlock()
	if len(m.battles) > m.max {
		log.Println("Create() the match is full!")
		return nil
	}
	s := m.number()
	b := New(m, s)
	if b == nil {
		log.Println("Create() failed!")
	} else {
		m.battles[s] = b
	}
	return b
}

// 索引房号为 s 的房间
func (m *Match) Get(s string) *Battle {
	return m.battles[s]
}

// 删除房号为 s 的房间
func (m *Match) Remove(s string) {
	delete(m.battles, s)
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

func (m *Match) number() string {
	for {
		n := rand.Intn(900000) + 100000
		s := strconv.Itoa(n)
		_, ok := m.battles[s]
		if !ok {
			return s
		}
	}
	return ""
}
