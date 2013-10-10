package pool

import (
	"log"
	"sync"
	"time"
)

var pools = make(map[string]*Pool)

type Pool struct {
	mx        sync.RWMutex
	Id        uint
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
func Initialize(name string, max uint, min uint, create func() interface{}, destroy func(interface{})) {
	p := new(Pool)
	p.max = max
	p.min = min
	p.resources = make(chan interface{}, max)
	p.create = create
	p.destroy = destroy
	for i := uint(0); i < min; i++ {
		p.New()
		pools[name] = p
	}
}

func (p *Pool) New() {
	p.mx.Lock()
	p.Count++
	p.Id = p.Count
	resource := p.create()
	p.resources <- resource
	p.mx.Unlock()
}

func Get(name string) (p *Pool) {
	return pools[name]
}

/*
 * Obtain a resource from the Pool.  Wait indefinately until there is a
 * resource available.
 */
func (p *Pool) Acquire() interface{} {
Waiting:
	p.mx.Lock()
	defer p.mx.Unlock()
	if p.Inuse < p.Count {
		p.Inuse++
	} else if p.Count < p.max {
		resource := p.create()
		p.resources <- resource
		p.Count++
		p.Id = p.Count
		p.Inuse++
	} else if p.Count >= p.max {
		var i interface{}
		return i
	} else {
		p.mx.Unlock()
		goto Waiting
	}
	if p.Count < p.min {
		go p.New()
	}
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
	log.Println("Start Release: ", resource)
	p.mx.Lock()
	defer p.mx.Unlock()
	log.Println("Release Unlocked: ", resource)
	if p.Count-p.Inuse > p.min {
		log.Println("Release Destroy: ", resource)
		p.destroy(resource)
		p.Count--
	} else {
		log.Println("Release add to resources: ", resource)
		p.resources <- resource
	}
	p.Inuse--
	log.Println("Resource Count: ", p.Count)
}

/*
 * Remove a resource from the Pool.  This is helpful if the resource
 * has gone bad.  A new resource will be created in it's place.
 */
func (p *Pool) Destroy(resource interface{}) {
	p.mx.Lock()
	defer p.mx.Unlock()
	p.destroy(resource)
	p.Count--
	p.Inuse--
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
