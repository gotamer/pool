package pool

import (
	"testing"
	"time"
)

func TestIntialize(t *testing.T) {
	create := func() interface{} {
		return "test"
	}
	destroy := func(interface{}) {
	}
	p := Initialize(1, create, destroy)
	msg := p.Acquire()
	if msg.(string) != "test" {
		t.Errorf("did not receive \"test\" == %s", msg.(string))
	}
}

func TestAcquireRelease(t *testing.T) {
	create := func() interface{} {
		return "test"
	}
	destroy := func(interface{}) {
	}
	p := Initialize(5, create, destroy)
	if len(p.resources) != 5 {
		t.Errorf("Pool size incorrect. Should be 5 but is %d", len(p.resources))
	}
	msg := p.Acquire()
	if len(p.resources) != 4 {
		t.Errorf("Pool size incorrect. Should be 4 but is %d", len(p.resources))
	}
	p.Release(msg)
	if len(p.resources) != 5 {
		t.Errorf("Pool size incorrect. Should be 5 but is %d", len(p.resources))
	}
}

func TestDrain(t *testing.T) {
	i := 0
	create := func() interface{} {
		return "test"
	}
	destroy := func(interface{}) {
		i++
	}
	p := Initialize(5, create, destroy)
	p.Drain()
	if i != 5 {
		t.Errorf("Drain did not call the destroy function 5 times.  Destroy was called %d times", i)
	}
}

func TestAcquireWithTimeout(t *testing.T) {
	create := func() interface{} {
		return "test"
	}
	destroy := func(interface{}) {
	}
	p := Initialize(1, create, destroy)
	p.Acquire()
	r2 := p.AcquireWithTimeout(time.Millisecond * 1)
	if r2 != nil {
		t.Errorf("A timed out acquire should return nil")
	}
}
