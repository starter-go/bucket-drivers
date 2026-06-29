package liboss

import (
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"hash"
	"io"
	"strings"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/starter-go/afs"
	"github.com/starter-go/base/lang"
	"github.com/starter-go/buckets"
	"github.com/starter-go/vlog"
)

// 参考:
//   https://help.aliyun.com/zh/oss/developer-reference/manual-for-go-sdk-v2/

////////////////////////////////////////////////////////////////////////////////

type innerOssBucket struct {
	context    context.Context
	client     *oss.Client
	config     *oss.Config
	bucketName string
}

// ForFiles implements buckets.Bucket.
func (inst *innerOssBucket) ForFiles() buckets.BucketFileAPI {
	api := new(innerOssFileApi)
	api.bucket = inst
	return api
}

// ForSum implements buckets.Bucket.
func (inst *innerOssBucket) ForSum() buckets.BucketNativeSumAPI {
	api := new(innerOssNativeSumApi)
	api.bucket = inst
	return api
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

	sum := buckets.SUM{}
	ctx := inst.context
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
	o.Sum = inst.innerReadSumFromETag(res.ETag)

	return o, nil
}

func (inst *innerOssBucket) innerPrepareWriteMeta(src buckets.MetaMap) map[string]string {
	const empty = ""
	dst := make(map[string]string)
	for k, v := range src {
		name := k.String()
		if name == empty || v == empty {
			continue
		}
		dst[name] = v
	}
	return dst
}

func (inst *innerOssBucket) Put(o *buckets.Object) (*buckets.Object, error) {

	client := inst.client
	body := o.Data
	bucketName := inst.bucketName
	objectName := o.Name.String()
	sc := oss.StorageClassStandard
	acl := oss.ObjectACLPrivate

	ctx := inst.innerPrepareContext(o.Context)
	meta := inst.innerPrepareWriteMeta(o.Meta)

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
	client := inst.client

	req := &oss.GetObjectMetaRequest{
		Bucket: oss.Ptr(bucketName),
		Key:    oss.Ptr(objectName),
	}

	res, err := client.GetObjectMeta(ctx, req)
	if err != nil {
		return nil, err
	}

	size := res.ContentLength

	// if vlog.IsDebugEnabled() {
	// 	vlog.Debug("object.etag = %s", etag)
	// }

	o.Size = size
	o.Sum = inst.innerReadSumFromETag(res.ETag)

	return o, nil
}

// SetMeta implements buckets.Bucket.
func (inst *innerOssBucket) SetMeta(o *buckets.Object) (*buckets.Object, error) {

	vlog.Warn("innerOssBucket.SetMeta() : unsupported")

	// return nil, fmt.Errorf("no impl")
	return o, nil
}

func (inst *innerOssBucket) Exists(o *buckets.Object) (bool, error) {

	client := inst.client
	bucketName := inst.bucketName
	objectName := o.Name.String()
	ctx := inst.innerPrepareContext(o.Context)

	return client.IsObjectExist(ctx, bucketName, objectName)
}

func (inst *innerOssBucket) innerReadSumFromETag(etagPtr *string) buckets.SUM {

	sum := buckets.SUM{}
	etag := ""

	if etagPtr == nil {
		return sum
	}

	etag = *etagPtr
	etag = strings.ReplaceAll(etag, "\"", "")
	etag = strings.TrimSpace(etag)
	etag = strings.ToLower(etag)

	sum.Algorithm = buckets.AlgorithmMD5
	sum.Value = lang.Hex(etag)

	return sum
}

////////////////////////////////////////////////////////////////////////////////

type innerOssNativeSumApi struct {
	bucket *innerOssBucket
}

// Algorithm implements buckets.BucketNativeSumAPI.
func (i *innerOssNativeSumApi) Algorithm() buckets.CheckSumAlgorithm {
	return buckets.AlgorithmMD5
}

// Bucket implements buckets.BucketNativeSumAPI.
func (i *innerOssNativeSumApi) Bucket() buckets.Bucket {
	return i.bucket
}

// Hash implements buckets.BucketNativeSumAPI.
func (i *innerOssNativeSumApi) Hash() hash.Hash {
	return md5.New()
}

////////////////////////////////////////////////////////////////////////////////

type innerOssFileApi struct {
	bucket *innerOssBucket
}

// Bucket implements buckets.BucketFileAPI.
func (i *innerOssFileApi) Bucket() buckets.Bucket {
	return i.bucket
}

func (i *innerOssFileApi) innerMakeDirForFile(file afs.Path) error {

	dir := file.GetParent()
	if dir.Exists() {
		return nil
	}
	om := new(afs.OptionsMaker)
	om.Reset().SetMode(7, 5, 5)
	opt := om.Options()

	return dir.Mkdirs(&opt)
}

// FetchFile implements buckets.BucketFileAPI.
func (i *innerOssFileApi) FetchFile(o1 *buckets.ObjectFile) (*buckets.ObjectFile, error) {

	ctx := o1.Context
	client := i.bucket.client
	dl := client.NewDownloader()
	file := o1.Path
	name := o1.Name
	bucketName := i.bucket.bucketName
	objectName := name.String()

	err := i.innerMakeDirForFile(file)
	if err != nil {
		return nil, err
	}

	ctx = i.bucket.innerPrepareContext(ctx)
	req := &oss.GetObjectRequest{
		Bucket: oss.Ptr(bucketName),
		Key:    oss.Ptr(objectName),
	}

	res, err := dl.DownloadFile(ctx, req, file.GetPath())
	if err != nil {
		return nil, err
	}

	if vlog.IsDebugEnabled() {
		count := res.Written
		vlog.Debug("innerOssFileApi.[FetchFile res_Written:%d]", count)
	}

	o2 := new(buckets.ObjectFile)
	o2.Path = o1.Path
	o2.Context = ctx
	o2.Name = name
	o2.Bucket = i.bucket
	o2.Size = res.Written
	o2.Existed = true

	return o2, nil
}

// PutFile implements buckets.BucketFileAPI.
func (i *innerOssFileApi) PutFile(o1 *buckets.ObjectFile) (*buckets.ObjectFile, error) {

	ctx := o1.Context
	client := i.bucket.client
	upper := client.NewUploader()
	path := o1.Path.GetPath()
	name := o1.Name
	bucketName := i.bucket.bucketName
	objectName := name.String()
	meta := i.bucket.innerPrepareWriteMeta(o1.Meta)

	sc := oss.StorageClassStandard
	acl := oss.ObjectACLPrivate

	req := &oss.PutObjectRequest{
		Bucket:       oss.Ptr(bucketName),
		Key:          oss.Ptr(objectName),
		StorageClass: sc,
		Acl:          acl,
		Metadata:     meta,
		Body:         nil,
	}

	ctx = i.bucket.innerPrepareContext(ctx)

	res, err := upper.UploadFile(ctx, req, path)
	if err != nil {
		return nil, err
	}

	if vlog.IsDebugEnabled() {
		status := res.Status
		etagPtr := res.ETag
		etag := ""
		if etagPtr != nil {
			etag = *etagPtr
		}
		vlog.Debug("[innerOssFileApi.PutFile status:'%s' etag:'%s']", status, etag)
	}

	o2 := &buckets.ObjectFile{}
	o2.Path = o1.Path
	o2.Name = o1.Name
	o2.Context = ctx
	o2.Bucket = i.bucket

	return o2, nil
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
