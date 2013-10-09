// pool.go
package main

import (
	"fmt"
	"github.com/gotamer/pool"
)

var (
	P  *pool.Pool
	no int
)

func init() {
	create := func() interface{} {
		return test()
	}
	destroy := func(resource interface{}) {
		// clean up resource
	}
	P = pool.Initialize(5, create, destroy) // create a pool of 5 resources
}

func final() {
	P.Drain() // free up all resources
}

func main() {
	for i := 0; i < 20; i++ {
		resource := P.Acquire() // obtain the resource
		fmt.Println("Resource: ", resource)
		if i%2 == 0 {
			P.Destroy(resource)
		} else {
			P.Release(resource)
		}
	}
	final()
}

func test() int {
	no++
	return no
}
