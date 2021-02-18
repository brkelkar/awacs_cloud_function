package functions

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"time"

	"awacs.com/awcacs_cloud_function/models"
	"awacs.com/awcacs_cloud_function/utils"

	//bt "github.com/brkelkar/common_utils/batch"
	cr "github.com/brkelkar/common_utils/configreader"
	//db "github.com/brkelkar/common_utils/databases"
)

//OutstandingAttar as model
type OutstandingAttar struct {
	cAttar CommonAttr
}

func (o *OutstandingAttar) initOutstanding(cfg cr.Config) {
	o.cAttar.colMap = make(map[string]int)
	o.cAttar.colName = []string{"CUSTOMERCODE", "DOCUMENTNUMBER", "DOCUMENTDATE", "AMOUNT", "ADJUSTEDAMOUNT", "PENDINGAMOUNT", "DUEDATE"}

	for _, val := range o.cAttar.colName {
		o.cAttar.colMap[val] = -1
	}

	apiPath = "/api/outstanding"
	URLPath = utils.GetHostURL(cfg) + apiPath
}

//OutstandingCloudFunction used to load outstanding file to database
func (o *OutstandingAttar) OutstandingCloudFunction(g *utils.GcsFile, cfg cr.Config) (err error) {
	log.Printf("Starting outstanding file upload for :%v/%v ", g.FilePath, g.FileName)
	o.initOutstanding(cfg)
	g.FileType = "O"
	// reader := csv.NewReader(g.GcsClient.GetReader())
	// reader.Comma = '|'
	var reader *bufio.Reader
	reader = bufio.NewReader(g.GcsClient.GetReader())
	flag := 1
	var Outstanding []models.Outstanding

	for {
		//fileRow, err := reader.Read()
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			g.ErrorMsg = "Error while reading file"
			g.LogFileDetails(false)
			return err
		}
		var tempOutstanding models.Outstanding
		line = strings.TrimSpace(line)
		lineSlice := strings.Split(line, "|")
		for i, val := range lineSlice {
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
					tempOutstanding.DocumentDate = val
				case o.cAttar.colMap["AMOUNT"]:
					tempOutstanding.Amount, _ = strconv.ParseFloat(val, 64)
				case o.cAttar.colMap["ADJUSTEDAMOUNT"]:
					tempOutstanding.AdjustedAmount, _ = strconv.ParseFloat(val, 64)
				case o.cAttar.colMap["PENDINGAMOUNT"]:
					tempOutstanding.PendingAmount, _ = strconv.ParseFloat(val, 64)
				case o.cAttar.colMap["DUEDATE"]:
					tempOutstanding.DueDate = val
				}
			}
		}
		tempOutstanding.UserId = g.DistributorCode
		if flag == 0 {
			Outstanding = append(Outstanding, tempOutstanding)
		}
		flag = 0

		if err == io.EOF {
			break
		}
	}

	outstandingMap := make(map[string]models.CustomerOutstanding)
	for _, val := range Outstanding {
		key := val.UserId + val.CustomerCode
		if _, ok := outstandingMap[key]; !ok {
			var tout models.CustomerOutstanding
			tout.OutstandingJson = GetJsonstring(val) //`"[{"CustomerCode":"` + val.CustomerCode + `","DocumentNumber":"` + val.DocumentNumber + `","DocumentDate":"` + val.DocumentDate + `","Amount":"` + fmt.Sprintf("%f", val.Amount) + `","PendingAmount":"` + fmt.Sprintf("%f", val.PendingAmount) + `","AdjustedAmount":"` + fmt.Sprintf("%f", val.AdjustedAmount) + `","DueDate":"` + val.DueDate + `"}"`
			tout.Outstanding = val.PendingAmount
			tout.UserId = val.UserId
			tout.CustomerCode = val.CustomerCode
			tout.LastUpdated = time.Now()

			outstandingMap[key] = tout
		} else {
			t, _ := outstandingMap[key]
			t.OutstandingJson = t.OutstandingJson + "," + GetJsonstring(val) //`",{"CustomerCode":"` + val.CustomerCode + `","DocumentNumber":"` + val.DocumentNumber + `","DocumentDate":"` + val.DocumentDate + `","Amount":"` + fmt.Sprintf("%f", val.Amount) + `","PendingAmount":"` + fmt.Sprintf("%f", val.PendingAmount) + `","AdjustedAmount":"` + fmt.Sprintf("%f", val.AdjustedAmount) + `","DueDate":"` + val.DueDate + `"}]"`
			t.Outstanding = t.Outstanding + val.PendingAmount
			t.UserId = val.UserId
			t.CustomerCode = val.CustomerCode
			t.LastUpdated = time.Now()

			outstandingMap[key] = t
		}
	}

	var customerOutstanding []models.CustomerOutstanding
	for _, val := range outstandingMap {
		val.OutstandingJson = "[" + val.OutstandingJson + "]"
		customerOutstanding = append(customerOutstanding, val)
	}

	recordCount := len(customerOutstanding)
	jsonValue, _ := json.Marshal(customerOutstanding)
	if recordCount > 0 {
		err = utils.WriteToSyncService(URLPath, jsonValue)
		if err != nil {
			// If upload service
			// var d db.DbObj
			// dbPtr, err := d.GetConnection("smartdb", cfg)
			if err != nil {
				log.Print(err)
				g.GcsClient.MoveObject(g.FileName, "error_Files/"+g.FileName, "balatestawacs")
				log.Println("Porting Error :" + g.FileName)
				g.ErrorMsg = "Error while connecting to db"
				g.LogFileDetails(false)
				return err
			}

			// dbPtr.AutoMigrate(&models.CustomerOutstanding{})
			// //Insert records to temp table
			// totalRecordCount := recordCount
			// batchSize := bt.GetBatchSize(customerOutstanding[0])

			// if totalRecordCount <= batchSize {
			// 	err = dbPtr.Save(customerOutstanding).Error
			// 	if err != nil {
			// 		g.ErrorMsg = "Error while writing records to db"
			// 		g.LogFileDetails(false)
			// 		return err
			// 	}
			// } else {
			// 	remainingRecords := totalRecordCount
			// 	updateRecordLastIndex := batchSize
			// 	startIndex := 0
			// 	for {
			// 		if remainingRecords < 1 {
			// 			break
			// 		}
			// 		updateStockBatch := customerOutstanding[startIndex:updateRecordLastIndex]
			// 		err = dbPtr.Save(updateStockBatch).Error
			// 		if err != nil {
			// 			g.ErrorMsg = "Error while writing records to db"
			// 			g.LogFileDetails(false)
			// 			return err
			// 		}
			// 		remainingRecords = remainingRecords - batchSize
			// 		startIndex = updateRecordLastIndex
			// 		if remainingRecords < batchSize {
			// 			updateRecordLastIndex = updateRecordLastIndex + remainingRecords
			// 		} else {
			// 			updateRecordLastIndex = updateRecordLastIndex + batchSize
			// 		}
			// 	}
			// }
		}
	}

	// If either of the loading is successful move file to ported
	g.GcsClient.MoveObject(g.FileName, "ported/"+g.FileName, "balatestawacs")
	log.Println("Porting Done :" + g.FileName)
	g.Records = recordCount
	g.LogFileDetails(true)
	return
}

//GetJsonstring concat json string
func GetJsonstring(outstanding models.Outstanding) (jsonString string) {
	jsonString = `{"CustomerCode":"` + outstanding.CustomerCode + `","DocumentNumber":"` +
		outstanding.DocumentNumber + `","DocumentDate":"` + outstanding.DocumentDate + `","Amount":"` +
		fmt.Sprintf("%f", outstanding.Amount) + `","PendingAmount":"` + fmt.Sprintf("%f", outstanding.PendingAmount) +
		`","AdjustedAmount":"` + fmt.Sprintf("%f", outstanding.AdjustedAmount) + `","DueDate":"` + outstanding.DueDate + `"}`
	return
}
