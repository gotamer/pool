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
	Initialize("db", 1, create, destroy)
	p := Name("db")
	msg := p.Acquire()
	if msg.(string) != "test" {
		t.Errorf("did not receive \"test\" == %s", msg.(string))
	}
}

var dbuse = make(map[int]interface{})

func TestAcquireRelease(t *testing.T) {
	create := func() interface{} {
		return "test"
	}
	destroy := func(interface{}) {
	}
	Initialize("db", 5, create, destroy)
	p := Name("db")
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

	for i := 0; i < 6; i++ {
		dbuse[i] = p.Acquire()
	}
	for _, v := range dbuse {
		p.Release(v)
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
	Initialize("db", 5, create, destroy)
	p := Name("db")
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
	Initialize("db", 1, create, destroy)
	p := Name("db")
	p.Acquire()
	r2 := p.AcquireWithTimeout(time.Millisecond * 1)
	if r2 != nil {
		t.Errorf("A timed out acquire should return nil")
	}
}
