package main4tencent
import (
    pa0a0e3aad "github.com/starter-go/bucket-drivers/tencent/lib/libcos"
     "github.com/starter-go/application"
)

// type pa0a0e3aad.TencentCosDriver in package:github.com/starter-go/bucket-drivers/tencent/lib/libcos
//
// id:com-a0a0e3aad71bd2f4-libcos-TencentCosDriver
// class:class-262c04a06c32904104382e2b8d56c279-DriverRegistry
// alias:
// scope:singleton
//
type pa0a0e3aad7_libcos_TencentCosDriver struct {
}

func (inst* pa0a0e3aad7_libcos_TencentCosDriver) register(cr application.ComponentRegistry) error {
	r := cr.NewRegistration()
	r.ID = "com-a0a0e3aad71bd2f4-libcos-TencentCosDriver"
	r.Classes = "class-262c04a06c32904104382e2b8d56c279-DriverRegistry"
	r.Aliases = ""
	r.Scope = "singleton"
	r.NewFunc = inst.new
	r.InjectFunc = inst.inject
	return r.Commit()
}

func (inst* pa0a0e3aad7_libcos_TencentCosDriver) new() any {
    return &pa0a0e3aad.TencentCosDriver{}
}

func (inst* pa0a0e3aad7_libcos_TencentCosDriver) inject(injext application.InjectionExt, instance any) error {
	ie := injext
	com := instance.(*pa0a0e3aad.TencentCosDriver)
	nop(ie, com)

	
    com.Enabled = inst.getEnabled(ie)
    com.Priority = inst.getPriority(ie)


    return nil
}


func (inst*pa0a0e3aad7_libcos_TencentCosDriver) getEnabled(ie application.InjectionExt)bool{
    return ie.GetBool("${buckets-driver.tencent-cos.enabled}")
}


func (inst*pa0a0e3aad7_libcos_TencentCosDriver) getPriority(ie application.InjectionExt)int{
    return ie.GetInt("${buckets-driver.tencent-cos.priority}")
}


