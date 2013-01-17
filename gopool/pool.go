package gopool

import "errors"

type Resource struct {
	Payload interface{}
//	UseTime int
}

type pool struct {
	max int
//	min int // TODO support min
	freeResources chan Resource
	inUseResources chan Resource
//	maxIdleTime int // TODO implement?
	create func() interface{}
	destroy func()
}

func Initialize(max int, create func() interface{}, destroy func()) (pool) {
	p := pool{}
	p.create = create
	p.max = max
	p.destroy = destroy
	p.freeResources = make(chan Resource)
	p.inUseResources = make(chan Resource)
	return p
}

func (p *pool) Drain() {

}

func (p *pool) Acquire() (error, Resource) {
	select {
	case resource := <-p.freeResources:
		// resource available
		p.inUseResources <- resource
		return nil, resource
	default:
		// resource not immediately available
		if len(p.freeResources) + len(p.inUseResources) < p.max - 1 {
			// if we have less resources than max 
			// attempt to create a new resource
			resource := Resource{}
			val := p.create()
			if (val != nil) {
				p.inUseResources <- val.(Resource)
				resource.Payload = val
				return nil, resource
			} else {
				return errors.New("Could not create new resource"), Resource{}
			}
		} else {
			// wait for a resource to become available
			resource := <-p.freeResources
			p.inUseResources <- resource
			return nil, resource
		}
	}
	return errors.New("Could not acquire or create resource"), Resource{}
}

func (p *pool) Release(resource Resource) {
}
