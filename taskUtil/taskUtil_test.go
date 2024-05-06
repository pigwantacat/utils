package taskutil

import (
	"testing"
	"time"
)

func TestTimingWheel(t *testing.T) {
	tw, _ := NewTimingWheel(time.Second, 12)
	start := time.Now()
	t.Log("start ---> ", start)
	err := tw.AddTask("task_test", func(key string) {
		t.Log(key, "doTask...")
		t.Log("now ---> ", time.Since(start))
	}, 2*time.Second, 3)
	if err != nil {
		return
	}
	t.Log("now ---> ", time.Now())
	time.Sleep(5 * time.Second)
	err = tw.RemoveTask("task_test")
	if err != nil {
		return
	}
	t.Log("now ---> ", time.Now())
}
