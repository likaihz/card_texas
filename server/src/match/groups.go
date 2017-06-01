package match

import (
	"../lib/xx"
	"log"
	"os"
	"sync"
)

type Receptor interface {
	Receive(msg map[string]interface{})
	Send(msg map[string]interface{})
	Index() int
	SetConfig(cfg map[string]interface{})
	Occupancy() int
	Enterable() bool
	Enter(uid string, ch client)
	Leave(uid string, ch client)
}

type Groups struct {
	sync.Mutex
	arr  []Receptor
	log  *log.Logger
	file *os.File
}

func NewGroups(num int) *Groups {
	g := new(Groups)
	g.arr = make([]Receptor, num)
	g.log, g.file = xx.Log2file("groups.log")
	return g
}

//创建组，返回组编号
func (g *Groups) Create(uid string, ch chan string) Receptor {
	g.Lock()
	defer g.Unlock()
	for i, e := range g.arr {
		if e == nil {
			e = NewGroup(uid, ch, i)
			g.arr[i] = e
			return e
		}
	}
	return nil
}

// 索引编号为 i 的组
func (g *Groups) Get(i int) Receptor {
	g.Lock()
	defer g.Unlock()
	return g.arr[i]
}

// 删除编号为 i 的组
func (g *Groups) Remove(e Receptor) {
	g.Lock()
	defer g.Unlock()
	i := e.Index()
	g.arr[i] = nil
}

// 随机匹配可进入的组
func (g *Groups) Pop() Receptor {
	for _, e := range g.arr {
		if e.Enterable() {
			return e
		}
	}
	return nil
}
