package unittestcases

import (
	"context"

	"github.com/starter-go/afs"
	"github.com/starter-go/buckets"
	"github.com/starter-go/units"
)

type CaseTryFileAPI struct {

	//starter:component

	DM      buckets.DriverManager //starter:inject("#")
	Service buckets.Service       //starter:inject("#")
	DirMan  units.DirManager      //starter:inject("#")
}

// ListRegistrations implements units.Unit.
func (inst *CaseTryFileAPI) ListRegistrations(list []*units.Registration) []*units.Registration {

	// old := list

	r1 := &units.Registration{
		Name:     "case-try-file-api",
		Enabled:  true,
		Priority: 0,
		Do:       inst.run,
	}

	list = append(list, r1)

	return list
}

func (inst *CaseTryFileAPI) _impl() units.Unit {
	return inst
}

func (inst *CaseTryFileAPI) run(ctx context.Context) error {

	holder := new(buckets.BucketHolder)

	holder.SetService(inst.Service)
	holder.SetName("demo1")
	holder.SetLazy(true)
	holder.SetContext(ctx)

	inst.getTmpDir(ctx)

	return nil
}

func (inst *CaseTryFileAPI) getTmpDir(ctx context.Context) (afs.Path, error) {

	h1 := &units.DirHolder{
		Context: ctx,
		Scope:   units.DirScopeModule,
		Key:     units.DirKeyConfig,
	}

	h2, err := inst.DirMan.GetDir(h1)
	if err != nil {
		return nil, err
	}

	return h2.Path, nil
}

////////////////////////////////////////////////////////////////////////////////
// EOF
