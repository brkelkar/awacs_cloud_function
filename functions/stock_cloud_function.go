package functions

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"awacs.com/awcacs_cloud_function/models"
	"awacs.com/awcacs_cloud_function/utils"

	//bt "github.com/brkelkar/common_utils/batch"
	cr "github.com/brkelkar/common_utils/configreader"
	//db "github.com/brkelkar/common_utils/databases"
)

//StockAttr used for update Stock file in database
type StockAttr struct {
	cAttr CommonAttr
}

func (s *StockAttr) stockInit(cfg cr.Config) {
	s.cAttr.colMap = make(map[string]int)
	s.cAttr.colName = []string{"USERID", "PRODUCTCODE", "CLOSING"}

	for _, val := range s.cAttr.colName {
		s.cAttr.colMap[val] = -1
	}
	apiPath = "/api/stocks"
	URLPath = utils.GetHostURL(cfg) + apiPath
}

//StockCloudFunction used to load stock file to database
func (s *StockAttr) StockCloudFunction(g *utils.GcsFile, cfg cr.Config) (err error) {
	log.Printf("Starting stock file upload for :%v/%v ", g.FilePath, g.FileName)
	g.FileType = "S"
	s.stockInit(cfg)
	// reader := csv.NewReader(g.GcsClient.GetReader())
	// reader.Comma = '|'
	var reader *bufio.Reader
	reader = bufio.NewReader(g.GcsClient.GetReader())
	flag := 1
	var stock []models.Stocks
	productMap := make(map[string]models.Stocks)

	for {
		//fileRow, err := reader.Read()
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			fmt.Println(err)
		}
		var tempStock models.Stocks
		var strproductCode string
		line = strings.TrimSpace(line)
		lineSlice := strings.Split(line, "|")
		for i, val := range lineSlice {
			if flag == 1 {
				s.cAttr.colMap[strings.ToUpper(val)] = i
			} else {
				switch i {
				case -1:
					break
				case s.cAttr.colMap["PRODUCTCODE"]:
					strproductCode = val
					tempStock.ProductCode = val
				case s.cAttr.colMap["CLOSING"]:
					if s, err := strconv.ParseFloat(val, 64); err == nil {
						tempStock.Closing = s
					}
				}
				tempStock.UserId = g.DistributorCode
			}
		}

		if flag == 0 {
			val, ok := productMap[strproductCode]
			if ok == true {
				val.Closing = val.Closing + tempStock.Closing
				productMap[strproductCode] = val
			} else {
				productMap[strproductCode] = tempStock
			}
		}
		flag = 0

		if err == io.EOF {
			break
		}
	}

	for _, val := range productMap {
		stock = append(stock, val)
	}
	recordCount := len(stock)
	jsonValue, _ := json.Marshal(stock)
	if recordCount > 0 {
		err = utils.WriteToSyncService(URLPath, jsonValue)
		if err != nil {
			//var d db.DbObj
			//dbPtr, err := d.GetConnection("smartdb", cfg)
			// if err != nil {
			log.Print(err)
			g.GcsClient.MoveObject(g.FileName, "error_Files/"+g.FileName, "balaawacstest")
			log.Println("Porting Error :" + g.FileName)
			g.LogFileDetails(false)
			//return err
			// }
			// dbPtr.AutoMigrate(&models.Stocks{})

			// totalRecordCount := recordCount
			// batchSize := bt.GetBatchSize(stock[0])
			// if totalRecordCount <= batchSize {
			// 	err = dbPtr.Save(stock).Error
			// 	if err != nil {
			// 		g.LogFileDetails(false)
			// 		return err
			// 	}
			// } else {
			// 	// remainingRecords := totalRecordCount
			// 	// updateRecordLastIndex := batchSize
			// 	// startIndex := 0
			// 	// for {
			// 	// 	if remainingRecords < 1 {
			// 	// 		break
			// 	// 	}
			// 	// 	updateStockBatch := stock[startIndex:updateRecordLastIndex]
			// 	// 	err = dbPtr.Save(updateStockBatch).Error
			// 	// 	if err != nil {
			// 	// 		g.LogFileDetails(false)
			// 	// 		return err
			// 	// 	}
			// 	// 	remainingRecords = remainingRecords - batchSize
			// 	// 	startIndex = updateRecordLastIndex
			// 	// 	if remainingRecords < batchSize {
			// 	// 		updateRecordLastIndex = updateRecordLastIndex + remainingRecords
			// 	// 	} else {
			// 	// 		updateRecordLastIndex = updateRecordLastIndex + batchSize
			// 	// 	}
			// 	// }
			// }
		}
	}
	// If either of the loading is successful move file to ported
	g.GcsClient.MoveObject(g.FileName, "ported/"+g.FileName, "balatestawacs")
	g.Records = recordCount
	g.LogFileDetails(true)
	log.Println("Porting Done :" + g.FileName)
	return nil
}
