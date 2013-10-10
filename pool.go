package pool

import (
	"time"
)

var pools = make(map[string]*pool)

type pool struct {
	max       uint
	resources chan interface{}
	create    func() interface{}
	destroy   func(interface{})
}

/*
 * Creates a new resource pool
 */
func Initialize(name string, max uint, create func() interface{}, destroy func(interface{})) {
	p := new(pool)
	p.max = max
	p.create = create
	p.destroy = destroy
	p.resources = make(chan interface{}, max)
	for i := uint(0); i < max; i++ {
		resource := create()
		p.resources <- resource
	}
	pools[name] = p
}

func (p *pool) add() {
	resource := p.create()
	p.resources <- resource
}

func Name(name string) *pool {
	return pools[name]
}

/*
 * Obtain a resource from the pool.  Wait indefinately until there is a
 * resource available.
 */
func (p *pool) Acquire() interface{} {
	return <-p.resources
}

/*
 * Obtain a resource from the pool but only wait for a specified duration.
 * If the duration expires return nil.
 */
func (p *pool) AcquireWithTimeout(timeout time.Duration) interface{} {
	var resource interface{}
	select {
	case resource = <-p.resources:
	case <-time.After(timeout):
		return nil
	}
	return resource
}

/*
 * Returns a resource back in to the pool
 */
func (p *pool) Release(resource interface{}) {
	p.resources <- resource
}

/*
 * Remove a resource from the pool.  This is helpful if the resource
 * has gone bad.  A new resource will be created in it's place.
 */
func (p *pool) Destroy(resource interface{}) {
	p.destroy(resource)
	p.resources <- p.create()
}

/*
 * Remove all resources from the pool and call the destroy method on each of
 * them.
 */
func (p *pool) Drain() {
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
