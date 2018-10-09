// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
//

package gset

import (
	"fmt"
	"sync"
)

type InterfaceSet struct {
	mu sync.RWMutex
	m  map[interface{}]struct{}
}

func NewInterfaceSet() *InterfaceSet {
	return &InterfaceSet{m: make(map[interface{}]struct{})}
}

// 给定回调函数对原始内容进行遍历
func (this *InterfaceSet) Iterator(f func(v interface{})) {
	this.mu.RLock()
	for k, _ := range this.m {
		f(k)
	}
	this.mu.RUnlock()
}

// 添加
func (this *InterfaceSet) Add(item interface{}) *InterfaceSet {
	this.mu.Lock()
	this.m[item] = struct{}{}
	this.mu.Unlock()
	return this
}

// 批量添加
func (this *InterfaceSet) BatchAdd(items []interface{}) *InterfaceSet {
	this.mu.Lock()
	for _, item := range items {
		this.m[item] = struct{}{}
	}
	this.mu.Unlock()
	return this
}

// 键是否存在
func (this *InterfaceSet) Contains(item interface{}) bool {
	this.mu.RLock()
	_, exists := this.m[item]
	this.mu.RUnlock()
	return exists
}

// 删除键值对
func (this *InterfaceSet) Remove(key interface{}) {
	this.mu.Lock()
	delete(this.m, key)
	this.mu.Unlock()
}

// 大小
func (this *InterfaceSet) Size() int {
	this.mu.RLock()
	l := len(this.m)
	this.mu.RUnlock()
	return l
}

// 清空set
func (this *InterfaceSet) Clear() {
	this.mu.Lock()
	this.m = make(map[interface{}]struct{})
	this.mu.Unlock()
}

// 转换为数组
func (this *InterfaceSet) Slice() []interface{} {
	this.mu.RLock()
	i := 0
	ret := make([]interface{}, len(this.m))
	for item := range this.m {
		ret[i] = item
		i++
	}
	this.mu.RUnlock()
	return ret
}

// 转换为字符串
func (this *InterfaceSet) String() string {
	return fmt.Sprint(this.Slice())
}
