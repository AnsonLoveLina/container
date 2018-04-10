package smap

import (
	"fmt"
	"sync"
)

type Map interface {
	// Set 添加一组键值对
	Set(key, value interface{})

	// Remove 移除指定的 key 及其 value
	Remove(key interface{})

	// RemoveAll 移除所有的 key 及 value
	RemoveAll()

	// Exists 判断指定的 key 是否存在
	Exists(key interface{}) bool

	// Contains 判断指定的 key 列表是否存在,只有当所有的 key 都存在的时候,才会返回 true
	Contains(keys ...interface{}) bool

	// Len 返回元素的个数
	Len() int

	// Value 获取指定 key 的 value
	Value(key interface{}) interface{}

	// Keys 返回由所有 key 组成的 Slice
	Keys() []interface{}

	// Values 返回由所有 value 组成的 Slice
	Values() []interface{}

	// Iter 返回所有 key 及 value
	Iter() <-chan MapValue
}

type MapValue struct {
	Key   interface{}
	Value interface{}
}

type SyncMap struct {
	M     map[interface{}]interface{}
	Mu    sync.RWMutex
	Block bool
}

func New(block bool) Map {
	var sm = &SyncMap{}
	sm.Block = block
	sm.M = make(map[interface{}]interface{})
	return sm
}

func (this *SyncMap) lock() {
	if this.Block {
		this.Mu.Lock()
	}
}

func (this *SyncMap) unlock() {
	if this.Block {
		this.Mu.Unlock()
	}
}

func (this *SyncMap) rLock() {
	if this.Block {
		this.Mu.RLock()
	}
}

func (this *SyncMap) rUnlock() {
	if this.Block {
		this.Mu.RUnlock()
	}
}

func (this *SyncMap) Set(key, value interface{}) {
	this.lock()
	defer this.unlock()

	this.M[key] = value
}

func (this *SyncMap) Remove(key interface{}) {
	this.rLock()
	defer this.rUnlock()

	delete(this.M, key)
}

func (this *SyncMap) RemoveAll() {
	this.lock()
	defer this.unlock()

	for k := range this.M {
		delete(this.M, k)
	}
}

func (this *SyncMap) Exists(key interface{}) bool {
	this.rLock()
	defer this.rUnlock()

	_, found := this.M[key]
	return found
}

func (this *SyncMap) Contains(keys ...interface{}) bool {
	this.rLock()
	defer this.rUnlock()

	for _, k := range keys {
		if _, found := this.M[k]; !found {
			return false
		}
	}
	return true
}

func (this *SyncMap) Len() int {
	this.rLock()
	defer this.rUnlock()

	return this.len()
}

func (this *SyncMap) len() int {
	return len(this.M)
}

func (this *SyncMap) Value(key interface{}) interface{} {
	this.rLock()
	defer this.rUnlock()

	return this.M[key]
}

func (this *SyncMap) Keys() []interface{} {
	this.rLock()
	defer this.rUnlock()
	var keys = make([]interface{}, 0, 0)

	for k := range this.M {
		keys = append(keys, k)
	}
	return keys
}

func (this *SyncMap) Values() []interface{} {
	this.rLock()
	defer this.rUnlock()
	var values = make([]interface{}, 0, 0)

	for _, v := range this.M {
		values = append(values, v)
	}
	return values
}

func (this *SyncMap) Iter() <-chan MapValue {
	var iv = make(chan MapValue)

	go func(m *SyncMap) {
		if m.Block {
			m.rLock()
		}

		for k, v := range this.M {
			iv <- MapValue{k, v}
		}

		close(iv)

		if m.Block {
			m.rUnlock()
		}
	}(this)

	return iv
}

func (this *SyncMap) String() string {
	return fmt.Sprint(this.M)
}
