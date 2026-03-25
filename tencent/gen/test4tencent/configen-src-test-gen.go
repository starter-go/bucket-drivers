package test4tencent
import (
    p8ae6d2a6c "github.com/starter-go/bucket-drivers/tencent/src/test/golang/unittestcases"
    p262c04a06 "github.com/starter-go/buckets"
     "github.com/starter-go/application"
)

// type p8ae6d2a6c.ExampleCase in package:github.com/starter-go/bucket-drivers/tencent/src/test/golang/unittestcases
//
// id:com-8ae6d2a6ced684e6-unittestcases-ExampleCase
// class:
// alias:
// scope:singleton
//
type p8ae6d2a6ce_unittestcases_ExampleCase struct {
}

func (inst* p8ae6d2a6ce_unittestcases_ExampleCase) register(cr application.ComponentRegistry) error {
	r := cr.NewRegistration()
	r.ID = "com-8ae6d2a6ced684e6-unittestcases-ExampleCase"
	r.Classes = ""
	r.Aliases = ""
	r.Scope = "singleton"
	r.NewFunc = inst.new
	r.InjectFunc = inst.inject
	return r.Commit()
}

func (inst* p8ae6d2a6ce_unittestcases_ExampleCase) new() any {
    return &p8ae6d2a6c.ExampleCase{}
}

func (inst* p8ae6d2a6ce_unittestcases_ExampleCase) inject(injext application.InjectionExt, instance any) error {
	ie := injext
	com := instance.(*p8ae6d2a6c.ExampleCase)
	nop(ie, com)

	
    com.DM = inst.getDM(ie)
    com.Service = inst.getService(ie)


    return nil
}


func (inst*p8ae6d2a6ce_unittestcases_ExampleCase) getDM(ie application.InjectionExt)p262c04a06.DriverManager{
    return ie.GetComponent("#alias-262c04a06c32904104382e2b8d56c279-DriverManager").(p262c04a06.DriverManager)
}


func (inst*p8ae6d2a6ce_unittestcases_ExampleCase) getService(ie application.InjectionExt)p262c04a06.Service{
    return ie.GetComponent("#alias-262c04a06c32904104382e2b8d56c279-Service").(p262c04a06.Service)
}


