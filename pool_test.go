package pool

import (
	"testing"
	//"time"
)

func TestIntialize(t *testing.T) {
	create := func() interface{} {
		return "test"
	}
	destroy := func(interface{}) {
	}
	Initialize("db", 3, 1, create, destroy)
	p := Get("db")
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
	Initialize("db", 50, 10, create, destroy)
	p := Get("db")
	if len(p.resources) != 10 {
		t.Errorf("Pool size incorrect. Should be 10 but is %d", len(p.resources))
	}
	msg := p.Acquire()
	if len(p.resources) != 9 {
		t.Errorf("Pool size incorrect. Should be 9 but is %d", len(p.resources))
	}
	p.Release(msg)
	if len(p.resources) != 10 {
		t.Errorf("Pool size incorrect. Should be 10 but is %d", len(p.resources))
	}
	var iface interface{}
	for i := 0; i < 55; i++ {
		if r := p.Acquire(); r != iface {
			dbuse[i] = r
		}
	}
	for _, v := range dbuse {
		p.Release(v)
	}

	if len(p.resources) < 10 {
		t.Errorf("Pool size incorrect. Should be 10 or more but is %d", len(p.resources))
	}
	if len(p.resources) > 50 {
		t.Errorf("Pool size incorrect. Should be 50 or less but is %d", len(p.resources))
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
	Initialize("db", 50, 5, create, destroy)
	p := Get("db")
	p.Drain()
	if i != 5 {
		t.Errorf("Drain did not call the destroy function 5 times.  Destroy was called %d times", i)
	}
}

/*
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
*/
