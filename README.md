# pool - A generic resource pool for go (golang)

### Example Usage

```go
package main

import "pool"

func init(){
	create := func() (interface{}) {
	  return "resource" // return a resource
	}
	destroy := func(resource interface{}) {
	  // clean up resource
	}
	// create a pool named "myDB" with min 5 and max 10 resources
	pool.Initialize("myDB", 5, 10, create, destroy) 	
}

func main() {
	rp := pool.Name("myDB")
	resource := rp.Resource() // obtain the resource
	// use resource ... db what ever
	rp.Release(resource) // return resource to the pool
	rp.Drain() // free up all resources
}
```

### Features

- Simple interface
- Very little code
- Lazy creation of resources (create resources when needed, but keep a min on hand)

