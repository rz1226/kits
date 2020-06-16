package kits

import (
	"fmt"
	"testing"

	"time"
)

func Test_timerkit(t *testing.T) {
	tk := NewTimerKit("a", "b")
	for i := 0; i < 100; i++ {

		tick := tk.Start("调用abc接口")
		tick2 := tk.Start("调用abc接口1")
		tick3 := tk.Start("调用abc接口2")
		time.Sleep(time.Millisecond * 10)
		tk.End(tick)
		time.Sleep(time.Millisecond * 10)
		tk.End(tick2)
		time.Sleep(time.Millisecond * 10)
		tk.End(tick3)

	}
	fmt.Println(tk.Info())
	fmt.Println(tk.Info())
	fmt.Println(tk.Info())

}
