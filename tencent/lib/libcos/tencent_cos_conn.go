package libcos

import (
	"context"
	"crypto/md5"
	"fmt"
	"hash"
	"strings"

	cos "github.com/tencentyun/cos-go-sdk-v5"

	"github.com/starter-go/afs"
	"github.com/starter-go/base/lang"
	"github.com/starter-go/buckets"
	"github.com/starter-go/vlog"
)

// 参考文档: https://cloud.tencent.com/document/product/436/31215

////////////////////////////////////////////////////////////////////////////////

type innerCosBucketBuilder struct {
	context    context.Context
	bucketName string
	urls       *cos.BaseURL
	client     *cos.Client
}

func (inst *innerCosBucketBuilder) build() (buckets.Bucket, error) {

	b := new(innerCosBucket)

	b.context = inst.context
	b.bucketName = inst.bucketName
	b.urls = inst.urls
	b.client = inst.client

	return b, nil

}

////////////////////////////////////////////////////////////////////////////////

type innerCosBucket struct {
	context    context.Context
	bucketName string
	urls       *cos.BaseURL
	client     *cos.Client
}

// ForSum implements buckets.Bucket.
func (inst *innerCosBucket) ForSum() buckets.BucketNativeSumAPI {
	fapi := new(innerCosBucketSumApi)
	fapi.bucket = inst
	return fapi
}

// ForFiles implements buckets.Bucket.
func (inst *innerCosBucket) ForFiles() buckets.BucketFileAPI {
	fapi := new(innerCosBucketFileApi)
	fapi.bucket = inst
	return fapi
}

// Delete implements buckets.Bucket.
func (inst *innerCosBucket) Delete(o *buckets.Object) error {

	ctx := o.Context
	name := o.Name.String()
	client := inst.client

	ctx = inst.innerPrepareContext(ctx)
	resp, err := client.Object.Delete(ctx, name)
	if err != nil {
		return err
	}

	if vlog.IsDebugEnabled() {
		status := resp.Status
		vlog.Debug("[COS do:'Object.Delete' status:'%s' name:'%s']", status, name)
	}

	return nil
}

// Exists implements buckets.Bucket.
func (inst *innerCosBucket) Exists(o *buckets.Object) (bool, error) {

	ctx := o.Context
	name := o.Name.String()
	client := inst.client

	ctx = inst.innerPrepareContext(ctx)

	return client.Object.IsExist(ctx, name)
}

// Fetch implements buckets.Bucket.
func (inst *innerCosBucket) Fetch(o *buckets.Object) (*buckets.Object, error) {

	ctx := o.Context
	opt := &cos.ObjectGetOptions{}
	name := o.Name.String()
	client := inst.client

	ctx = inst.innerPrepareContext(ctx)
	resp, err := client.Object.Get(ctx, name, opt)
	if err != nil {
		return nil, err
	}

	data := resp.Body
	defer func() {
		if data != nil {
			data.Close()
		}
	}()

	if vlog.IsDebugEnabled() {
		// etag := result.ETag
		status := resp.Status
		vlog.Debug("[COS do:'Object.Put' status:'%s' name:'%s']", status, name)
	}

	err = inst.innerReadMeta(resp, o)
	if err != nil {
		return nil, err
	}

	o.Data = data
	data = nil

	return o, nil
}

// GetContext implements buckets.Bucket.
func (inst *innerCosBucket) GetContext() context.Context {
	c := inst.context
	if c == nil {
		c = context.Background()
	}
	return c
}

// GetMeta implements buckets.Bucket.
func (inst *innerCosBucket) GetMeta(o *buckets.Object) (*buckets.Object, error) {

	ctx := o.Context
	opt := &cos.ObjectHeadOptions{}
	name := o.Name.String()
	client := inst.client

	ctx = inst.innerPrepareContext(ctx)
	resp, err := client.Object.Head(ctx, name, opt)
	if err != nil {
		return nil, err
	}

	err = inst.innerReadMeta(resp, o)
	if err != nil {
		return nil, err
	}

	return o, nil
}

// SetMeta implements buckets.Bucket.
func (inst *innerCosBucket) SetMeta(o *buckets.Object) (*buckets.Object, error) {

	// return nil, fmt.Errorf("unsupported")

	vlog.Warn("innerCosBucket.SetMeta() : unsupported")

	return o, nil
}

func (inst *innerCosBucket) innerReadMeta(resp *cos.Response, dst *buckets.Object) error {

	name := dst.Name
	contentType := resp.Header.Get("Content-Type")
	contentLength := resp.Header.Get("Content-Length")
	etag := resp.Header.Get("ETag")
	reqid := resp.Header.Get("X-Cos-Request-Id")

	etag = strings.ReplaceAll(etag, "\"", "")
	etag = strings.TrimSpace(etag)

	if vlog.IsDebugEnabled() {
		// etag := result.ETag
		status := resp.Status
		vlog.Debug("[COS do:'Object.Head(GetMeta)' status:'%s' name:'%s']", status, name)
	}

	meta := buckets.MetaMap{
		"content-type":   contentType,
		"content-length": contentLength,
		"etag":           etag,
		"request-id":     reqid,
	}

	meta[buckets.MetaObjectType] = contentType
	meta[buckets.MetaObjectLength] = contentLength
	meta[buckets.MetaObjectSum] = "md5(" + etag + ")"

	dst.Meta = meta
	dst.Sum.Algorithm = "md5"
	dst.Sum.Value = lang.Hex(etag)

	return nil
}

