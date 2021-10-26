package utils

import "sync"

type BeeMap struct {
	Lock *sync.RWMutex
	BM   map[string]interface{}
}

// BewBeeMap 新建Map集合
func BewBeeMap() *BeeMap {
	return &BeeMap{
		Lock: new(sync.RWMutex),
		BM:   make(map[string]interface{}),
	}
}

// Get 通过key获取value
func (m *BeeMap) Get(k string) interface{} {
	m.Lock.RLock()
	defer m.Lock.RUnlock()
	if val, ok := m.BM[k]; ok {
		return val
	}
	return nil
}

// Set 不覆盖存元素
func (m *BeeMap) Set(k string, v interface{}) bool {
	m.Lock.RLock()
	defer m.Lock.RUnlock()
	if val, ok := m.BM[k]; !ok {
		m.BM[k] = v
	} else if val != v {
		m.BM[k] = v
	} else {
		return false
	}
	return true
}

// ReSet 覆盖存元素
func (m *BeeMap) ReSet(k string, v interface{}) bool {
	m.Delete(k)
	m.Set(k, v)
	return true
}

// Check 判断是否存在该key
func (m *BeeMap) Check(k string) bool {
	m.Lock.RLock()
	defer m.Lock.RUnlock()
	if _, ok := m.BM[k]; !ok {
		return false
	}
	return true
}

// Delete 通过key删除元素
func (m *BeeMap) Delete(k string) {
	m.Lock.RLock()
	defer m.Lock.RUnlock()
	delete(m.BM, k)
}

// Size 获取元素个数
func (m *BeeMap) Size() int {
	m.Lock.RLock()
	defer m.Lock.RUnlock()

	return len(m.BM)
}

// GetFirst 只读第一个
func (m *BeeMap) GetFirst() interface{} {
	m.Lock.RLock()
	defer m.Lock.RUnlock()

	for _, v := range m.BM {
		if v != nil {
			return v
		}
	}

	return nil
}

// DetachFirst 返回第一个，且从map中删除
func (m *BeeMap) DetachFirst() (string, interface{}) {
	m.Lock.Lock()
	defer m.Lock.Unlock()
	for k, v := range m.BM {
		if v != nil {
			delete(m.BM, k)
			return k, v
		}
	}
	return "", nil
}
