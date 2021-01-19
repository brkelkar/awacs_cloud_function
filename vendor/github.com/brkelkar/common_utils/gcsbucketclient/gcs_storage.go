package gcsbucketclient

import (
	"context"
	"io"

	"cloud.google.com/go/storage"
)

var err error

//GcsBucketClient holds all required attributes to do gsc bucket operation
type GcsBucketClient struct {
	client     *storage.Client
	bucketName string
	w          io.Writer
	r          io.Reader
	ctx        context.Context
	// status indicates that one or more of the demo steps failed.
	status bool
}

func (g *GcsBucketClient) errorf(format string, args ...interface{}) {
	g.status = false
}

//InitClient intilialzie gcs client
func (g *GcsBucketClient) InitClient(ctx ...context.Context) *GcsBucketClient {
	if ctx == nil {
		g.ctx = context.Background()
	} else {

		g.ctx = ctx[0]
	}

	g.client, err = storage.NewClient(g.ctx)
	if err != nil {
		g.errorf("failed to get default GCS client: %v", err)

	} else {
		g.status = true
	}
	return g

}

//SetBucketName sets bucketname
func (g *GcsBucketClient) SetBucketName(bucket string) *GcsBucketClient {
	g.bucketName = bucket
	return g
}

//SetNewReader populates reader in struct
func (g *GcsBucketClient) SetNewReader(object string) *GcsBucketClient {

	g.r, err = g.client.Bucket(g.bucketName).Object(object).NewReader(g.ctx)
	if err != nil {
		g.errorf("failed to get default GCS bucket name: %v", err)
	} else {
		g.status = true

	}
	return g
}

//MoveObject moves obejct from one bucket to other
func (g *GcsBucketClient) MoveObject(srcObject string, destObject string, destBucket string) *GcsBucketClient {

	src := g.client.Bucket(g.bucketName).Object(srcObject)
	dst := g.client.Bucket(destBucket).Object(destObject)

	if _, err := dst.CopierFrom(src).Run(g.ctx); err != nil {
		g.errorf("failed to get default GCS bucket name: %v", err)
	} else {
		g.status = true

	}
	if err := src.Delete(g.ctx); err != nil {
		g.errorf("failed to get default GCS bucket name: %v", err)
	} else {
		g.status = true

	}
	return g
}

//GetReader expose reader
func (g *GcsBucketClient) GetReader() io.Reader {
	return g.r
}

//GetLastStatus expose reader
func (g *GcsBucketClient) GetLastStatus() bool {
	return g.status
}
