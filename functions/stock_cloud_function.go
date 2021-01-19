package functions

import (
	"awacs.com/awcacs_cloud_function/models"
	"awacs.com/awcacs_cloud_function/utils"

	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	bt "github.com/brkelkar/common_utils/batch"
	db "github.com/brkelkar/common_utils/databases"

	cr "github.com/brkelkar/common_utils/configreader"
)

//StockAttr used for update Stock file in database
type StockAttr struct {
	cAttr CommonAttr
}

func (s *StockAttr) stockInit() {
	s.cAttr.colMap = make(map[string]int)
	s.cAttr.colName = []string{"USERID", "PRODUCTCODE", "CLOSING"}

	for _, val := range s.cAttr.colName {
		s.cAttr.colMap[val] = -1
	}
}

//StockCloudFunction used to load stock file to database
func (s *StockAttr) StockCloudFunction(g *utils.GcsFile, cfg cr.Config) (err error) {
	log.Printf("Starting Invoice file upload for :%v/%v ", g.FilePath, g.FileName)

	s.stockInit()
	reader := csv.NewReader(g.GcsClient.GetReader())
	reader.Comma = '|'
	flag := 1
	var stock []models.Stocks
	productMap := make(map[string]models.Stocks)

	for {
		fileRow, err := reader.Read()

		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(err)
		}
		var tempStock models.Stocks
		var strproductCode string

		for i, val := range fileRow {
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
	}

	for _, val := range productMap {
		stock = append(stock, val)
	}
	jsonValue, _ := json.Marshal(stock)
	uploadServerURL := "http://" + cfg.Server.Host + ":" + strconv.Itoa(cfg.Server.Port) + "/api/stocks"
	resp, err := http.Post(uploadServerURL, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Println("Error while calling request", err)
	}
	if resp.Status != "200 OK" {
		var d db.DbObj
		dbPtr, err := d.GetConnection("smartdb", cfg)
		if err != nil {
			log.Print(err)
			g.GcsClient.MoveObject(g.FileName, "error_Files/"+g.FileName, "balaawacstest")
			log.Println("Porting Error :" + g.FileName)

			return err
		}
		dbPtr.AutoMigrate(&models.Stocks{})

		totalRecordCount := len(stock)
		batchSize := bt.GetBatchSize(stock[0])
		if totalRecordCount <= batchSize {
			dbPtr.Save(stock)
		} else {
			remainingRecords := totalRecordCount
			updateRecordLastIndex := batchSize
			startIndex := 0
			for {
				if remainingRecords < 1 {
					break
				}
				updateStockBatch := stock[startIndex:updateRecordLastIndex]
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
	return
}
