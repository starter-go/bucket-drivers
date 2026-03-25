package test4aliyun
import (
    pa384215c1 "github.com/starter-go/bucket-drivers/aliyun/src/test/golang/unittestcases"
    p262c04a06 "github.com/starter-go/buckets"
     "github.com/starter-go/application"
)

// type pa384215c1.ExampleCase in package:github.com/starter-go/bucket-drivers/aliyun/src/test/golang/unittestcases
//
// id:com-a384215c1b53a20b-unittestcases-ExampleCase
// class:
// alias:
// scope:singleton
//
type pa384215c1b_unittestcases_ExampleCase struct {
}

func (inst* pa384215c1b_unittestcases_ExampleCase) register(cr application.ComponentRegistry) error {
	r := cr.NewRegistration()
	r.ID = "com-a384215c1b53a20b-unittestcases-ExampleCase"
	r.Classes = ""
	r.Aliases = ""
	r.Scope = "singleton"
	r.NewFunc = inst.new
	r.InjectFunc = inst.inject
	return r.Commit()
}

func (inst* pa384215c1b_unittestcases_ExampleCase) new() any {
    return &pa384215c1.ExampleCase{}
}

func (inst* pa384215c1b_unittestcases_ExampleCase) inject(injext application.InjectionExt, instance any) error {
	ie := injext
	com := instance.(*pa384215c1.ExampleCase)
	nop(ie, com)

	
    com.DM = inst.getDM(ie)
    com.Service = inst.getService(ie)


    return nil
}


func (inst*pa384215c1b_unittestcases_ExampleCase) getDM(ie application.InjectionExt)p262c04a06.DriverManager{
    return ie.GetComponent("#alias-262c04a06c32904104382e2b8d56c279-DriverManager").(p262c04a06.DriverManager)
}


func (inst*pa384215c1b_unittestcases_ExampleCase) getService(ie application.InjectionExt)p262c04a06.Service{
    return ie.GetComponent("#alias-262c04a06c32904104382e2b8d56c279-Service").(p262c04a06.Service)
}


