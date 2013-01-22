package main

import (
	 "./gopool"
	 "time"
	 "fmt"
)

func main() {
	create := func() (interface{}) {
		t := time.Now()
		return t.String()
	}
	destroy := func(r interface{}) {
		fmt.Printf("destroyed: %s\n", r.(string))
	}
	pool := gopool.Initialize(5, create, destroy)

	r1 := pool.Acquire()
	fmt.Println(r1.(string))
	go func() {
		time.Sleep(1000 * time.Millisecond)
		pool.Release(r1)
	}()
	r2 := pool.Acquire()
	fmt.Println(r2.(string))
	
	r3 := pool.Acquire()
	fmt.Println(r3.(string))
	
	r4 := pool.Acquire()
	fmt.Println(r4.(string))
	
	r5 := pool.Acquire()
	fmt.Println(r5.(string))
	
	r6 := pool.Acquire()
	fmt.Println(r6.(string))
	
	//pool.Destroy(r6)

	pool.Release(r6)
	pool.Release(r5)
	pool.Release(r4)
	pool.Release(r3)
	pool.Release(r2)

	pool.Drain()


}




