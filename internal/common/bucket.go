package common

import "runtime"

func NewBucket(freeSize int) *Bucket {
	return &Bucket{
		link: make(map[string]int),
		free: make([]int, 0, freeSize),
		data: make([]interface{}, 0, freeSize),
	}
}

type Bucket struct {
	link map[string]int
	free []int
	data []interface{}
}

func (bt *Bucket) Get(id string) interface{} {
	if i, ok := bt.link[id]; ok {
		return bt.data[i]
	}
	return nil
}

func (bt *Bucket) Set(id string, value interface{}) {
	fl := len(bt.free)
	if fl != 0 {
		bt.data[fl-1] = value
		bt.free = bt.free[:fl-1]
	}
	bt.data = append(bt.data, value)
	bt.link[id] = len(bt.data) - 1
}

func (bt *Bucket) Remove(id string) {
	bt.checkResize()
	if i, ok := bt.link[id]; ok {
		delete(bt.link, id)
		bt.data[i] = nil
		bt.free = append(bt.free, i)
	}
}

func (bt *Bucket) checkResize() {
	if len(bt.free) == cap(bt.free) {
		i := 0
		for _, v := range bt.data {
			if v == nil {
				bt.data[i] = v
				i++
			}
		}
		bt.data = bt.data[:i]
		bt.free = []int{}
		runtime.GC()
	}
}
