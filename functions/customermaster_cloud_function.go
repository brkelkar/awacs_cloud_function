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

//CustomerMasterAttar as model
type CustomerMasterAttar struct {
	cAttar CommonAttr
}

func (o *CustomerMasterAttar) initCustomerMaster() {
	o.cAttar.colMap = make(map[string]int)
	o.cAttar.colName = []string{"CODE", "COMPANIONCODE", "NAME", "ADDRESS1", "ADDRESS2", "ADDRESS3", "CITY", "STATE", "AREA", "PINCODE", "KEYPERSON", "CELL", "PHONE", "EMAIL", "DRUGLIC1", "DRUGLIC2", "DRUGLIC3", "DRUGLIC4", "DRUGLIC5", "DRUGLIC6", "GSTIN"}

	for _, val := range o.cAttar.colName {
		o.cAttar.colMap[val] = -1
	}
}

//CustomerMasterCloudFunction used to load outstanding file to database
func (o *CustomerMasterAttar) CustomerMasterCloudFunction(g *utils.GcsFile, cfg cr.Config) (err error) {
	log.Printf("Starting customer master file upload for :%v/%v ", g.FilePath, g.FileName)

	o.initCustomerMaster()
	reader := csv.NewReader(g.GcsClient.GetReader())
	reader.Comma = '|'
	flag := 1
	var Customermaster []models.CustomerMaster

	for {
		fileRow, err := reader.Read()

		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(err)
		}
		var tempCustomermaster models.CustomerMaster

		for i, val := range fileRow {
			if flag == 1 {
				o.cAttar.colMap[strings.ToUpper(val)] = i
			} else {
				switch i {
				case -1:
					break
				case o.cAttar.colMap["CODE"]:
					tempCustomermaster.Code = val
				case o.cAttar.colMap["COMPANIONCODE"]:
					tempCustomermaster.CompanionCode = val
				case o.cAttar.colMap["NAME"]:
					tempCustomermaster.Name = val
				case o.cAttar.colMap["ADDRESS1"]:
					tempCustomermaster.Address1 = val
				case o.cAttar.colMap["ADDRESS2"]:
					tempCustomermaster.Address2 = val
				case o.cAttar.colMap["ADDRESS3"]:
					tempCustomermaster.Address3 = val
				case o.cAttar.colMap["CITY"]:
					tempCustomermaster.City = val
				case o.cAttar.colMap["STATE"]:
					tempCustomermaster.State = val
				case o.cAttar.colMap["AREA"]:
					tempCustomermaster.Area = val
				case o.cAttar.colMap["PINCODE"]:
					tempCustomermaster.Pincode = val
				case o.cAttar.colMap["KEYPERSON"]:
					tempCustomermaster.KeyPerson = val
				case o.cAttar.colMap["CELL"]:
					tempCustomermaster.Cell = val
				case o.cAttar.colMap["PHONE"]:
					tempCustomermaster.Phone = val
				case o.cAttar.colMap["EMAIL"]:
					tempCustomermaster.Email = val
				case o.cAttar.colMap["DRUGLIC1"]:
					tempCustomermaster.DrugLic1 = val
				case o.cAttar.colMap["DRUGLIC2"]:
					tempCustomermaster.DrugLic2 = val
				case o.cAttar.colMap["DRUGLIC3"]:
					tempCustomermaster.DrugLic3 = val
				case o.cAttar.colMap["DRUGLIC4"]:
					tempCustomermaster.DrugLic4 = val
				case o.cAttar.colMap["DRUGLIC5"]:
					tempCustomermaster.DrugLic5 = val
				case o.cAttar.colMap["DRUGLIC6"]:
					tempCustomermaster.DrugLic6 = val
				case o.cAttar.colMap["GSTIN"]:
					tempCustomermaster.GSTIN = val
				}
			}
		}
		tempCustomermaster.UserId=g.DistributorCode
		if flag == 0 {
			Customermaster = append(Customermaster, tempCustomermaster)
		}
		flag = 0
	}
	recordCount := len(Customermaster)
	if recordCount > 0 {
		jsonValue, _ := json.Marshal(Customermaster)
		resp, err := http.Post("http://"+cfg.Server.Host+":"+strconv.Itoa(cfg.Server.Port)+"/api/customermaster", "application/json", bytes.NewBuffer(jsonValue))
		if err != nil || resp.Status != "200 OK" {
			fmt.Println("Error while calling request", err)

			// If upload service
			var d db.DbObj
			dbPtr, err := d.GetConnection("smartdb", cfg)
			if err != nil {
				log.Print(err)
				g.GcsClient.MoveObject(g.FileName, "error_Files/"+g.FileName, "balatestawacs")
				log.Println("Porting Error :" + g.FileName)

				return err
			}

			dbPtr.AutoMigrate(&models.CustomerMaster{})
			//Insert records to temp table
			totalRecordCount := recordCount
			batchSize := bt.GetBatchSize(Customermaster[0])

			if totalRecordCount <= batchSize {
				dbPtr.Save(Customermaster)
			} else {
				remainingRecords := totalRecordCount
				updateRecordLastIndex := batchSize
				startIndex := 0
				for {
					if remainingRecords < 1 {
						break
					}
					updateStockBatch := Customermaster[startIndex:updateRecordLastIndex]
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
		g.GcsClient.MoveObject(g.FileName, "ported/customermaster/"+g.FileName, "balatestawacs")
		log.Println("Porting Done :" + g.FileName)
	}
	return
}
