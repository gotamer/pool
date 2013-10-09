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
  rpool := pool.Initialize(5, create, destroy) // create a pool of 5 resources
  resource := rpool.Acquire() // obtain the resource
  // use resource ...
  rpool.Release(resource) // return resource to the pool
  rpool.Drain() // free up all resources
}
```

### Features

- Simple interface
- Very little code

### Possible Future Features

- Lazy creation of resources (only create resources when needed)
- Ability to handle inability to create resources by periodically trying to create resources in a goroutine (example: allow an app to start up without a database connection then periodically attempt to create the database connections)
