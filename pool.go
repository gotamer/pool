// A generic resource pool
package pool

import (
	"sync"
	"time"
)

var pools = make(map[string]*pool)

type pool struct {
	mx        sync.RWMutex
	count     uint
	inuse     uint
	max       uint
	min       uint // Minimum Available
	resources chan interface{}
	create    func() interface{}
	destroy   func(interface{})
}

/*
 * Creates a new resource Pool
 */
func Initialize(name string, min uint, max uint, create func() interface{}, destroy func(interface{})) {
	p := new(pool)
	p.max = max
	p.min = min
	p.resources = make(chan interface{}, max)
	p.create = create
	p.destroy = destroy
	for i := uint(0); i < min; i++ {
		p.count++
		resource := p.create()
		p.resources <- resource
	}
	pools[name] = p
}

func (p *pool) New() {
	p.mx.Lock()
	p.count++
	resource := p.create()
	p.resources <- resource
	p.mx.Unlock()
}

func Name(name string) (p *pool) {
	return pools[name]
}

/*
 * Obtain a resource from the Pool.
 * Returns nil if there are no more resources available (Set pool.max)
 */
func (p *pool) Acquire() interface{} {
Waiting:
	p.mx.Lock()
	defer p.mx.Unlock()
	if p.inuse < p.count {
		p.inuse++
	} else if p.count < p.max {
		resource := p.create()
		p.resources <- resource
		p.count++
		p.inuse++
	} else if p.count >= p.max {
		return nil // No resources availabe
	} else {
		p.mx.Unlock()
		goto Waiting
	}
	if p.count < p.min {
		go p.New()
	}
	return <-p.resources
}

/*
 * Obtain a resource from the Pool but only wait for a specified duration.
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
 * Returns a resource back in to the Pool
 */
func (p *pool) Release(resource interface{}) {
	p.mx.Lock()
	defer p.mx.Unlock()
	if p.count-p.inuse > p.min {
		p.destroy(resource)
		p.count--
	} else {
		p.resources <- resource
	}
	p.inuse--
}

/*
 * Remove a resource from the Pool.  This is helpful if the resource
 * has gone bad.  A new resource will be created in it's place.
 */
func (p *pool) Destroy(resource interface{}) {
	p.mx.Lock()
	defer p.mx.Unlock()
	p.destroy(resource)
	p.count--
	p.inuse--
}

/*
 * Remove all resources from the Pool and call the destroy method on each of
 * them.
 */
func (p *pool) Drain() {
	for {
		select {
		case r := <-p.resources:
			p.inuse--
			p.count--
			p.destroy(r)
		default:
			return
		}
	}
	close(p.resources)
}
