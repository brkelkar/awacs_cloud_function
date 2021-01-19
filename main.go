package invoiceupload

import (
	"awacs.com/awcacs_cloud_function/functions"
	"awacs.com/awcacs_cloud_function/models"
	"awacs.com/awcacs_cloud_function/utils"

	"context"

	cr "github.com/brkelkar/common_utils/configreader"
	gc "github.com/brkelkar/common_utils/gcsbucketclient"

	"log"
	"strings"
)

var (
	dateFormatMap map[string]string
	err           error
	cfg           cr.Config
	gcsFileAttr   utils.GcsFile
	gcsObj        gc.GcsBucketClient
)

func init() {
	cfg.ReadGcsFile("gs://awacs_config/config.yml")
}

//SyncFileUpload cloud funtion to upload file
func SyncFileUpload(ctx context.Context, e models.GCSEvent) (err error) {

	log.Println("Porting Start File Name = " + e.Name)

	g := gcsFileAttr.HandleGCSEvent(ctx, e)
	log.Print(g)
	if strings.Contains(strings.ToLower(e.Bucket), "invoice") {
		log.Println("Calling Invoice upload method")
		var invoiceObj functions.InvoiceAttr
		err = invoiceObj.InvoiceCloudFunction(g, cfg)
		return
	}

	if strings.Contains(strings.ToLower(e.Bucket), "stock") {
		log.Println("Calling Stock upload method")
		var stockObj functions.StockAttr
		err = stockObj.StockCloudFunction(g, cfg)

	}

	return
}
