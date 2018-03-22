package base

import "sync"


// 线程安全的map
type SafeMap struct {
	sync.RWMutex
	Map map[string]int64
}




func (sm *SafeMap) readMap(key string) int64 {
	sm.RLock()
	value := sm.Map[key]
	sm.RUnlock()
	return value
}

func (sm *SafeMap) add(key string,) {
	sm.Lock()
	sm.Map[key]++
	sm.Unlock()
}

func (sm *SafeMap) setValue(key string, value int64) {
	sm.Lock()
	sm.Map[key] = value
	sm.Unlock()
}

