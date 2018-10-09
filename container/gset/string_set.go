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

type StringSet struct {
	mu sync.RWMutex
	m  map[string]struct{}
}

func NewStringSet() *StringSet {
	return &StringSet{m: make(map[string]struct{})}
}

// 给定回调函数对原始内容进行遍历
func (this *StringSet) Iterator(f func(v string)) {
	this.mu.RLock()
	for k, _ := range this.m {
		f(k)
	}
	this.mu.RUnlock()
}

// 设置键
func (this *StringSet) Add(item string) *StringSet {
	this.mu.Lock()
	this.m[item] = struct{}{}
	this.mu.Unlock()
	return this
}

// 批量添加设置键
func (this *StringSet) BatchAdd(items []string) *StringSet {
	this.mu.Lock()
	for _, item := range items {
		this.m[item] = struct{}{}
	}
	this.mu.Unlock()
	return this
}

// 键是否存在
func (this *StringSet) Contains(item string) bool {
	this.mu.RLock()
	_, exists := this.m[item]
	this.mu.RUnlock()
	return exists
}

// 删除键值对
func (this *StringSet) Remove(key string) {
	this.mu.Lock()
	delete(this.m, key)
	this.mu.Unlock()
}

// 大小
func (this *StringSet) Size() int {
	this.mu.RLock()
	l := len(this.m)
	this.mu.RUnlock()
	return l
}

// 清空set
func (this *StringSet) Clear() {
	this.mu.Lock()
	this.m = make(map[string]struct{})
	this.mu.Unlock()
}

// 转换为数组
func (this *StringSet) Slice() []string {
	this.mu.RLock()
	ret := make([]string, len(this.m))
	i := 0
	for item := range this.m {
		ret[i] = item
		i++
	}

	this.mu.RUnlock()
	return ret
}

// 转换为字符串
func (this *StringSet) String() string {
	return fmt.Sprint(this.Slice())
}
