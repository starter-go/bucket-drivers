package main4aliyun
import (
    p93fc28538 "github.com/starter-go/bucket-drivers/aliyun/lib/liboss"
     "github.com/starter-go/application"
)

// type p93fc28538.AliyunOssDriver in package:github.com/starter-go/bucket-drivers/aliyun/lib/liboss
//
// id:com-93fc28538a0fa181-liboss-AliyunOssDriver
// class:class-262c04a06c32904104382e2b8d56c279-DriverRegistry
// alias:
// scope:singleton
//
type p93fc28538a_liboss_AliyunOssDriver struct {
}

func (inst* p93fc28538a_liboss_AliyunOssDriver) register(cr application.ComponentRegistry) error {
	r := cr.NewRegistration()
	r.ID = "com-93fc28538a0fa181-liboss-AliyunOssDriver"
	r.Classes = "class-262c04a06c32904104382e2b8d56c279-DriverRegistry"
	r.Aliases = ""
	r.Scope = "singleton"
	r.NewFunc = inst.new
	r.InjectFunc = inst.inject
	return r.Commit()
}

func (inst* p93fc28538a_liboss_AliyunOssDriver) new() any {
    return &p93fc28538.AliyunOssDriver{}
}

func (inst* p93fc28538a_liboss_AliyunOssDriver) inject(injext application.InjectionExt, instance any) error {
	ie := injext
	com := instance.(*p93fc28538.AliyunOssDriver)
	nop(ie, com)

	


    return nil
}


