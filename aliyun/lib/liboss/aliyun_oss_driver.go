package liboss

import (
	"github.com/starter-go/application"
	"github.com/starter-go/buckets"
)

const theDriverName = "aliyun-oss"

////////////////////////////////////////////////////////////////////////////////

type AliyunOssDriver struct {

	//starter:component

	_as func(buckets.DriverRegistry) //starter:as(".")

	Enabled  bool
	Priority int
}

func (inst *AliyunOssDriver) _impl() (application.Lifecycle, buckets.DriverRegistry, buckets.Driver) {
	return inst, inst, inst
}

func (inst *AliyunOssDriver) Life() *application.Life {

	l := new(application.Life)

	// l.OnCreate = inst.onCreate
	// l.OnStart = inst.onStart
	// l.OnLoop = inst.onLoop
	// l.OnStop = inst.onStop
	// l.OnDestroy = inst.onDestroy

	return l
}

func (inst *AliyunOssDriver) ListDriverRegistrations() []*buckets.DriverRegistration {
	dr1 := inst.GetRegistration()
	return []*buckets.DriverRegistration{dr1}
}

func (inst *AliyunOssDriver) GetLoader() buckets.Loader {

	return new(innerLoader)

}

func (inst *AliyunOssDriver) GetRegistration() *buckets.DriverRegistration {

	inst.Enabled = true
	inst.Priority = 1

	dr1 := new(buckets.DriverRegistration)

	dr1.Name = theDriverName
	dr1.Driver = inst
	dr1.Enabled = inst.Enabled
	dr1.Priority = inst.Priority

	return dr1
}

func (inst *AliyunOssDriver) Accept(cfg *buckets.Configuration) bool {
	if cfg == nil {
		return false
	}
	name := cfg.Driver
	return name == theDriverName
}

////////////////////////////////////////////////////////////////////////////////

type innerLoader struct {
}

func (inst *innerLoader) Open(cfg *buckets.Configuration, options *buckets.OpenOptions) (buckets.Bucket, error) {

	cfg, err := inst.innerPrepareConfig(cfg)
	if err != nil {
		return nil, err
	}

	bu := new(innerOssBucket)
	err = bu.open(cfg, options)
	if err != nil {
		return nil, err
	}
	return bu, nil
}

func (inst *innerLoader) innerPrepareConfig(cfg *buckets.Configuration) (*buckets.Configuration, error) {

	if cfg == nil {
		cfg = new(buckets.Configuration)
	}

	if cfg.MaxObjectSize == 0 {

	} else if cfg.MaxObjectSize < 0 {
	}

	return cfg, nil
}

////////////////////////////////////////////////////////////////////////////////
