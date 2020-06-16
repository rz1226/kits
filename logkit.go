package kits

import (
	"bytes"
	"fmt"
	"strconv"
	"sync/atomic"
	"time"
)

const (
	sIZE      = 500  //  log存储空间总长度
	fetchSIZE = 40   //  一次拿多少
	strSIZE   = 3000 //  单条日志的最大限制
)

type LogKit struct {
	logs   *CircleQueue
	name   string //  有了这个show的时候不用锁
	readme string //  名称的注释
	used   uint32
}

func NewLogKit(name, readme string) *LogKit {
	lk := &LogKit{}
	lk.readme = readme
	lk.logs = NewCircleQueue(sIZE)
	lk.name = name
	lk.used = 0
	return lk
}
func (lk *LogKit) isUsed() bool {
	v := atomic.LoadUint32(&lk.used)
	return v != 0
}
func (lk *LogKit) setUsed() {
	atomic.StoreUint32(&lk.used, 1)
}

//  展示数据
func (lk *LogKit) Show() string {
	if !lk.isUsed() {
		return ""
	}
	str := ""
	str += "\n----------------------\n日志名称: " + lk.name + "   \n" + "日志说明:" + lk.readme + "\n"
	str += lk.FetchContents(fetchSIZE)
	return str
}

//  显示最近的count条记录
func (lk *LogKit) FetchContents(count int) string {
	values, newestID := lk.logs.GetSeveral(count)
	return formatFetchedLog(values, newestID)
}

//  把日志信息放入，然后返回格式化后的字符串
func (lk *LogKit) PutContentsAndFormat(a ...interface{}) string {
	lk.setUsed()
	cq := lk.logs
	buffer := bytes.Buffer{}
	now := time.Now()
	buffer.WriteString(now.Format("2006-01-02 15:04:05") + " 0." + fmt.Sprint(now.Nanosecond()))
	buffer.WriteString("\n")
	for _, v := range a {
		str := fmt.Sprint(v)
		if len(str) > strSIZE {
			str = str[0:strSIZE] + "......后面的内容过长截断......"
		}
		buffer.WriteString(str)
		//  每一个参数后面都加一个换行
		buffer.WriteString("\n")
	}

	logStr := buffer.String()
	cq.Put(logStr + "\n")
	return logStr
}

//  美化输出
func formatFetchedLog(values []interface{}, id uint64) string {
	buffer := bytes.Buffer{}
	buffer.WriteString("序号: " + strconv.FormatUint(id, 10) + "\n")
	for _, v := range values {
		str, ok := v.(string)
		if ok {
			buffer.WriteString(str)

		}
	}
	return buffer.String()
}
