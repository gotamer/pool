# pool - A generic resource pool for golang

### Example Usage

```go
package main

import "pool"

func main() {
  create := func() (interface{}) {
    return "resource" // return a resource
  }
  destroy := func(resource interface{}) {
    // clean up resource
  }
  // create a pool named "myDB" with min 5 and max 10 resources
  rpool := pool.Initialize("myDB", 5, 10, create, destroy) 
  resource := rpool.Acquire() // obtain the resource
  // use resource ...
  rpool.Release(resource) // return resource to the pool
  rpool.Drain() // free up all resources
}
```

### Features

- Simple interface
- Very little code
- Lazy creation of resources (create resources when needed, but keep a min on hand)