// GetObject implements buckets.Bucket.
func (inst *innerCosBucket) GetObject(name buckets.ObjectName) *buckets.Object {
	obj := new(buckets.Object)
	obj.Name = name
	obj.Context = inst.context
	obj.Bucket = inst
	return obj
}

// Put implements buckets.Bucket.
func (inst *innerCosBucket) Put(o *buckets.Object) (*buckets.Object, error) {

	ctx := o.Context
	opt := &cos.ObjectPutOptions{}
	name := o.Name.String()
	src := o.Data
	client := inst.client

	ctx = inst.innerPrepareContext(ctx)
	resp, err := client.Object.Put(ctx, name, src, opt)
	if err != nil {
		return nil, err
	}

	if vlog.IsDebugEnabled() {
		// etag := result.ETag
		status := resp.Status
		vlog.Debug("[COS do:'Object.Put' status:'%s' name:'%s']", status, name)
	}

	return o, nil
}

// SetContext implements buckets.Bucket.
func (inst *innerCosBucket) SetContext(ctx context.Context) buckets.Bucket {
	if ctx == nil {
		ctx = context.Background()
	}
	inst.context = ctx
	return inst
}

func (inst *innerCosBucket) innerPrepareContext(c context.Context) context.Context {
	if c == nil {
		c = inst.context
	}
	if c == nil {
		c = context.Background()
	}
	return c
}

func (inst *innerCosBucket) _impl() buckets.Bucket {
	return inst
}

////////////////////////////////////////////////////////////////////////////////

type innerCosBucketSumApi struct {
	bucket *innerCosBucket

	// algorithm = md5
}

// Algorithm implements buckets.BucketNativeSumAPI.
func (inst *innerCosBucketSumApi) Algorithm() buckets.CheckSumAlgorithm {
	return buckets.AlgorithmMD5
}

// Bucket implements buckets.BucketNativeSumAPI.
func (inst *innerCosBucketSumApi) Bucket() buckets.Bucket {
	return inst.bucket
}

// Hash implements buckets.BucketNativeSumAPI.
func (inst *innerCosBucketSumApi) Hash() hash.Hash {
	return md5.New()
}

func (inst *innerCosBucketSumApi) _impl() buckets.BucketNativeSumAPI {
	return inst
}

////////////////////////////////////////////////////////////////////////////////

type innerCosBucketFileApi struct {
	bucket *innerCosBucket
}

// Bucket implements buckets.BucketFileAPI.
func (inst *innerCosBucketFileApi) Bucket() buckets.Bucket {
	return inst.bucket
}

// FetchFile implements buckets.BucketFileAPI.
func (inst *innerCosBucketFileApi) FetchFile(o *buckets.ObjectFile) (*buckets.ObjectFile, error) {

	ctx := o.Context
	opt := &cos.MultiDownloadOptions{}
	name := o.Name.String()
	file := o.Path.GetPath()
	client := inst.bucket.client

	err := inst.innerPrepareDirToFetchFile(o)
	if err != nil {
		return nil, err
	}

	ctx = inst.bucket.innerPrepareContext(ctx)
	resp, err := client.Object.Download(ctx, name, file, opt)
	if err != nil {
		return nil, err
	}

	if vlog.IsDebugEnabled() {
		// etag := result.ETag
		status := resp.Status
		vlog.Debug("[COS do:'Object.Download' status:'%s' name:'%s']", status, name)
	}

	err = inst.bucket.innerReadMeta(resp, &o.Object)
	if err != nil {
		return nil, err
	}

	return o, nil
}

func (inst *innerCosBucketFileApi) innerPrepareDirToFetchFile(o *buckets.ObjectFile) error {

	file := o.Path
	if file == nil {
		return fmt.Errorf("innerCosBucketFileApi: object-file is nil")
	}

	dir := file.GetParent()
	if dir.Exists() {
		return nil
	}

	om := new(afs.OptionsMaker)
	om.SetMode(7, 5, 5)
	opt := om.Options()

	return dir.Mkdirs(&opt)
}

// PutFile implements buckets.BucketFileAPI.
func (inst *innerCosBucketFileApi) PutFile(o *buckets.ObjectFile) (*buckets.ObjectFile, error) {

	ctx := o.Context
	opt := &cos.ObjectPutOptions{}
	name := o.Name.String()
	file := o.Path.GetPath()
	client := inst.bucket.client

	ctx = inst.bucket.innerPrepareContext(ctx)
	resp, err := client.Object.PutFromFile(ctx, name, file, opt)
	if err != nil {
		return nil, err
	}

	if vlog.IsDebugEnabled() {
		// etag := result.ETag
		status := resp.Status
		vlog.Debug("[COS do:'Object.PutFromFile' status:'%s' name:'%s']", status, name)
	}

	return o, nil
}

func (inst *innerCosBucketFileApi) _impl() buckets.BucketFileAPI {
	return inst
}

////////////////////////////////////////////////////////////////////////////////
// EOF
