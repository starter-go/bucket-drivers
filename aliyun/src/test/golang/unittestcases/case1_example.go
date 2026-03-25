package unittestcases

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"strings"

	"github.com/starter-go/application"
	"github.com/starter-go/base/lang"
	"github.com/starter-go/buckets"
	"github.com/starter-go/vlog"
)

type ExampleCase struct {

	//starter:component

	DM buckets.DriverManager //starter:inject("#")

	Service buckets.Service //starter:inject("#")

}

func (inst *ExampleCase) _impl() application.Lifecycle {
	return inst
}

func (inst *ExampleCase) Life() *application.Life {
	l := &application.Life{
		OnLoop: inst.run,
	}
	return l
}

func (inst *ExampleCase) run() error {

	// dm := inst.DM

	ctx := context.Background()
	holder := new(buckets.BucketHolder)
	// obj1 := new(buckets.Object)

	holder.SetService(inst.Service)
	holder.SetName("demo1")
	holder.SetLazy(true)
	holder.SetContext(ctx)

	bucket1, err := holder.GetBucket()
	if err != nil {
		return err
	}

	taskset := &tryUseOssBucket{
		bucket: bucket1,
	}
	tasklist := make([]func() error, 0)

	tasklist = append(tasklist, taskset.init)
	tasklist = append(tasklist, taskset.doExists)
	tasklist = append(tasklist, taskset.doInsert)
	tasklist = append(tasklist, taskset.doExists)
	tasklist = append(tasklist, taskset.doFetchMeta)
	tasklist = append(tasklist, taskset.doUpdate)
	tasklist = append(tasklist, taskset.doFetchMeta)
	tasklist = append(tasklist, taskset.doFetchData)
	tasklist = append(tasklist, taskset.doDelete)
	tasklist = append(tasklist, taskset.doExists)

	for i, fn := range tasklist {

		vlog.Trace("step: %v", i)

		err := fn()
		if err != nil {
			return err
		}
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////

type tryUseOssBucket struct {
	oName  buckets.ObjectName
	bucket buckets.Bucket
}

func (inst *tryUseOssBucket) init() error {

	now := lang.Now()
	name := fmt.Sprintf(".test/oss/client/demo1/object-%v.txt", now.Int())
	inst.oName = buckets.ObjectName(name)
	return nil
}

func (inst *tryUseOssBucket) prepareObject() *buckets.Object {
	o := &buckets.Object{
		Name: inst.oName,
	}
	return o
}

func (inst *tryUseOssBucket) prepareObjectData(index int) []byte {

	b := &strings.Builder{}
	now := lang.Now()

	b.WriteString("hello,bucket")
	b.WriteString(fmt.Sprintf(".index=%d", index))
	b.WriteString(fmt.Sprintf(".now=%v", now))

	str := b.String()
	return []byte(str)
}

func (inst *tryUseOssBucket) checkDataSum(o *buckets.Object) {

	src := o.Data
	if src == nil {
		return
	}
	defer src.Close()

	data, err := io.ReadAll(src)

	if err != nil {
		vlog.Error("checkDataSum:error: %s", err.Error())
		return
	}

	sum := sha256.Sum256(data)
	hex := lang.HexFromBytes(sum[:])

	vlog.Info("object.data.sha256sum = %v", hex.String())

}

func (inst *tryUseOssBucket) doInsert() error {

	vlog.Trace("tryUseOssBucket.doInsert()")

	data := inst.prepareObjectData(1)
	o1 := inst.prepareObject()
	src := bytes.NewReader(data)
	bucket := inst.bucket

	o1.Data = io.NopCloser(src)
	o2, err := bucket.Put(o1)
	if err != nil {
		return err
	}

	vlog.Info("insert object : %v", o2.Sum)
	vlog.Info("  object.name = %v", o2.Name)

	inst.checkDataSum(o2)

	return nil
}

func (inst *tryUseOssBucket) doUpdate() error {

	vlog.Trace("tryUseOssBucket.doUpdate()")

	data := inst.prepareObjectData(2)
	o1 := inst.prepareObject()
	src := bytes.NewReader(data)
	bucket := inst.bucket

	o1.Data = io.NopCloser(src)
	o2, err := bucket.Put(o1)
	if err != nil {
		return err
	}

	// vlog.Info("update object : %v", o2.Sum)

	inst.checkDataSum(o2)

	return nil
}

func (inst *tryUseOssBucket) doDelete() error {

	vlog.Trace("tryUseOssBucket.doDelete()")

	o1 := inst.prepareObject()
	bucket := inst.bucket
	return bucket.Delete(o1)
}

func (inst *tryUseOssBucket) doFetchMeta() error {

	vlog.Trace("tryUseOssBucket.doFetchMeta()")

	o1 := inst.prepareObject()
	bucket := inst.bucket

	o2, err := bucket.GetMeta(o1)
	if err != nil {
		return err
	}

	inst.checkDataSum(o2)

	return nil
}

func (inst *tryUseOssBucket) doFetchData() error {

	vlog.Trace("tryUseOssBucket.doFetchData()")

	o1 := inst.prepareObject()
	bucket := inst.bucket

	o2, err := bucket.Fetch(o1)
	if err != nil {
		return err
	}

	inst.checkDataSum(o2)

	return nil
}

func (inst *tryUseOssBucket) doExists() error {

	obj := inst.prepareObject()
	ok, err := inst.bucket.Exists(obj)

	if err != nil {
		return err
	}

	vlog.Info("oss.client.ExistObject() : %v", ok)
	return nil
}

////////////////////////////////////////////////////////////////////////////////
