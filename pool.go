// A generic resource pool for databases etc
package pool

import (
	"fmt"
	"log"
	"sync"
	"time"
)

var pools = make(map[string]*ResourcePool)

// ResourcePool allows you to use a pool of resources.
type ResourcePool struct {
	mx        sync.RWMutex
	min       uint // Minimum Available resources
	inUse     uint
	resources chan resourceWrapper
	resOpen   func() (interface{}, error)
	resClose  func(interface{}) error
}

type resourceWrapper struct {
	resource interface{}
	timeUsed time.Time
}

/*
 * Creates a new resource Pool
 */
func NewResourcePool(name string, min uint, max uint, o func() (interface{}, error), c func(interface{}) error) (err error) {
	p := new(ResourcePool)
	p.min = min
	p.resources = make(chan resourceWrapper, max)
	p.resOpen = o
	p.resClose = c

	for i := uint(0); i < min; i++ {
		resource, err := p.resOpen()
		if err == nil {
			var wrapper resourceWrapper
			wrapper.resource = resource
			p.resources <- wrapper
		} else {
			break
		}
	}
	for i := uint(0); i < max-min; i++ {
		p.resources <- resourceWrapper{}
	}
	if err == nil {
		pools[name] = p
	}
	return
}

func (p *ResourcePool) add() (err error) {
	p.mx.Lock()
	defer p.mx.Unlock()
	// make sure we are not going over limit
	if p.Cap() > p.Count() {
		resource, err := p.resOpen()
		var ok bool
		if err == nil {
			var wrapper resourceWrapper
			wrapper.resource = resource
			if ok {
				p.resources <- wrapper
			}
		}
	}
	return
}

// Will return a pool
func Name(name string) *ResourcePool {
	return pools[name]
}

// Get will return the next available resource. If capacity
// has not been reached, it will create a new one using the factory. Otherwise,
// it will indefinitely wait untill the next resource becomes available.
func (p *ResourcePool) Get() (resource resourceWrapper, err error) {
	return p.get()
}

// Fetch a new resource
func (p *ResourcePool) get() (resource resourceWrapper, err error) {
	p.mx.Lock()
	defer p.mx.Unlock()

	var wrapper resourceWrapper
	var ok bool
	for {
		if p.AvailableNow() != 0 {
			select {
			case wrapper, ok = <-p.resources:
			default:
				wrapper, ok = <-p.resources
			}
			if !ok {
				return wrapper, fmt.Errorf("ResourcePool is closed")
			}
			break
		} else if p.AvailableMax() != 0 {
			wrapper.resource, err = p.resOpen()
			if err == nil {
				p.resources <- wrapper
			}
			break
		}
	}

	if p.AvailableNow() < p.min && p.Cap() > p.Count() {
		go p.add()
	}
	p.inUse++
	return wrapper, err
}

/*
 * Returns a resource back in to the Pool
 */
func (p *ResourcePool) Release(wrapper resourceWrapper) {
	p.mx.Lock()
	defer p.mx.Unlock()
	if p.AvailableNow() > p.min {
		if err := p.resClose(wrapper.resource); err != nil {
			log.Println("Resource close error: ", err)
		} else {
			p.inUse--
		}
	} else {
		p.resources <- wrapper
		p.inUse--
	}
}

/*
 * Remove a resource from the Pool.  This is helpful if the resource
 * has gone bad.  A new resource will be created in it's place.
 */
func (p *ResourcePool) Destroy(wrapper resourceWrapper) {
	p.mx.Lock()
	defer p.mx.Unlock()
	if err := p.resClose(wrapper.resource); err != nil {
		log.Println("Resource close error: ", err)
	} else {
		p.inUse--
	}
}

// Remove all resources from the Pool.
// Then close the pool.
// Then delete the pool from pools
func (p *ResourcePool) Close(name string) {
	for {
		select {
		case resource := <-p.resources:
			p.resClose(resource)
		default:
			close(p.resources)
			delete(pools, name)
			return
		}
	}
}

// Resources already obtained and available for use
func (p *ResourcePool) AvailableNow() uint {
	return uint(len(p.resources))
}

// Total # of resoureces including the once we haven't yet created - whats in use
func (p *ResourcePool) AvailableMax() uint {
	return p.Cap() - p.inUse
}

// Count of resources open (should be less theen Cap())
func (p *ResourcePool) Count() uint {
	return p.inUse + p.AvailableNow()
}

// Resources being used right now
func (p *ResourcePool) InUse() uint {
	return p.inUse
}

// Max resources the pool allows; all in use, obtained, and not obtained.
func (p *ResourcePool) Cap() uint {
	return uint(cap(p.resources))
}

// Reterns how many resources we need to add to the pool, to get the reserve to reach min
func (p *ResourcePool) Short() (need uint) {
	an := p.AvailableNow()
	if an < p.min {
		need = p.min - an
	}
	return
}
