// pool.go
package main

import (
	"fmt"
	"github.com/gotamer/pool"
	"sync"
	"time"
)

var (
	no int
	wg sync.WaitGroup
)

type resource_symulator struct {
	id int
}

func init() {
	create := func() interface{} {
		return resourceNew()
	}
	destroy := func(resource interface{}) {
		resourceDel(resource.(resource_symulator))
	}
	pool.Initialize("db", 3, 5, create, destroy) // create a pool of 5 resources
}

func final() {
	P := pool.Name("db")
	P.Drain() // free up all resources
}

func main() {
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			work(i)
		}(i)
	}
	wg.Wait()
	final()
}

func work(i int) {
	fmt.Println("Start Work id: ", i)
	P := pool.Name("db")
	resource := P.Acquire() // obtain the resource
	if resource == nil {
		fmt.Println("Not enough resources to run: ", i)
	} else {
		if i%2 == 0 {
			time.Sleep(time.Microsecond * 1)
		} else {
			time.Sleep(time.Microsecond * 2)
		}
		fmt.Println("Work id: ", i, " Resource id: ", resource)
		P.Release(resource)
		fmt.Println("End Work id: ", i)
	}
}

func resourceNew() (r resource_symulator) {
	no++
	r.id = no
	time.Sleep(time.Microsecond * 1)
	fmt.Println("Resource New: ", r.id)
	return
}

func resourceDel(r resource_symulator) {
	fmt.Println("Resource Del: ", r.id)
}
