package libcos

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/starter-go/application"
	"github.com/starter-go/buckets"
	"github.com/starter-go/vlog"
	cos "github.com/tencentyun/cos-go-sdk-v5"
)

const theDriverName = "tencent-cos"

////////////////////////////////////////////////////////////////////////////////

type TencentCosDriver struct {

	//starter:component

	_as func(buckets.DriverRegistry) //starter:as(".")

	Enabled  bool //starter:inject("${buckets-driver.tencent-cos.enabled}")
	Priority int  //starter:inject("${buckets-driver.tencent-cos.priority}")
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
	loader := &innerCosLoader{
		driver: inst,
	}
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

type innerCosLoader struct {
	driver *TencentCosDriver
}

// Open implements buckets.Loader.
func (inst *innerCosLoader) Open(cfg *buckets.Configuration, options *buckets.OpenOptions) (buckets.Bucket, error) {

	loading := new(innerCosLoading)
	loading.config = cfg
	loading.options = options

	steps := make([]func(loading *innerCosLoading) error, 0)

	steps = append(steps, inst.innerDoConfigure)
	steps = append(steps, inst.innerDoOpenClient)
	steps = append(steps, inst.innerDoMakeBucket)
	steps = append(steps, inst.innerDoCheckClient)

	for _, st := range steps {
		err := st(loading)
		if err != nil {
			return nil, err
		}
	}

	return loading.getBucket()
}

func (inst *innerCosLoader) innerDoConfigure(loading *innerCosLoading) error {

	cfg := loading.config
	options := loading.options

	if cfg == nil {
		return fmt.Errorf("innerCosBucket: buckets.Configuration is nil")
	}
	if options == nil {
		options = new(buckets.OpenOptions)
		options.Timeout = time.Second * 30
	}

	ctx := options.Context
	if ctx == nil {
		ctx = context.Background()
	}

	at := &cos.AuthorizationTransport{
		SecretID:     cfg.AccessKeyID,
		SecretKey:    cfg.AccessKeySecret,
		SessionToken: "",
	}

	bucketURL, err := inst.innerGetBucketURL(cfg)
	if err != nil {
		return err
	}

	urlset := &cos.BaseURL{
		BucketURL: bucketURL,
	}

	loading.urlset = urlset
	loading.transport = at
	loading.bucketName = cfg.Name
	loading.context = ctx

	return nil
}

func (inst *innerCosLoader) innerGetBucketURL(cfg *buckets.Configuration) (*url.URL, error) {
	str := cfg.URL
	return url.Parse(str)
}

func (inst *innerCosLoader) innerDoCheckClient(loading *innerCosLoading) error {

	// 通过 'get-bucket-list' 来检验 client 是否ready

	ctx := context.Background()
	client := loading.client
	res, _, err := client.Service.Get(ctx)

	if err != nil {
		return err
	}

	list := res.Buckets
	for index, bucket := range list {
		name := bucket.Name
		vlog.Debug("  bucket[%d].name = %s", index, name)
	}

	return nil
}

func (inst *innerCosLoader) innerDoOpenClient(loading *innerCosLoading) error {

	urls := loading.urlset
	transport := loading.transport

	if urls == nil {
		return fmt.Errorf("innerCosBucket: urls is nil")
	}

	if transport == nil {
		return fmt.Errorf("innerCosBucket: transport is nil")
	}

	htclient := &http.Client{
		Transport: transport,
	}
	client := cos.NewClient(urls, htclient)
	loading.client = client
	return nil
}

func (i *innerCosLoader) innerDoMakeBucket(loading *innerCosLoading) error {

	b := new(innerCosBucketBuilder)

	// bucket := new(innerCosBucket)
	// bucket.client = loading.client
	// bucket.context = loading.context
	// bucket.bucketName = loading.bucketName
	// bucket.urls = loading.urlset

	b.client = loading.client
	b.context = loading.context
	b.bucketName = loading.bucketName
	b.urls = loading.urlset

	bucket, err := b.build()
	if err != nil {
		return err
	}

	loading.bucket = bucket
	return nil
}

////////////////////////////////////////////////////////////////////////////////

type innerCosLoading struct {
	config *buckets.Configuration

	options *buckets.OpenOptions

	bucket buckets.Bucket

	client *cos.Client

	urlset *cos.BaseURL

	transport http.RoundTripper

	bucketName string

	context context.Context
}

func (inst *innerCosLoading) getBucket() (buckets.Bucket, error) {
	b := inst.bucket
	if b == nil {
		return nil, fmt.Errorf("innerCosLoading: bucket is nil")
	}
	return b, nil
}

////////////////////////////////////////////////////////////////////////////////
// EOF
