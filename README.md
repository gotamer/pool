# pool - A generic resource pool for go (golang)

### Example Usage

```go
package main

import "pool"

func init(){
	resourceOpen := func() (interface{}, error) {
		db, err := sql.Open("sqlite3", "/tmp/fool.db")
		return db, err
	}
	resourceClose := func(r interface{}) error {
		db := resource.(*sql.DB)
		return db.Close() 
	}
	err = Register("db1", 10, 20, resourceOpen, resourceClose)
}

func main() {
	rp := pool.Name("db1")
	resource, err := rp.Get() // obtain a resource
	// use resource ... db what ever
	rp.Release(resource) // return resource to the pool
	...
	rp.Close("db1") // free up all resources and delete the pool
}
```

### Features

- Simple interface
- Lazy creation of resources (create resources when needed, but keep a min on hand)

### Notes
- Not backwards compadable with version 1,2 and 3
- db.resourceDel() "as in db.Close()" must return errors just as go sql databases do

### ToDo
- Find a better way to handle errors when closing resources

### Credits
- This was originally started by https://github.com/braindev/gopool, I fixed some bugs and added more futures

________________________________________________________

The MIT License (MIT)
=====================
Copyright © 2012 - 2013 GoTamer <http://www.robotamer.com>

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 