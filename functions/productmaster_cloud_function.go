package functions

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"strconv"
	"strings"

	"awacs.com/awcacs_cloud_function/models"
	"awacs.com/awcacs_cloud_function/utils"
	bt "github.com/brkelkar/common_utils/batch"
	cr "github.com/brkelkar/common_utils/configreader"
	db "github.com/brkelkar/common_utils/databases"
)

//ProductMasterAttar as model
type ProductMasterAttar struct {
	cAttar CommonAttr
}

func (o *ProductMasterAttar) initProductMaster(cfg cr.Config) {
	o.cAttar.colMap = make(map[string]int)
	o.cAttar.colName = []string{"USERID", "UPC", "PRODUCTCODE", "CODE", "FAVCODE",
		"PRODUCTNAME", "NAME", "BOXPACK", "CASEPACK", "PRODUCTPACK", "PACK", "COMPANYNAME", "COMPANYCODE",
		"COMPANYCODE", "COMPANYNAME", "COMPANY", "DIVISIONCODE", "DIVISION", "DIVISIONNAME",
		"PRODUCTCATEGORY", "CATEGORY", "PTS", "PTR", "MRP", "HSN",
		"CONTENT", "ISACTIVE", "LASTUPDATEDTIME", "CLOSING", "MINDISCOUNT", "MAXDISCOUNT", "ISLOCKED"}

	for _, val := range o.cAttar.colName {
		o.cAttar.colMap[val] = -1
	}

	apiPath = "/api/productmaster"
	URLPath = utils.GetHostURL(cfg) + apiPath
}

//ProductMasterCloudFunction used to load outstanding file to database
func (o *ProductMasterAttar) ProductMasterCloudFunction(g *utils.GcsFile, cfg cr.Config) (err error) {
	log.Printf("Starting product master file upload for :%v/%v ", g.FilePath, g.FileName)
	o.initProductMaster(cfg)
	g.FileType = "P"
	
	var reader *bufio.Reader
	reader = bufio.NewReader(g.GcsClient.GetReader())
	
	flag := 1
	var Productmaster []models.ProductMaster

	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			g.ErrorMsg = "Error while reading file"
			g.LogFileDetails(false)
			return err
		}
		line = strings.TrimSpace(line)
		lineSlice := strings.Split(line, "|")
		var tempProductmaster models.ProductMaster
		for i, val := range lineSlice {
			if flag == 1 {
				o.cAttar.colMap[strings.ToUpper(val)] = i
			} else {
				switch i {
				case -1:
					break
				case o.cAttar.colMap["USERID"]:
					tempProductmaster.UserId = val
				case o.cAttar.colMap["UPC"]:
					tempProductmaster.UPC = val
				case o.cAttar.colMap["CODE"], o.cAttar.colMap["PRODUCTCODE"]:
					tempProductmaster.Code = val
				case o.cAttar.colMap["FAVCODE"]:
					tempProductmaster.FavCode = val
				case o.cAttar.colMap["NAME"], o.cAttar.colMap["PRODUCTNAME"]:
					tempProductmaster.Name = val
				case o.cAttar.colMap["BOXPACK"]:
					tempProductmaster.BoxPack = val
				case o.cAttar.colMap["CASEPACK"]:
					tempProductmaster.CasePack = val
				case o.cAttar.colMap["PACK"], o.cAttar.colMap["PRODUCTPACK"]:
					tempProductmaster.Pack = val
				case o.cAttar.colMap["COMPANYCODE"]:
					tempProductmaster.CompanyCode = val
				case o.cAttar.colMap["COMPANY"], o.cAttar.colMap["COMPANYNAME"]:
					tempProductmaster.CompanyName = val
				case o.cAttar.colMap["DIVISIONCODE"]:
					tempProductmaster.DivisionCode = val
				case o.cAttar.colMap["DIVISION"], o.cAttar.colMap["DIVISIONNAME"]:
					tempProductmaster.DivisionName = val
				case o.cAttar.colMap["CATEGORY"], o.cAttar.colMap["PRODUCTCATEGORY"]:
					tempProductmaster.Category = val
				case o.cAttar.colMap["PTS"]:
					tempProductmaster.PTS, _ = strconv.ParseFloat(val, 64)
				case o.cAttar.colMap["PTR"]:
					tempProductmaster.PTR, _ = strconv.ParseFloat(val, 64)
				case o.cAttar.colMap["MRP"]:
					tempProductmaster.MRP, _ = strconv.ParseFloat(val, 64)
				case o.cAttar.colMap["HSN"]:
					tempProductmaster.HSN = val
				case o.cAttar.colMap["CONTENT"]:
					tempProductmaster.Content = val
				case o.cAttar.colMap["ISACTIVE"]:
					tempProductmaster.IsActive, _ = strconv.ParseBool(val)
				case o.cAttar.colMap["LASTUPDATEDTIME"]:
					tempProductmaster.LastUpdatedTime, _ = utils.ConvertDate(val)
				case o.cAttar.colMap["CLOSING"]:
					tempProductmaster.Closing, _ = strconv.ParseFloat(val, 64)
				case o.cAttar.colMap["MINDISCOUNT"]:
					tempProductmaster.MinDiscount, _ = strconv.ParseFloat(val, 64)
				case o.cAttar.colMap["MAXDISCOUNT"]:
					tempProductmaster.MaxDiscount, _ = strconv.ParseFloat(val, 64)
				case o.cAttar.colMap["ISLOCKED"]:
					tempProductmaster.IsLocked, _ = strconv.ParseBool(val)
				}
			}
		}
		tempProductmaster.UserId = g.DistributorCode
		if flag == 0 {
			Productmaster = append(Productmaster, tempProductmaster)
		}
		flag = 0
	}
	recordCount := len(Productmaster)
	jsonValue, _ := json.Marshal(Productmaster)
	if recordCount > 0 {
		err = utils.WriteToSyncService(URLPath, jsonValue)
		if err != nil {
			// If upload service
			var d db.DbObj
			dbPtr, err := d.GetConnection("smartdb", cfg)
			if err != nil {
				log.Print(err)
				g.GcsClient.MoveObject(g.FileName, "error_Files/"+g.FileName, "balatestawacs")
				log.Println("Porting Error :" + g.FileName)
				g.ErrorMsg = "Error while connecting to db"
				g.LogFileDetails(false)
				return err
			}

			dbPtr.AutoMigrate(&models.ProductMaster{})
			//Insert records to temp table
			totalRecordCount := recordCount
			batchSize := bt.GetBatchSize(Productmaster[0])

			if totalRecordCount <= batchSize {
				err = dbPtr.Save(Productmaster).Error
				if err != nil {
					g.ErrorMsg = "Error while writing records to db"
					g.LogFileDetails(false)
					return err
				}
			} else {
				remainingRecords := totalRecordCount
				updateRecordLastIndex := batchSize
				startIndex := 0
				for {
					if remainingRecords < 1 {
						break
					}
					updateProductBatch := Productmaster[startIndex:updateRecordLastIndex]
					err = dbPtr.Save(updateProductBatch).Error
					if err != nil {
						g.ErrorMsg = "Error while writing records to db"
						g.LogFileDetails(false)
						return err
					}
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
	}

	// If either of the loading is successful move file to ported
	g.GcsClient.MoveObject(g.FileName, "ported/productmaster/"+g.FileName, "balatestawacs")
	log.Println("Porting Done :" + g.FileName)
	g.Records = recordCount
	g.LogFileDetails(true)
	return
}
