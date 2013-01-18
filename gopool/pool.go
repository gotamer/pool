package gopool

import "errors"

type Resource struct {
	Payload interface{}
//	UseTime int
}

type pool struct {
	count int
	max int
//	min int // TODO support min?
	resources chan Resource
//	maxIdleTime int // TODO support an idle time?
	create func() (interface{})
	destroy func(Resource)
}

func Initialize(max int, create func() interface{}, destroy func(Resource)) (pool) {
	p := pool{}
	p.count = 0
	p.create = create
	p.max = max
	p.destroy = destroy
	p.resources = make(chan Resource, 10)
	return p
}

func (p *pool) Drain() {
	println("draining pool")
	println(len(p.resources))
	for len(p.resources) > 0 {
		resource := <-p.resources
		p.destroy(resource)
	}
}

func (p *pool) Acquire() (error, Resource) {
	select {
	case resource := <-p.resources:
		// resource available
		println("Obtained free resource")
		return nil, resource
	default:
		// resource not immediately available
		if p.count < p.max {
			// if we have less resources than max attempt to create a new resource
			resource := Resource{}
			val := p.create()
			resource.Payload = val
			if (val != nil) {
				p.count++
				resource.Payload = val
				println("Resource created and returned")
				return nil, resource
			} else {
				println("Could not create new resource")
				return errors.New("Could not create new resource"), Resource{}
			}
		} else {
			// wait for a resource to become available
			resource := <-p.resources
			println("Previously created resource returned")
			return nil, resource
		}
	}
	return errors.New("Could not acquire or create resource"), Resource{}
}

func (p *pool) Release(resource Resource) {
	go func(){ p.resources <- resource; println("put back"); println(len(p.resources)) }()
	println("Resource returned to pool")
}
