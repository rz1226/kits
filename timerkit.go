package kits

//  时间记录
import (
	"fmt"
	"strconv"
	"sync/atomic"
	"time"
)

/*
		tick := tk.Start("a","调用abc接口")
		tick2 := tk.Start("b","调用abc接口1")
		tick3 := tk.Start("c","调用abc接口2")
		time.Sleep(time.Millisecond*10)
		tk.End( tick)
		time.Sleep(time.Millisecond*10)
		tk.End(tick2)
		time.Sleep(time.Millisecond*10)
		tk.End(tick3)
// 记录最近的时间比较长的操作
*/
const (
	sLOWLOGNAME   = "slowest"
	lASTLOGNAME   = "last"
	sLOWNANO      = 10000000 //  耗时大于这个+平均数算作慢   (nano单位)
	TIMEKITMAXKIT = 1000
)

type TimerKit struct {
	node   *timerKitNode
	name   string
	readme string //  名称的注释
	used   uint32
}

func NewTimerKit(name, readme string) *TimerKit {

	tk := &TimerKit{}
	tk.readme = readme
	tk.node = newTimerKitNode(name)
	tk.name = name
	tk.used = 0
	return tk
}
func (tk *TimerKit) isUsed() bool {
	v := atomic.LoadUint32(&tk.used)
	return v != 0
}
func (tk *TimerKit) setUsed() {
	atomic.StoreUint32(&tk.used, 1)
}
func (tk *TimerKit) Show() string {
	if !tk.isUsed() {
		return ""
	}
	str := ""
	str += "----------------------\n计时器名称:" + tk.name + "\n计时器说明:" + tk.readme + "\n"
	str += tk.Info()
	return str
}

func (tk *TimerKit) Start(tickInfo string) *Tick {

	return tk.node.start(tickInfo)
}
func (tk *TimerKit) End(tick *Tick) {
	tk.setUsed()
	if tick == nil {
		return
	}
	if tick.timerName != tk.name {
		return
	}

	tk.node.end(tick)
}
func (tk *TimerKit) Info() string {
	timerKitNode := tk.node
	//  这里不用考虑组装字符串的性能，因为没必要， 而且数据很小的常见，+的低性能劣势并不明显
	resStr := ""
	resStr += "count:" + strconv.FormatInt(timerKitNode.getCount(), 10) + "次\n"
	resStr += "sum:" + strconv.FormatFloat(timerKitNode.getSum(), 'f', 6, 64) + "秒\n"
	resStr += "avg:" + strconv.FormatFloat(timerKitNode.getAvg(), 'f', 6, 64) + "秒每次\n"
	resStr += "-----\n高耗时记录: \n" + timerKitNode.showslow()
	resStr += "-----\n最近记录: \n" + timerKitNode.showlast()
	return resStr
}

type timerKitNode struct {
	name    string  //  计时器的名字
	count   int64   //  经过了多少次计数
	sum     int64   //  总计时时间
	slowest *LogKit //  该计时器的慢操作列表
	last    *LogKit
}

// 单次计时操作
type Tick struct {
	timerName string
	info      string
	startTime int64
}

func newTimerKitNode(name string) *timerKitNode {
	tkn := &timerKitNode{}
	tkn.name = name
	tkn.count = 0
	tkn.sum = 0
	tkn.slowest = NewLogKit(sLOWLOGNAME, "最慢记录")
	tkn.last = NewLogKit(lASTLOGNAME, "最近记录")
	return tkn
}

func (t *timerKitNode) start(tickInfo string) *Tick {
	te := &Tick{}
	te.info = tickInfo
	te.timerName = t.name
	te.startTime = time.Now().UnixNano()
	return te
}

func (t *timerKitNode) end(tick *Tick) {
	du := time.Now().UnixNano() - tick.startTime
	atomic.AddInt64(&t.count, 1)
	atomic.AddInt64(&t.sum, du)
	count := atomic.LoadInt64(&t.count)
	t.last.PutContentsAndFormat("操作是:"+tick.info, "耗时秒是:"+fmt.Sprint(float64(du)/float64(time.Second)))
	sum := atomic.LoadInt64(&t.sum)
	if du > sLOWNANO+sum/(count+1) {
		t.slowest.PutContentsAndFormat("操作是:"+tick.info, "耗时秒是:"+fmt.Sprint(float64(du)/float64(time.Second)))
	}
}
func (t *timerKitNode) getCount() int64 {
	return atomic.LoadInt64(&t.count)
}

func (t *timerKitNode) getSum() float64 {
	return float64(atomic.LoadInt64(&t.sum)) / float64(time.Second)
}
func (t *timerKitNode) getAvg() float64 {
	sum := atomic.LoadInt64(&t.sum)
	count := atomic.LoadInt64(&t.count)
	if count == 0 {
		return 0
	}
	avg := sum / count
	return float64(avg) / float64(time.Second)
}

const FETCHLEN = 10

func (t *timerKitNode) showslow() string {
	return t.slowest.FetchContents(FETCHLEN)
}
func (t *timerKitNode) showlast() string {
	return t.last.FetchContents(FETCHLEN)
}
