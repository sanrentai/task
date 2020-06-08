package main

import (
	"errors"
	"fmt"
	"runtime"
	"time"

	"github.com/sanrentai/task"
)

func main() {
	task.InitTaskReceiver(runtime.NumCPU())
	fmt.Println(runtime.NumCPU())
	mytask1 := task.NewTask(map[string]interface{}{
		"a": "1234",
		"b": 123,
		"c": "mytask1",
	}, func(p map[string]interface{}) (interface{}, error) {
		fmt.Println("mytask1:", p)

		time.Sleep(1 * time.Second)
		fmt.Println(p)
		return p["b"], nil
	}, 3*time.Second)

	defer mytask1.Close()
	// id := AddTask(mytask1)
	mytask1.Start()
	// time.Sleep(3 * time.Second)
	mytask2 := task.NewTask(map[string]interface{}{
		"a": "1234",
		"b": 123,
		"c": "mytask2",
	}, func(p map[string]interface{}) (interface{}, error) {
		fmt.Println("mytask2:", p)

		time.Sleep(2 * time.Second)
		fmt.Println(p)
		return p["b"], nil
	}, 3*time.Second)

	defer mytask2.Close()

	mytask2.Start()

	mytask3 := task.NewTask(map[string]interface{}{
		"a": "1234",
		"b": 123,
		"c": "mytask3",
	}, func(p map[string]interface{}) (interface{}, error) {
		fmt.Println("mytask3:", p)

		time.Sleep(1 * time.Second)
		fmt.Println(p)
		return p["b"], errors.New("错误测试")
	}, 3*time.Second)

	defer mytask3.Close()
	mytask3.Start()

	fmt.Println("mytask1:", mytask1.ID)
	res1, err := mytask1.GetResult()
	if err != nil {
		fmt.Println("res1 err:", err.Error())
	} else {
		fmt.Println("res1:", res1)
	}
	fmt.Println("mytask2:", mytask2.ID)
	res2, err := mytask2.GetResult()
	if err != nil {
		fmt.Println("res2 err:", err.Error())
	} else {
		fmt.Println("res2:", res2)
	}
	fmt.Println("mytask3:", mytask3.ID)
	res3, err := mytask3.GetResult()
	if err != nil {
		fmt.Println("res3 err:", err.Error())
	} else {
		fmt.Println("res3:", res3)
	}
}
