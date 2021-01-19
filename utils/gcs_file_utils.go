package utils

import (
	"context"
	"log"
	"strings"
	"time"

	"awacs.com/awcacs_cloud_function/models"

	gc "github.com/brkelkar/common_utils/gcsbucketclient"
)

type GcsFile struct {
	FileName        string
	FilePath        string
	BucketName      string
	DistributorCode string
	LastUpdateTime  time.Time
	ProcessingTime  string
	Records         int
	GcsClient       *gc.GcsBucketClient
}

//HandleGCSEvent  parse file name and set all required attributes for the file
func (g *GcsFile) HandleGCSEvent(ctx context.Context, e models.GCSEvent) *GcsFile {

	var gcsObj gc.GcsBucketClient
	g.GcsClient = gcsObj.InitClient(ctx).SetBucketName(e.Bucket).SetNewReader(e.Name)

	if !g.GcsClient.GetLastStatus() {
		log.Print("Error while reading file")
	}
	g.FilePath = e.Bucket + "/" + e.Name
	g.FileName = e.Name
	g.BucketName = e.Bucket
	fileSplitSlice := strings.Split(e.Name, "/")
	g.DistributorCode = fileSplitSlice[0]
	g.LastUpdateTime = e.Updated
	g.ProcessingTime = e.Updated.Format("2006-01-02")
	return g
}
