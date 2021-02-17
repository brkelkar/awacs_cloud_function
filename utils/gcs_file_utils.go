package utils

import (
	"context"
	"log"
	"strconv"
	"strings"
	"time"

	"awacs.com/awcacs_cloud_function/models"
	gc "github.com/brkelkar/common_utils/gcsbucketclient"
	"github.com/brkelkar/common_utils/logger"
	"go.uber.org/zap"
)

//GcsFile gcs file attributes
type GcsFile struct {
	FileName        string
	FilePath        string
	BucketName      string
	DistributorCode string
	LastUpdateTime  time.Time
	ProcessingTime  string
	Records         int
	FileType        string
	FileSize        int
	ErrorMsg        string
	GcsClient       *gc.GcsBucketClient
}

//HandleGCSEvent  parse file name and set all required attributes for the file
func (g *GcsFile) HandleGCSEvent(ctx context.Context, e models.GCSEvent) *GcsFile {

	var gcsObj gc.GcsBucketClient
	g.GcsClient = gcsObj.InitClient(ctx).SetBucketName(e.Bucket).SetNewReader(e.Name)

	if !g.GcsClient.GetLastStatus() {
		log.Print("Error while reading file")
	}
	g.FileSize, _ = strconv.Atoi(e.Size)
	g.FilePath = e.Bucket + "/" + e.Name
	g.FileName = e.Name
	g.BucketName = e.Bucket
	fileSplitSlice := strings.Split(e.Name, "/")
	g.DistributorCode = fileSplitSlice[0]
	g.LastUpdateTime = e.Updated
	g.ProcessingTime = e.Updated.Format("2006-01-02")
	return g
}

//LogFileDetails file details logger
func (g *GcsFile) LogFileDetails(status bool) {
	logger.Info("CF", zap.String("distributor_code", g.DistributorCode),
		zap.String("FileName", g.FileName),
		zap.String("FileSize", g.FileName),
		zap.String("FileType", g.FileType),
		zap.String("ProcessingTime", g.ProcessingTime),
		zap.Bool("Proting_status", status),
		zap.String("ErrorMsg", g.ErrorMsg),
		zap.Int("record_count", g.Records))
}
