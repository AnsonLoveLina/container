package set

import (
	"fmt"
	"sync"
)

type Set interface {
	// 向集合添加元素
	Add(values ...interface{})

	// 从集合移除指定的元素
	Remove(values ...interface{})

	// 从集合移除所有的元素
	RemoveAll()

	// 判断集合是否包含指定元素
	Exists(v interface{}) bool

	// 判断集合是否包含指定的元素,包含所有的元素才会返回 true, 否则返回 false
	Contains(values ...interface{}) bool

	// 返回集合元素的长度
	Len() int

	// 返回集合元素组成的 Slice
	Values() []interface{}

	// 交集
	Intersect(s Set) Set

	// 并集
	Union(s Set) Set

	// 差集
	Difference(s Set) Set
}

type set struct {
	m     map[interface{}]struct{}
	rw    sync.RWMutex
	block bool
}

func NewSet(values ...interface{}) Set {
	return newSet(false, values...)
}

func NewBlockSet(values ...interface{}) Set {
	return newSet(true, values...)
}

func newSet(block bool, values ...interface{}) Set {
	var s = &set{}
	s.block = block
	s.m = make(map[interface{}]struct{})
	if len(values) > 0 {
		s.Add(values...)
	}
	return s
}

func (this *set) lock() {
	if this.block {
		this.rw.Lock()
	}
}

func (this *set) unlock() {
	if this.block {
		this.rw.Unlock()
	}
}

func (this *set) rLock() {
	if this.block {
		this.rw.RLock()
	}
}

func (this *set) rUnlock() {
	if this.block {
		this.rw.RUnlock()
	}
}

func (this *set) Add(values ...interface{}) {
	this.lock()
	defer this.unlock()

	for _, v := range values {
		this.m[v] = struct{}{}
	}
}

func (this *set) Remove(values ...interface{}) {
	this.lock()
	defer this.unlock()

	for _, v := range values {
		delete(this.m, v)
	}
}

func (this *set) RemoveAll() {
	this.lock()
	defer this.unlock()

	for k, _ := range this.m {
		delete(this.m, k)
	}
}

func (this *set) Exists(v interface{}) bool {
	this.rLock()
	defer this.rUnlock()

	_, found := this.m[v]
	return found
}

func (this *set) Contains(values ...interface{}) bool {
	this.rLock()
	defer this.rUnlock()

	for _, v := range values {
		if _, found := this.m[v]; !found {
			return false
		}
	}
	return true
}

func (this *set) Len() int {
	this.rLock()
	defer this.rUnlock()

	return this.len()
}

func (this *set) len() int {
	return len(this.m)
}

func (this *set) Values() []interface{} {
	this.rLock()
	defer this.rUnlock()

	var ns = make([]interface{}, 0, this.len())
	for k, _ := range this.m {
		ns = append(ns, k)
	}
	return ns
}

func (this *set) Intersect(s Set) Set {
	this.rLock()
	defer this.rUnlock()

	var ns = NewSet()
	var vs = s.Values()
	for _, v := range vs {
		_, exists := this.m[v]
		if exists {
			ns.Add(v)
		}
	}
	return ns
}

func (this *set) Union(s Set) Set {
	this.rLock()
	defer this.rUnlock()

	var ns = NewSet()
	ns.Add(this.Values()...)
	ns.Add(s.Values()...)
	return ns
}

func (this *set) Difference(s Set) Set {
	this.rLock()
	defer this.rUnlock()

	var ns = NewSet()
	for k, _ := range this.m {
		if !s.Contains(k) {
			ns.Add(k)
		}
	}
	return ns
}

func (this *set) String() string {
	return fmt.Sprint(this.Values()...)
}