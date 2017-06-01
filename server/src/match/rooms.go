package match

import (
	"../lib/xx"
	"log"
	"os"
	"sync"
)

type Rooms struct {
	sync.Mutex
	arr  []*Room
	log  *log.Logger
	file *os.File
}

func NewRooms(num int) *Rooms {
	r := new(Rooms)
	r.arr = make([]*Room, num)
	r.log, r.file = xx.Log2file("Rooms.log")
	return r
}

//创建组，返回组编号
func (r *Rooms) Create() *Room {
	r.Lock()
	defer r.Unlock()
	for i, e := range r.arr {
		if e == nil {
			n := r.number(i)
			e = NewRoom(n)
			r.arr[i] = e
			return e
		}
	}
	return nil
}

// 索引组号为 n 的组
func (r *Rooms) Get(n int) (*Room, bool) {
	r.Lock()
	defer r.Unlock()
	i := r.index(n)
	if i < 0 || i >= len(r.arr) {
		return nil, false
	}
	return r.arr[i], true
}

// 删除编号为 i 的组
func (r *Rooms) Remove(e *Room) {
	r.Lock()
	defer r.Unlock()
	n := e.Number()
	i := r.index(n)
	r.arr[i] = nil
}

// 随机匹配可进入的组
func (r *Rooms) Pop() *Room {
	for _, e := range r.arr {
		if e.Enterable() {
			return e
		}
	}
	return nil
}

// implementation
func (r *Rooms) index(num int) int {
	return num - 100000
}

func (r *Rooms) number(idx int) int {
	return idx + 100000
}
