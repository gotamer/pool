package pool

import (
	"sync"
	"time"
)

type Pool struct {
	mx        sync.RWMutex
	Count     uint
	Inuse     uint
	max       uint
	min       uint // Minimum Available
	resources chan interface{}
	create    func() interface{}
	destroy   func(interface{})
}

/*
 * Creates a new resource Pool
 */
func Initialize(max uint, min uint, create func() interface{}, destroy func(interface{})) *Pool {
	p := new(Pool)
	p.max = max
	p.min = min
	p.resources = make(chan interface{}, max)
	p.create = create
	p.destroy = destroy
	for i := uint(0); i < min; i++ {
		p.New()
	}
	return p
}

func (p *Pool) New() {
	p.mx.RLock()
	if p.Count < p.max {
		p.mx.RUnlock()
		p.mx.Lock()
		p.Count++
		resource := p.create()
		p.resources <- resource
		p.mx.Unlock()
	} else {
		p.mx.RUnlock()
	}
}

/*
 * Obtain a resource from the Pool.  Wait indefinately until there is a
 * resource available.
 */
func (p *Pool) Acquire() interface{} {
	p.mx.RLock()
	if p.Inuse < p.Count {
		p.mx.RUnlock()
		p.mx.Lock()
		p.Inuse++
		p.mx.Unlock()
		return <-p.resources
	} else {
		p.mx.RUnlock()
		p.New()
		return <-p.resources
	}
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
	p.mx.Lock()
	p.Inuse--
	if p.Count-p.Inuse >= p.min {
		p.Inuse++
		p.mx.Unlock()
		p.Destroy(resource)
	} else {
		p.resources <- resource
		p.mx.Unlock()
	}
}

/*
 * Remove a resource from the Pool.  This is helpful if the resource
 * has gone bad.  A new resource will be created in it's place.
 */
func (p *Pool) Destroy(resource interface{}) {
	p.mx.Lock()
	defer p.mx.Unlock()
	p.Count--
	p.Inuse--
	p.destroy(resource)
}

/*
 * Remove all resources from the Pool and call the destroy method on each of
 * them.
 */
func (p *Pool) Drain() {
	for {
		select {
		case r := <-p.resources:
			p.Inuse--
			p.Count--
			p.destroy(r)
		default:
			return
		}
	}
	close(p.resources)
}
