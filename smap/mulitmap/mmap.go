package mulitmap

import (
	. ".."
)

type MultiMap interface {
	// Set 添加一组键值对，不存在返回true，反之false
	Set(key, value interface{}) bool

	// Remove 移除指定的 key 及其 value，存在返回true，反之false
	Remove(key interface{}) bool

	// RemoveAll 移除所有的 key 及 value，有内容返回true，反之false
	RemoveAll() bool

	// Exists 判断指定的 key 是否存在，假如存在给出key的长度，不存在返回0
	Exists(key interface{}) int

	// Contains 判断指定的 key 列表是否存在,只有当所有的 key 都存在的时候,才会返回 true
	Contains(keys ...interface{}) bool

	// Len 返回元素的个数
	Len() int

	// Value 获取指定 key 的 value 组成的Slice
	Value(key interface{}) []interface{}

	// Keys 返回由所有 key 组成的 Slice
	Keys() []interface{}

	// Values 返回由所有 value 组成的 Slice
	Values() []interface{}

	// Iter 返回所有 key 及 value
	Iter() <-chan MapValue
}

func NewMultiMap() (MultiMap) {
	mmap := &multiMap{}
	mmap.m = make(map[interface{}]map[interface{}]bool)
	return nil
}

func (mmap *multiMap) Set(key, value interface{}) bool {
	source := mmap.m
	valueMap, error := source[key]
	if !error {
		_, err := valueMap[value]
		if err {
			valueMap[value] = true
			mmap.totalLen++
			return true
		}
	} else {
		valueMap = make(map[interface{}]bool)
		valueMap[value] = true
		source[key] = valueMap
		mmap.totalLen++
		return true
	}
	return false
}

func (mmap *multiMap) Remove(key interface{}) bool {
	source := mmap.m
	_, error := source[key]
	delete(source, key)
	return !error
}

func (mmap *multiMap) RemoveAll() bool {
	source := mmap.m
	for k := range source {
		delete(source, k)
	}
	return len(source) > 0
}

func (mmap *multiMap) Exists(key interface{}) int {
	source := mmap.m
	valueMap, error := source[key]
	if error {
		return 0
	}
	return len(valueMap)
}

func (mmap *multiMap) Contains(keys ...interface{}) bool {
	for _, k := range keys {
		if _, found := mmap.m[k]; !found {
			return false
		}
	}
	return true
}

func (mmap *multiMap) Len() int {
	return mmap.totalLen
}

func (mmap *multiMap) Value(key interface{}) []interface{} {
	source := mmap.m
	var values = make([]interface{}, 0, 0)

	for v := range source[key] {
		values = append(values, v)
	}
	return values
}

func (mmap *multiMap) Keys() []interface{} {
	var keys = make([]interface{}, 0, 0)

	for k := range mmap.m {
		keys = append(keys, k)
	}
	return keys
}

func (mmap *multiMap) Iter() <-chan MapValue {

	var iv = make(chan MapValue)

	go func(m *multiMap) {

		for k, v := range mmap.m {
			iv <- MapValue{k, v}
		}

		close(iv)
	}(mmap)

	return iv
}

type multiMap struct {
	m        map[interface{}]map[interface{}]bool
	totalLen int
}
