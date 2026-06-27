package test4tencent
import (
    p8ae6d2a6c "github.com/starter-go/bucket-drivers/tencent/src/test/golang/unittestcases"
    p262c04a06 "github.com/starter-go/buckets"
    p0dc072ed4 "github.com/starter-go/units"
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



// type p8ae6d2a6c.CaseTryCRUD in package:github.com/starter-go/bucket-drivers/tencent/src/test/golang/unittestcases
//
// id:com-8ae6d2a6ced684e6-unittestcases-CaseTryCRUD
// class:
// alias:
// scope:singleton
//
type p8ae6d2a6ce_unittestcases_CaseTryCRUD struct {
}

func (inst* p8ae6d2a6ce_unittestcases_CaseTryCRUD) register(cr application.ComponentRegistry) error {
	r := cr.NewRegistration()
	r.ID = "com-8ae6d2a6ced684e6-unittestcases-CaseTryCRUD"
	r.Classes = ""
	r.Aliases = ""
	r.Scope = "singleton"
	r.NewFunc = inst.new
	r.InjectFunc = inst.inject
	return r.Commit()
}

func (inst* p8ae6d2a6ce_unittestcases_CaseTryCRUD) new() any {
    return &p8ae6d2a6c.CaseTryCRUD{}
}

func (inst* p8ae6d2a6ce_unittestcases_CaseTryCRUD) inject(injext application.InjectionExt, instance any) error {
	ie := injext
	com := instance.(*p8ae6d2a6c.CaseTryCRUD)
	nop(ie, com)

	
    com.DM = inst.getDM(ie)
    com.Service = inst.getService(ie)


    return nil
}


func (inst*p8ae6d2a6ce_unittestcases_CaseTryCRUD) getDM(ie application.InjectionExt)p262c04a06.DriverManager{
    return ie.GetComponent("#alias-262c04a06c32904104382e2b8d56c279-DriverManager").(p262c04a06.DriverManager)
}


func (inst*p8ae6d2a6ce_unittestcases_CaseTryCRUD) getService(ie application.InjectionExt)p262c04a06.Service{
    return ie.GetComponent("#alias-262c04a06c32904104382e2b8d56c279-Service").(p262c04a06.Service)
}



// type p8ae6d2a6c.CaseTryFileAPI in package:github.com/starter-go/bucket-drivers/tencent/src/test/golang/unittestcases
//
// id:com-8ae6d2a6ced684e6-unittestcases-CaseTryFileAPI
// class:
// alias:
// scope:singleton
//
type p8ae6d2a6ce_unittestcases_CaseTryFileAPI struct {
}

func (inst* p8ae6d2a6ce_unittestcases_CaseTryFileAPI) register(cr application.ComponentRegistry) error {
	r := cr.NewRegistration()
	r.ID = "com-8ae6d2a6ced684e6-unittestcases-CaseTryFileAPI"
	r.Classes = ""
	r.Aliases = ""
	r.Scope = "singleton"
	r.NewFunc = inst.new
	r.InjectFunc = inst.inject
	return r.Commit()
}

func (inst* p8ae6d2a6ce_unittestcases_CaseTryFileAPI) new() any {
    return &p8ae6d2a6c.CaseTryFileAPI{}
}

func (inst* p8ae6d2a6ce_unittestcases_CaseTryFileAPI) inject(injext application.InjectionExt, instance any) error {
	ie := injext
	com := instance.(*p8ae6d2a6c.CaseTryFileAPI)
	nop(ie, com)

	
    com.DM = inst.getDM(ie)
    com.Service = inst.getService(ie)
    com.DirMan = inst.getDirMan(ie)


    return nil
}


func (inst*p8ae6d2a6ce_unittestcases_CaseTryFileAPI) getDM(ie application.InjectionExt)p262c04a06.DriverManager{
    return ie.GetComponent("#alias-262c04a06c32904104382e2b8d56c279-DriverManager").(p262c04a06.DriverManager)
}


func (inst*p8ae6d2a6ce_unittestcases_CaseTryFileAPI) getService(ie application.InjectionExt)p262c04a06.Service{
    return ie.GetComponent("#alias-262c04a06c32904104382e2b8d56c279-Service").(p262c04a06.Service)
}


func (inst*p8ae6d2a6ce_unittestcases_CaseTryFileAPI) getDirMan(ie application.InjectionExt)p0dc072ed4.DirManager{
    return ie.GetComponent("#alias-0dc072ed44b3563882bff4e657a52e62-DirManager").(p0dc072ed4.DirManager)
}



// type p8ae6d2a6c.CaseTrySumAPI in package:github.com/starter-go/bucket-drivers/tencent/src/test/golang/unittestcases
//
// id:com-8ae6d2a6ced684e6-unittestcases-CaseTrySumAPI
// class:
// alias:
// scope:singleton
//
type p8ae6d2a6ce_unittestcases_CaseTrySumAPI struct {
}

func (inst* p8ae6d2a6ce_unittestcases_CaseTrySumAPI) register(cr application.ComponentRegistry) error {
	r := cr.NewRegistration()
	r.ID = "com-8ae6d2a6ced684e6-unittestcases-CaseTrySumAPI"
	r.Classes = ""
	r.Aliases = ""
	r.Scope = "singleton"
	r.NewFunc = inst.new
	r.InjectFunc = inst.inject
	return r.Commit()
}

func (inst* p8ae6d2a6ce_unittestcases_CaseTrySumAPI) new() any {
    return &p8ae6d2a6c.CaseTrySumAPI{}
}

func (inst* p8ae6d2a6ce_unittestcases_CaseTrySumAPI) inject(injext application.InjectionExt, instance any) error {
	ie := injext
	com := instance.(*p8ae6d2a6c.CaseTrySumAPI)
	nop(ie, com)

	
    com.DM = inst.getDM(ie)
    com.Service = inst.getService(ie)


    return nil
}


func (inst*p8ae6d2a6ce_unittestcases_CaseTrySumAPI) getDM(ie application.InjectionExt)p262c04a06.DriverManager{
    return ie.GetComponent("#alias-262c04a06c32904104382e2b8d56c279-DriverManager").(p262c04a06.DriverManager)
}


func (inst*p8ae6d2a6ce_unittestcases_CaseTrySumAPI) getService(ie application.InjectionExt)p262c04a06.Service{
    return ie.GetComponent("#alias-262c04a06c32904104382e2b8d56c279-Service").(p262c04a06.Service)
}


