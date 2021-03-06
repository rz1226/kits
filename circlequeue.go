package kits

import (
	"sync"
	"sync/atomic"
)

/*
cq := kits.NewCircleQueue(100)
cq.Put("abc")
lists, newestId := cq.GetSeveral(10)

*/

// 只保留最近写入的部分数据
type CircleQueue struct {
	currentID    uint64
	size         uint32
	dataNodeList []dataNode
}

type dataNode struct {
	id      uint64
	content interface{}
	mu      *sync.RWMutex
}

func NewCircleQueue(size uint32) *CircleQueue {
	cq := &CircleQueue{}
	cq.currentID = 0
	cq.size = minQuantity(size)
	cq.dataNodeList = make([]dataNode, cq.size)
	for k, _ := range cq.dataNodeList {
		ele := &(cq.dataNodeList[k])
		ele.id = 0
		ele.mu = &sync.RWMutex{}
	}
	return cq
}

// 把任何类型的数据放入队列
func (c *CircleQueue) Put(val interface{}) {
	// func AddUint64(addr *uint64, delta uint64) (new uint64)
	nextID := atomic.AddUint64(&c.currentID, 1)
	//  & 相当于取模
	positionInList := nextID & uint64((c.size - 1))
	dataNode := &(c.dataNodeList[positionInList])
	dataNode.mu.Lock()
	defer dataNode.mu.Unlock()
	dataNode.id = nextID
	dataNode.content = val
}

// 从队列中取出count个数的数据，返回当初放进去的数据列表，以及最新的id
func (c *CircleQueue) GetSeveral(count int) (data []interface{}, fetchCount uint64) {
	resDataSli := make([]interface{}, count)
	currentID := atomic.LoadUint64(&c.currentID)
	for i := 0; i < count-1; i++ {
		//  & 相当于取模
		pos := (currentID - uint64(i)) & uint64((c.size - 1))
		dataNode := &(c.dataNodeList[pos])
		dataNode.mu.RLock()
		if dataNode.id == currentID-uint64(i) {
			resDataSli[i] = dataNode.content
			dataNode.mu.RUnlock()
		} else {
			dataNode.mu.RUnlock()
			break
		}
	}
	return resDataSli, currentID
}

//  round 到最近的2的倍数
func minQuantity(v uint32) uint32 {
	v--
	v |= v >> 1
	v |= v >> 2
	v |= v >> 4
	v |= v >> 8
	v |= v >> 16
	v++
	return v
}
