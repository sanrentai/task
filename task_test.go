package task

import (
	"errors"
	"runtime"
	"testing"
	"time"
)

func TestTask(t *testing.T) {
	InitTaskReceiver(runtime.NumCPU())
	t.Log(runtime.NumCPU())
	mytask1 := NewTask(map[string]interface{}{
		"a": "1234",
		"b": 123,
		"c": "mytask1",
	}, func(p map[string]interface{}) (interface{}, error) {
		t.Log("mytask1:", p)

		time.Sleep(1 * time.Second)
		t.Log(p)
		return p["b"], nil
	}, 3*time.Second)

	defer mytask1.Close()
	// id := AddTask(mytask1)
	mytask1.Start()
	// time.Sleep(3 * time.Second)
	mytask2 := NewTask(map[string]interface{}{
		"a": "1234",
		"b": 123,
		"c": "mytask2",
	}, func(p map[string]interface{}) (interface{}, error) {
		t.Log("mytask2:", p)

		time.Sleep(2 * time.Second)
		t.Log(p)
		return p["b"], nil
	}, 3*time.Second)

	defer mytask2.Close()

	mytask2.Start()

	mytask3 := NewTask(map[string]interface{}{
		"a": "1234",
		"b": 123,
		"c": "mytask3",
	}, func(p map[string]interface{}) (interface{}, error) {
		t.Log("mytask3:", p)

		time.Sleep(1 * time.Second)
		t.Log(p)
		return p["b"], errors.New("错误测试")
	}, 3*time.Second)

	defer mytask3.Close()
	mytask3.Start()

	t.Log("mytask1:", mytask1.ID)
	res1, err := mytask1.GetResult()
	if err != nil {
		t.Log(err)
	} else {
		t.Log(res1)
	}
	t.Log("mytask2:", mytask2.ID)
	res2, err := mytask2.GetResult()
	if err != nil {
		t.Log(err)
	} else {
		t.Log(res2)
	}
	t.Log("mytask3:", mytask3.ID)
	res3, err := mytask3.GetResult()
	if err != nil {
		t.Log(err)
	} else {
		t.Log(res3)
	}

}
