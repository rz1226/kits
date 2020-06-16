package kits

import (
	"fmt"
	"testing"
)

func Test_counterkit(t *testing.T) {
	c := NewCounterKit("a", "b")
	c.Inc()

	fmt.Println(c.Show())

}
