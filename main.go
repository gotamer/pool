package main

import "time"

import "./gopool"

func main() {
	create := func() (interface{}) {
		time.Sleep(time.Millisecond * 100)
		return "hello world"
	}
	destroy := func(gopool.Resource) {
		println("bye bye world")
	}

	myPool := gopool.Initialize(10, create, destroy)

	err, resource1 := myPool.Acquire()
	if (err == nil) {
		msg := resource1.Payload
		println(msg.(string))
	} else {
		println(err)
	}
	go func() { time.Sleep(time.Millisecond * 1000); myPool.Release(resource1) }()

	err, resource2 := myPool.Acquire()
	if (err == nil) {
		msg := resource2.Payload
		println(msg.(string))
	} else {
		println(err)
	}
	go func() { time.Sleep(time.Millisecond * 1000); myPool.Release(resource2) }()

	err, resource3 := myPool.Acquire()
	if (err == nil) {
		msg := resource3.Payload
		println(msg.(string))
	} else {
		println(err)
	}

	err, resource4 := myPool.Acquire()
	if (err == nil) {
		msg := resource4.Payload
		println(msg.(string))
	} else {
		println(err)
	}
	myPool.Release(resource3)
	myPool.Release(resource4)

	time.Sleep(time.Millisecond * 1000)

	myPool.Drain()

	time.Sleep(time.Millisecond * 1000)

}




