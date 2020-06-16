package kits

import "sync/atomic"

type CounterKit struct {
	data   *Counter
	name   string
	readme string //  名称的注释
	used   uint32 //  是否已经使用，只有用过的才会show的时候显示出来 0 没有使用 1 使用
}

func NewCounterKit(name, readme string) *CounterKit {
	c := &CounterKit{}
	c.name = name
	c.data = NewCounter()
	c.readme = readme
	c.used = 0
	return c
}
func (c *CounterKit) isUsed() bool {
	v := atomic.LoadUint32(&c.used)
	return v != 0
}
func (c *CounterKit) setUsed() {
	atomic.StoreUint32(&c.used, 1)
}

func (c *CounterKit) Show() string {
	if !c.isUsed() {
		return ""
	}
	str := ""
	str += "\n----------------------\n计数器名称:" + c.name + " : \n计数器信息:" + c.readme + "\n"
	str += c.Str()
	str += "\n"

	return str
}

func (c *CounterKit) Inc() {
	c.setUsed()
	c.data.Add(1)
}
func (c *CounterKit) IncBy(num int64) {
	c.setUsed()
	c.data.Add(num)
}
func (c *CounterKit) Get(name string) int64 {
	return c.data.Get()
}

func (c *CounterKit) Str() string {
	return c.data.Str()
}
