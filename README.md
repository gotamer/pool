# Gopool - A generic resource pool for golang

### Example Usage

```go
package main

import "gopool"

func main() {
  create := func() (interface{}) {
    return "resource" // return a resource
  }
  destroy := func(resource interface{}) {
    // clean up resource
  }
  pool := gopool.Initialize(5, create, destroy) // create a pool of 5 resources
  resource := pool.Acquire() // obtain the resource
  // use resource ...
  pool.Release(resource) // return resource to the pool
  pool.Drain() // free up all resources
}
```

### Features

- Simple interface
- Very little code

### Possible Future Features

- Lazy creation of resources (only create resources when needed)
- Ability to handle inability to create resources by periodically trying to create resources in a goroutine (example: allow an app to start up without a database connection then periodically attempt to create the database connections)
