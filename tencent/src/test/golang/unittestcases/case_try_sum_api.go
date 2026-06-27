package unittestcases

import (
	"context"

	"github.com/starter-go/buckets"
	"github.com/starter-go/units"
)

type CaseTrySumAPI struct {

	//starter:component

	DM buckets.DriverManager //starter:inject("#")

	Service buckets.Service //starter:inject("#")

}

// ListRegistrations implements units.Unit.
func (inst *CaseTrySumAPI) ListRegistrations(list []*units.Registration) []*units.Registration {

	// old := list

	r1 := &units.Registration{
		Name:     "case-try-sum-api",
		Enabled:  false,
		Priority: 0,
		Do:       inst.run,
	}

	list = append(list, r1)

	return list
}

func (inst *CaseTrySumAPI) _impl() units.Unit {
	return inst
}

func (inst *CaseTrySumAPI) run(cc context.Context) error {
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// EOF
