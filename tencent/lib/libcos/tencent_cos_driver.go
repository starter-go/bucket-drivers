package libcos

import (
	"github.com/starter-go/application"
	"github.com/starter-go/buckets"
)

const theDriverName = "tencent-cos"

////////////////////////////////////////////////////////////////////////////////

type TencentCosDriver struct {

	//starter:component

	_as func(buckets.DriverRegistry) //starter:as(".")

	Enabled  bool //starter:inject("${bucket-driver.tencent-cos.enabled}")
	Priority int  //starter:inject("${bucket-driver.tencent-cos.priority}")
}

// Accept implements buckets.Driver.
func (inst *TencentCosDriver) Accept(cfg *buckets.Configuration) bool {
	if cfg == nil {
		return false
	}
	drName := cfg.Driver
	return (drName == theDriverName)
}

// GetLoader implements buckets.Driver.
func (inst *TencentCosDriver) GetLoader() buckets.Loader {
	loader := new(innerTencentCosLoader)
	loader.driver = inst
	return loader
}

// GetRegistration implements buckets.Driver.
func (inst *TencentCosDriver) GetRegistration() *buckets.DriverRegistration {
	r1 := &buckets.DriverRegistration{
		Name:     theDriverName,
		Enabled:  inst.Enabled,
		Priority: inst.Priority,
		Driver:   inst,
	}
	return r1
}

func (inst *TencentCosDriver) Life() *application.Life {

	l := new(application.Life)

	// l.OnCreate = inst.onCreate
	// l.OnStart = inst.onStart
	// l.OnLoop = inst.onLoop
	// l.OnStop = inst.onStop
	// l.OnDestroy = inst.onDestroy

	return l
}

func (inst *TencentCosDriver) ListDriverRegistrations() []*buckets.DriverRegistration {
	dr1 := inst.GetRegistration()
	return []*buckets.DriverRegistration{dr1}
}

func (inst *TencentCosDriver) _impl() (application.Lifecycle, buckets.DriverRegistry, buckets.Driver) {
	return inst, inst, inst
}

////////////////////////////////////////////////////////////////////////////////

type innerTencentCosLoader struct {
	driver *TencentCosDriver
}

// Open implements buckets.Loader.
func (i *innerTencentCosLoader) Open(cfg *buckets.Configuration, options *buckets.OpenOptions) (buckets.Bucket, error) {

	bucket := new(innerTencentCosBucket)

	err := bucket.config(cfg, options)
	if err != nil {
		return nil, err
	}

	err = bucket.open()
	if err != nil {
		return nil, err
	}

	return bucket, nil
}

////////////////////////////////////////////////////////////////////////////////
