package liboss

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/starter-go/buckets"
	"github.com/starter-go/vlog"
)

////////////////////////////////////////////////////////////////////////////////

type innerOssBucket struct {
	context    context.Context
	client     *oss.Client
	config     *oss.Config
	bucketName string
}

// Delete implements buckets.Bucket.
func (inst *innerOssBucket) Delete(o *buckets.Object) error {

	bucketName := inst.bucketName
	objectName := o.Name.String()

	ctx := inst.innerPrepareContext(o.Context)
	req := &oss.DeleteObjectRequest{
		Bucket: oss.Ptr(bucketName),
		Key:    oss.Ptr(objectName),
	}

	client := inst.client
	res, err := client.DeleteObject(ctx, req)
	if err != nil {
		return err
	}

	if vlog.IsDebugEnabled() {
		vlog.Debug("oss.client.DeleteObject() : %s", res.Status)
	}

	return nil
}

// GetContext implements buckets.Bucket.
func (inst *innerOssBucket) GetContext() context.Context {

	ctx := inst.context
	if ctx == nil {
		ctx = context.Background()
		inst.context = ctx
	}
	return ctx
}

// SetContext implements buckets.Bucket.
func (inst *innerOssBucket) SetContext(ctx context.Context) buckets.Bucket {
	if ctx != nil {
		inst.context = ctx
	}
	return inst
}

func (inst *innerOssBucket) _impl() buckets.Bucket {
	return inst
}

func (inst *innerOssBucket) open(cfg *buckets.Configuration, options *buckets.OpenOptions) error {

	cfg2 := oss.LoadDefaultConfig()
	err := inst.innerPrepareConfig(cfg, cfg2, options)
	if err != nil {
		return err
	}

	client := oss.NewClient(cfg2)

	// inst.bucketName = cfg2.
	inst.client = client
	inst.config = cfg2

	return inst.innerCheck()
}

func (inst *innerOssBucket) innerCheck() error {

	const name = "tmp/test/object4check.txt"

	obj := inst.GetObject(name)
	ok, err := inst.Exists(obj)
	if err == nil {
		if ok {
			return nil
		}
	}

	// put a new object
	obj = inst.GetObject(name)
	data := bytes.NewBufferString("a test file for check")
	obj.Data = io.NopCloser(data)
	obj, err = inst.Put(obj)

	return err
}

func (inst *innerOssBucket) Close() error {

	// return fmt.Errorf("no impl")

	return nil
}

func (inst *innerOssBucket) innerResolveConfigURL(cfg *buckets.Configuration) error {

	l, err := buckets.ParseLocation(cfg.URL)
	if err != nil {
		return err
	}
	cfg.Location = *l
	return nil
}

func (inst *innerOssBucket) innerPrepareConfig(src *buckets.Configuration, dst *oss.Config, options *buckets.OpenOptions) error {

	err := inst.innerResolveConfigURL(src)
	if err != nil {
		return err
	}

	cp := new(innerCredentialsProvider)
	timeout := options.Timeout

	// keyID := src.AccessKeyID
	// keySecret := src.AccessKeySecret

	endpoint := src.Location.Host
	region := inst.innerGetConfigQuery(src, "region", true)
	bucketName := inst.innerGetConfigQuery(src, "bucket", true)

	cp.init(inst, src, options)
	dst.WithEndpoint(endpoint).WithRegion(region).WithCredentialsProvider(cp).WithReadWriteTimeout(timeout)
	inst.bucketName = bucketName

	return nil
}

func (inst *innerOssBucket) innerGetConfigQuery(cfg *buckets.Configuration, name string, required bool) string {

	q := cfg.Location.Query
	value := ""

	if q != nil {
		value = q[name]
	}

	if required && (value == "") {
		err := fmt.Errorf("innerOssBucket: no required URI.query item: '%s'", name)
		panic(err)
	}

	return value
}

func (inst *innerOssBucket) GetObject(name buckets.ObjectName) *buckets.Object {

	ctx := inst.context
	sum := buckets.SUM{}

	ctx = inst.innerPrepareContext(ctx)

	o1 := &buckets.Object{

		Context: ctx,
		Name:    name,
		Type:    "application/x-unknown",
		Sum:     sum,
		Size:    0,
		Data:    nil,
		Bucket:  inst,
		Existed: false,
	}

	return o1
}

