package functions

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"awacs.com/awcacs_cloud_function/models"
	"awacs.com/awcacs_cloud_function/utils"
	bt "github.com/brkelkar/common_utils/batch"
	cr "github.com/brkelkar/common_utils/configreader"
	db "github.com/brkelkar/common_utils/databases"
)

//OutstandingAttar as model
type OutstandingAttar struct {
	cAttar CommonAttr
}

func (o *OutstandingAttar) initOutstanding() {
	o.cAttar.colMap = make(map[string]int)
	o.cAttar.colName = []string{"CUSTOMERCODE", "DOCUMENTNUMBER", "DOCUMENTDATE", "AMOUNT", "ADJUSTEDAMOUNT", "PENDINGAMOUNT", "DUEDATE"}

	for _, val := range o.cAttar.colName {
		o.cAttar.colMap[val] = -1
	}
}

//OutstandingCloudFunction used to load outstanding file to database
func (o *OutstandingAttar) OutstandingCloudFunction(g *utils.GcsFile, cfg cr.Config) (err error) {
	log.Printf("Starting outstanding file upload for :%v/%v ", g.FilePath, g.FileName)

	o.initOutstanding()
	reader := csv.NewReader(g.GcsClient.GetReader())
	reader.Comma = '|'
	flag := 1
	var Outstanding []models.Outstanding

	for {
		fileRow, err := reader.Read()

		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(err)
		}
		var tempOutstanding models.Outstanding

		for i, val := range fileRow {
			if flag == 1 {
				o.cAttar.colMap[strings.ToUpper(val)] = i
			} else {
				switch i {
				case -1:
					break
				case o.cAttar.colMap["CUSTOMERCODE"]:
					tempOutstanding.CustomerCode = val
				case o.cAttar.colMap["DOCUMENTNUMBER"]:
					tempOutstanding.DocumentNumber = val
				case o.cAttar.colMap["DOCUMENTDATE"]:
					tempOutstanding.DocumentDate, _ = utils.ConvertDate(val)
				case o.cAttar.colMap["AMOUNT"]:
					tempOutstanding.Amount, _ = strconv.ParseFloat(val, 64)
				case o.cAttar.colMap["ADJUSTEDAMOUNT"]:
					tempOutstanding.AdjustedAmount, _ = strconv.ParseFloat(val, 64)
				case o.cAttar.colMap["PENDINGAMOUNT"]:
					tempOutstanding.PendingAmount, _ = strconv.ParseFloat(val, 64)
				case o.cAttar.colMap["DUEDATE"]:
					tempOutstanding.DueDate, _ = utils.ConvertDate(val)
				}

			}
		}
		if flag == 0 {
			Outstanding = append(Outstanding, tempOutstanding)
		}
		flag = 0
	}
	recordCount := len(Outstanding)
	if recordCount > 0 {

		jsonValue, _ := json.Marshal(Outstanding)
		resp, err := http.Post("http://"+cfg.Server.Host+":"+strconv.Itoa(cfg.Server.Port)+"/api/outstanding", "application/json", bytes.NewBuffer(jsonValue))
		if err != nil || resp.Status != "200 OK" {
			fmt.Println("Error while calling request", err)

			// If upload service
			var d db.DbObj
			dbPtr, err := d.GetConnection("awacs_smart", cfg)
			if err != nil {
				log.Print(err)
				g.GcsClient.MoveObject(g.FileName, "error_Files/"+g.FileName, "balatestawacs")
				log.Println("Porting Error :" + g.FileName)

				return err
			}

			dbPtr.AutoMigrate(&models.Outstanding{})
			//Insert records to temp table
			totalRecordCount := recordCount
			batchSize := bt.GetBatchSize(Outstanding[0])

			if totalRecordCount <= batchSize {
				dbPtr.Save(Outstanding)
			} else {
				remainingRecords := totalRecordCount
				updateRecordLastIndex := batchSize
				startIndex := 0
				for {
					if remainingRecords < 1 {
						break
					}
					updateStockBatch := Outstanding[startIndex:updateRecordLastIndex]
					dbPtr.Save(updateStockBatch)
					remainingRecords = remainingRecords - batchSize
					startIndex = updateRecordLastIndex
					if remainingRecords < batchSize {
						updateRecordLastIndex = updateRecordLastIndex + remainingRecords
					} else {
						updateRecordLastIndex = updateRecordLastIndex + batchSize
					}
				}
			}
		}
		// If either of the loading is successful move file to ported
		g.GcsClient.MoveObject(g.FileName, "ported/"+g.FileName, "balatestawacs")
		log.Println("Porting Done :" + g.FileName)
	}
	return
}
