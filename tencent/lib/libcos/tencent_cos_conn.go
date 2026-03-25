package libcos

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	cos "github.com/tencentyun/cos-go-sdk-v5"

	"github.com/starter-go/buckets"
)

// 参考文档: https://cloud.tencent.com/document/product/436/31215

////////////////////////////////////////////////////////////////////////////////

type innerTencentCosBucket struct {
	context    context.Context
	bucketName string
	urls       *cos.BaseURL
	client     *cos.Client
	transport  *cos.AuthorizationTransport
}

func (inst *innerTencentCosBucket) config(cfg *buckets.Configuration, options *buckets.OpenOptions) error {

	if cfg == nil {
		return fmt.Errorf("innerTencentCosBucket: buckets.Configuration is nil")
	}

	if options == nil {
		options = new(buckets.OpenOptions)
		options.Timeout = time.Second * 30
	}

	ctx := options.Context
	if ctx == nil {
		ctx = context.Background()
	}

	transport := &cos.AuthorizationTransport{
		SecretID:  cfg.AccessKeyID,
		SecretKey: cfg.AccessKeySecret,
	}

	bucketURL, err := inst.innerGetBucketURL(cfg)
	if err != nil {
		return
	}

	urls := &cos.BaseURL{
		BucketURL: bucketURL,
	}

	inst.urls = urls
	inst.transport = transport
	inst.bucketName = cfg.Name
	inst.context = ctx
	return nil
}

func (inst *innerTencentCosBucket) open() error {

	client, err := inst.innerMakeNewClient()
	if err != nil {
		return err
	}

	inst.client = client
	return nil
}

func (inst *innerTencentCosBucket) close() error {
	panic("unimplemented")
}

func (inst *innerTencentCosBucket) innerGetBucketURL(cfg *buckets.Configuration) (*url.URL, error) {}

func (inst *innerTencentCosBucket) innerMakeNewClient() (*cos.Client, error) {

	urls := inst.urls
	transport := inst.transport

	if urls == nil {
		return nil, fmt.Errorf("innerTencentCosBucket: urls is nil")
	}

	if transport == nil {
		return nil, fmt.Errorf("innerTencentCosBucket: transport is nil")
	}

	htclient := &http.Client{
		Transport: transport,
	}
	client := cos.NewClient(urls, htclient)
	return client, nil
}

// Delete implements buckets.Bucket.
func (inst *innerTencentCosBucket) Delete(o *buckets.Object) error {
	panic("unimplemented")
}

// Exists implements buckets.Bucket.
func (inst *innerTencentCosBucket) Exists(o *buckets.Object) (bool, error) {
	panic("unimplemented")
}

// Fetch implements buckets.Bucket.
func (inst *innerTencentCosBucket) Fetch(o *buckets.Object) (*buckets.Object, error) {
	panic("unimplemented")
}

// GetContext implements buckets.Bucket.
func (inst *innerTencentCosBucket) GetContext() context.Context {
	panic("unimplemented")
}

// GetMeta implements buckets.Bucket.
func (inst *innerTencentCosBucket) GetMeta(o *buckets.Object) (*buckets.Object, error) {
	panic("unimplemented")
}

// GetObject implements buckets.Bucket.
func (inst *innerTencentCosBucket) GetObject(name buckets.ObjectName) *buckets.Object {
	panic("unimplemented")
}

// Put implements buckets.Bucket.
func (inst *innerTencentCosBucket) Put(o *buckets.Object) (*buckets.Object, error) {
	panic("unimplemented")
}

// SetContext implements buckets.Bucket.
func (inst *innerTencentCosBucket) SetContext(ctx context.Context) buckets.Bucket {
	panic("unimplemented")
}

func (inst *innerTencentCosBucket) _impl() buckets.Bucket {
	return inst
}

////////////////////////////////////////////////////////////////////////////////

////////////////////////////////////////////////////////////////////////////////