func (inst *innerOssBucket) innerPrepareContext(ctx context.Context) context.Context {

	if ctx == nil {
		ctx = inst.GetContext()
	}

	if ctx == nil {
		ctx = context.Background()
	}

	return ctx
}

func (inst *innerOssBucket) Fetch(o *buckets.Object) (*buckets.Object, error) {

	client := inst.client
	bucketName := inst.bucketName
	objectName := o.Name.String()

	// body := o.Data
	// sc := oss.StorageClassStandard
	// acl := oss.ObjectACLPrivate

	ctx := inst.innerPrepareContext(o.Context)

	req := &oss.GetObjectRequest{
		Bucket: oss.Ptr(bucketName),
		Key:    oss.Ptr(objectName),
	}

	res, err := client.GetObject(ctx, req)
	if err != nil {
		return nil, err
	}

	etag := res.ETag
	vlog.Debug("object.etag = %s", *etag)

	o.Data = res.Body
	o.Bucket = inst
	o.Existed = true
	o.Size = res.ContentLength
	o.Type = buckets.ContentType(*res.ContentType)

	return o, nil
}

func (inst *innerOssBucket) Put(o *buckets.Object) (*buckets.Object, error) {

	client := inst.client
	body := o.Data
	bucketName := inst.bucketName
	objectName := o.Name.String()
	sc := oss.StorageClassStandard
	acl := oss.ObjectACLPrivate

	ctx := inst.innerPrepareContext(o.Context)

	meta := map[string]string{
		"sha256sum": "todo",
	}

	req := &oss.PutObjectRequest{
		Bucket:       oss.Ptr(bucketName),
		Key:          oss.Ptr(objectName),
		StorageClass: sc,
		Acl:          acl,
		Metadata:     meta,
		Body:         body,
	}

	res, err := client.PutObject(ctx, req)
	if err != nil {
		return nil, err
	}

	md5sum := res.ContentMD5
	vlog.Debug("content.md5 = %s", *md5sum)

	return o, nil
}

func (inst *innerOssBucket) GetMeta(o *buckets.Object) (*buckets.Object, error) {

	bucketName := inst.bucketName
	objectName := o.Name.String()

	ctx := inst.innerPrepareContext(o.Context)

	req := &oss.GetObjectMetaRequest{
		Bucket: oss.Ptr(bucketName),
		Key:    oss.Ptr(objectName),
	}

	client := inst.client
	res, err := client.GetObjectMeta(ctx, req)

	if err != nil {
		return nil, err
	}

	etag := res.ETag
	vlog.Debug("object.etag = %s", *etag)

	return o, nil
}

func (inst *innerOssBucket) Exists(o *buckets.Object) (bool, error) {

	client := inst.client
	bucketName := inst.bucketName
	objectName := o.Name.String()
	ctx := inst.innerPrepareContext(o.Context)

	return client.IsObjectExist(ctx, bucketName, objectName)
}

////////////////////////////////////////////////////////////////////////////////

type innerCredentialsProvider struct {

	// bucket *
	// cfg1   *buckets.Configuration

	cred credentials.Credentials
}

func (inst *innerCredentialsProvider) _impl() credentials.CredentialsProvider {
	return inst
}

func (inst *innerCredentialsProvider) init(bucket *innerOssBucket, cfg *buckets.Configuration, options *buckets.OpenOptions) {

	now := time.Now()
	age := time.Second * 30
	exp := now.Add(age)

	keyID := cfg.AccessKeyID
	keySec := cfg.AccessKeySecret
	stoken := bucket.innerGetConfigQuery(cfg, "security-token", true)

	stoken = ""

	cred := &credentials.Credentials{
		AccessKeyID:     keyID,
		AccessKeySecret: keySec,
		SecurityToken:   stoken,
		Expires:         &exp,
	}

	inst.cred = *cred
}

func (inst *innerCredentialsProvider) GetCredentials(ctx context.Context) (credentials.Credentials, error) {
	return inst.cred, nil
}

////////////////////////////////////////////////////////////////////////////////
