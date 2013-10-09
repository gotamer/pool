package pool

import (
	"time"
)

type Pool struct {
	max       int
	resources chan interface{}
	create    func() interface{}
	destroy   func(interface{})
}

/*
 * Creates a new resource Pool
 */
func Initialize(max int, create func() interface{}, destroy func(interface{})) *Pool {
	p := new(Pool)
	p.max = max
	p.resources = make(chan interface{}, max)
	for i := 0; i < max; i++ {
		resource := create()
		p.resources <- resource
	}
	p.create = create
	p.destroy = destroy

	return p
}

/*
 * Obtain a resource from the Pool.  Wait indefinately until there is a
 * resource available.
 */
func (p *Pool) Acquire() interface{} {
	return <-p.resources
}

/*
 * Obtain a resource from the Pool but only wait for a specified duration.
 * If the duration expires return nil.
 */
func (p *Pool) AcquireWithTimeout(timeout time.Duration) interface{} {
	var resource interface{}
	select {
	case resource = <-p.resources:
	case <-time.After(timeout):
		return nil
	}
	return resource
}

/*
 * Returns a resource back in to the Pool
 */
func (p *Pool) Release(resource interface{}) {
	p.resources <- resource
}

/*
 * Remove a resource from the Pool.  This is helpful if the resource
 * has gone bad.  A new resource will be created in it's place.
 */
func (p *Pool) Destroy(resource interface{}) {
	p.destroy(resource)
	p.resources <- p.create()
}

/*
 * Remove all resources from the Pool and call the destroy method on each of
 * them.
 */
func (p *Pool) Drain() {
	for {
		select {
		case r := <-p.resources:
			p.destroy(r)
		default:
			return
		}
	}
	close(p.resources)
}
