package pool

import (
	//"fmt"
	"testing"
	"time"
)

var (
	no int
)

type resource_symulator struct {
	id int
}

func resourceNew() (r *resource_symulator, err error) {
	no++
	r = new(resource_symulator)
	r.id = no
	time.Sleep(time.Microsecond * 1)
	return
}

func (r *resource_symulator) resourceDel() (err error) {
	r.id = 0
	time.Sleep(time.Microsecond * 1)
	return
}

func TestIntialize(t *testing.T) {
	var db *resource_symulator
	var err error
	create := func() (interface{}, error) {
		db, err = resourceNew()
		return db, err
	}
	destroy := func(r interface{}) error {
		return db.resourceDel()
	}
	err = NewResourcePool("db", 2, 5, create, destroy)
	if err != nil {
		t.Errorf("Resource error: %s", err.Error())
	}
	p := Name("db")
	msg, err := p.Get()
	if err != nil {
		t.Errorf("Get Resource error: %s", err.Error())
	}
	if msg.resource.(*resource_symulator).id != 1 {
		t.Errorf("Resource id should be on = %d", msg)
	}
}

func TestResourceRelease(t *testing.T) {
	var db *resource_symulator
	var err error
	create := func() (interface{}, error) {
		db, err = resourceNew()
		return db, err
	}
	destroy := func(r interface{}) error {
		return db.resourceDel()
	}
	var min, max uint
	min = 10
	max = 50
	err = NewResourcePool("db", min, max, create, destroy)
	p := Name("db")
	if p.Cap() != max {
		t.Errorf("Pool size incorrect. Should be %d but is %d", max, len(p.resources))
	}
	msg, err := p.get()
	if err != nil {
		t.Errorf("get error %d", err)
	}
	if p.AvailableNow() != max-1 {
		t.Errorf("Pool size incorrect. Should be %d but is %d", max-1, len(p.resources))
	}
	p.Release(msg)
	if max-1 != p.AvailableNow() {
		t.Errorf("Pool size incorrect. Should be %d but is %d", max-1, p.AvailableNow())
	}

	var dbuse = make(map[uint]resourceWrapper)
	for i := uint(0); i < max; i++ {
		dbuse[i], err = p.Get()
		if err != nil {
			t.Errorf("get error %d", err)
		}
	}
	for _, v := range dbuse {
		p.Destroy(v)
	}

	if p.Cap() != max {
		t.Errorf("Pool cap incorrect. Should be %d but is %d", max, p.Cap())
	}

	// pools test
	po := uint(60)
	for i := uint(0); i < po; i++ {
		dbuse[i], err = p.Get()
		if err != nil {
			t.Errorf("get error %d", err)
		}
	}
	if p.InUse() != po {
		t.Errorf("Pool InUse() before release incorrect. Should be 0 but is %d", p.InUse())
	}
	if p.AvailableMax() != p.Cap()-po {
		t.Errorf("Pool AvailableMax() incorrect. Should be  %d but is %d", max-po, p.AvailableMax())
	}
	for i := uint(0); i < po; i++ {
		p.Release(dbuse[i])
	}
	if p.InUse() != 0 {
		t.Errorf("Pool InUse() incorrect. Should be 0 but is %d", p.InUse())
	}
	if p.Cap() != max {
		t.Errorf("Pool Cap() incorrect. Should be %d but is %d", max, p.Cap())
	}
	if p.AvailableNow() < min || p.AvailableNow() > max {
		t.Errorf("Pool AvailableNow() incorrect. Should be min %d, max %d but is %d", min, max, p.AvailableNow())
	}
	if p.AvailableMax() != p.Cap() {
		t.Errorf("Pool AvailableMax() incorrect. Should be  %d but is %d", max, p.AvailableMax())
	}
}

/*
func TestDrain(t *testing.T) {
	i := 0
	create := func() interface{} {
		return "test"
	}
	destroy := func(interface{}) {
		i++
	}
	Initialize("db", 5, 50, create, destroy)
	p := Name("db")
	p.Drain()
	if i != 5 {
		t.Errorf("Drain did not call the destroy function 5 times.  Destroy was called %d times", i)
	}
}
*/
/*
func TestAcquireWithTimeout(t *testing.T) {
	create := func() interface{} {
		return "test"
	}
	destroy := func(interface{}) {
	}
	Initialize("db", 1, create, destroy)
	p := Name("db")
	p.Resource()
	r2 := p.AcquireWithTimeout(time.Millisecond * 1)
	if r2 != nil {
		t.Errorf("A timed out acquire should return nil")
	}
}
*/
