package functions

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"awacs.com/awcacs_cloud_function/models"
	"awacs.com/awcacs_cloud_function/utils"
	bt "github.com/brkelkar/common_utils/batch"
	cr "github.com/brkelkar/common_utils/configreader"
	db "github.com/brkelkar/common_utils/databases"
	gc "github.com/brkelkar/common_utils/gcsbucketclient"
)

//InvoiceAttr used to hold Invoice file parsing attributes required
type InvoiceAttr struct {
	developerID             string
	multiLinedistributorMap map[string]bool
	cAttr                   CommonAttr
}

func (i *InvoiceAttr) initInvoice(cfg cr.Config) {
	i.cAttr.colMap = make(map[string]int)
	i.cAttr.colName = []string{"USERID", "BILLNUMBER", "BILLDATE", "CHALLANNUMBER", "CHALLANDATE", "BUYERCODE", "REASON", "COMPANYNAME", "UPC", "PRODUCTCODE", "MRP", "BATCH", "EXPIRY", "QUANTITY", "FREEQUANTITY", "RATE", "AMOUNT", "DISCOUNT", "DISCOUNTAMOUNT", "ADDLSCHEME", "ADDLSCHEMEAMOUNT", "ADDLDISCOUNT", "ADDLDISCOUNTAMOUNT", "DEDUCTABLEBEFOREDISCOUNT", "MRPINCLUSIVETAX", "VATAPPLICATION", "VAT", "ADDLTAX", "CST", "SGST", "CGST", "IGST", "BASESCHEMEQUANTITY", "BASESCHEMEFREEQUANTITY", "NETINVOICEAMOUNT", "PAYMENTDUEDATE", "REMARKS", "SAVEDATE", "SYNDATE", "SYNCDATE", "PRODUCTNAME", "PRODUCTPACK", "EMONTH", "EXPMONTH", "CESS", "CESSAMOUNT", "SGSTAMOUNT", "CGSTAMOUNT", "IGSTAMOUNT", "TAXABLEAMOUNT", "HSN", "BARCODE", "ORDERNUMBER", "ORDERDATE", "LASTTRANSACTIONDATE"}
	for _, val := range i.cAttr.colName {
		i.cAttr.colMap[val] = -1
	}
	i.multiLinedistributorMap = getDistributorForMultiLineFile(cfg)
	apiPath = "/api/invoices"
	URLPath = utils.GetHostURL(cfg) + apiPath

}

//InvoiceCloudFunction used to load invoice file to database
func (i *InvoiceAttr) InvoiceCloudFunction(g *utils.GcsFile, cfg cr.Config) (err error) {
	log.Printf("Starting Invoice file upload for :%v/%v ", g.FilePath, g.FileName)
	i.initInvoice(cfg)
	g.FileType = "I"
	fileSplitSlice := strings.Split(g.FileName, "_")
	spiltLen := len(fileSplitSlice)

	// Check if file is in correct format or not
	if !(spiltLen == 7 || spiltLen == 6) {

		g.ErrorMsg = "Invalid file name"
		g.LogFileDetails(false)
		return errors.New("Invalid file name")
	}
	if spiltLen == 6 {
		i.developerID = fileSplitSlice[5]
	}

	replace, replaceDistributorCode := getReplaceStrings(g.DistributorCode)

	flag := 1

	var Invoice []models.Invoice
	var reader *bufio.Reader
	//Distributor with vendor code 113 sents files with multiline records
	//This code will handle these type of files by replaceing \n \r with
	// "" and then identify new line by distributor code
	if _, ok := i.multiLinedistributorMap[g.DistributorCode]; ok {

		data, _ := ioutil.ReadAll(g.GcsClient.GetReader())
		content := string(data)
		content = strings.ReplaceAll(content, "\n", "")
		content = strings.ReplaceAll(content, "\r", "")
		content = strings.ReplaceAll(content, g.DistributorCode, "\n"+g.DistributorCode)
		rc := strings.NewReader(content)
		reader = bufio.NewReader(rc)
	} else {
		reader = bufio.NewReader(g.GcsClient.GetReader())
	}

	// Start reading file line by line
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			g.ErrorMsg = "Error while reading file"
			g.LogFileDetails(false)
			return err
		}
		if _, ok := replaceDistributorCode[g.DistributorCode]; ok {
			for _, replaceVal := range replace {
				// Replace values we get from table (Company name which has '|' as charator in name)
				line = strings.ReplaceAll(line, replaceVal.Search_String, replaceVal.Replace_String)

			}
		}
		line = strings.TrimSpace(line)
		lineSlice := strings.Split(line, "|")
		var tempInvoice models.Invoice
		for index, val := range lineSlice {
			if flag == 1 {

				i.cAttr.colMap[strings.ToUpper(val)] = index

			} else {

				switch index {
				case -1:
					break
				case i.cAttr.colMap["BILLNUMBER"]:
					tempInvoice.BillNumber = val
				case i.cAttr.colMap["BILLDATE"]:
					tempInvoice.BillDate, _ = utils.ConvertDate(val)
				case i.cAttr.colMap["CHALLANNUMBER"]:
					tempInvoice.ChallanNumber = val
				case i.cAttr.colMap["CHALLANDATE"]:
					tempInvoice.ChallanDate, _ = utils.ConvertDate(val)
				case i.cAttr.colMap["BUYERCODE"]:
					tempInvoice.BuyerId = val
				case i.cAttr.colMap["REASON"]:
					tempInvoice.Reason = val
				case i.cAttr.colMap["UPC"]:
					tempInvoice.UPC = val
				case i.cAttr.colMap["PRODUCTCODE"]:
					tempInvoice.SupplierProductCode = val
				case i.cAttr.colMap["PRODUCTNAME"]:
					tempInvoice.SupplierProductName = val
				case i.cAttr.colMap["PRODUCTPACK"]:
					tempInvoice.SupplierProductPack = val
				case i.cAttr.colMap["MRP"]:
					tempInvoice.MRP, _ = strconv.ParseFloat(val, 64)
				case i.cAttr.colMap["BATCH"]:
					tempInvoice.Batch = val
				case i.cAttr.colMap["EXPIRY"]:
					tempInvoice.Expiry, _ = utils.ConvertDate(val)
				case i.cAttr.colMap["QUANTITY"]:
					tempInvoice.Quantity, _ = strconv.ParseFloat(val, 64)
				case i.cAttr.colMap["FREEQUANTITY"]:
					tempInvoice.FreeQuantity, _ = strconv.ParseFloat(val, 64)
				case i.cAttr.colMap["RATE"]:
					tempInvoice.Rate, _ = strconv.ParseFloat(val, 64)
				case i.cAttr.colMap["AMOUNT"]:
					tempInvoice.Amount, _ = strconv.ParseFloat(val, 64)
				case i.cAttr.colMap["DISCOUNT"]:
					tempInvoice.Discount, _ = strconv.ParseFloat(val, 64)
				case i.cAttr.colMap["DISCOUNTAMOUNT"]:
					tempInvoice.DiscountAmount, _ = strconv.ParseFloat(val, 64)
				case i.cAttr.colMap["ADDLSCHEME"]:
					tempInvoice.AddlScheme, _ = strconv.ParseFloat(val, 64)
				case i.cAttr.colMap["ADDLSCHEMEAMOUNT"]:
					tempInvoice.AddlSchemeAmount, _ = strconv.ParseFloat(val, 64)
				case i.cAttr.colMap["ADDLDISCOUNT"]:
					tempInvoice.AddlDiscount, _ = strconv.ParseFloat(val, 64)
				case i.cAttr.colMap["ADDLDISCOUNTAMOUNT"]:
					tempInvoice.AddlDiscountAmount, _ = strconv.ParseFloat(val, 64)
				case i.cAttr.colMap["DEDUCTABLEBEFOREDISCOUNT"]:
					tempInvoice.DeductableBeforeDiscount, _ = strconv.ParseFloat(val, 64)
				case i.cAttr.colMap["MRPINCLUSIVETAX"]:
					tempInvoice.MRPInclusiveTax, _ = strconv.Atoi(val)
				case i.cAttr.colMap["VATAPPLICATION"]:
					tempInvoice.VATApplication = val
				case i.cAttr.colMap["VAT"]:
					tempInvoice.VAT, _ = strconv.ParseFloat(val, 64)
				case i.cAttr.colMap["ADDLTAX"]:
					tempInvoice.AddlTax, _ = strconv.ParseFloat(val, 64)
				case i.cAttr.colMap["CST"]:
					tempInvoice.CST, _ = strconv.ParseFloat(val, 64)
				case i.cAttr.colMap["SGST"]:
					tempInvoice.SGST, _ = strconv.ParseFloat(val, 64)
				case i.cAttr.colMap["CGST"]:
					tempInvoice.CGST, _ = strconv.ParseFloat(val, 64)
				case i.cAttr.colMap["IGST"]:
					tempInvoice.IGST, _ = strconv.ParseFloat(val, 64)
				case i.cAttr.colMap["BASESCHEMEQUANTITY"]:
					tempInvoice.BaseSchemeQuantity, _ = strconv.ParseFloat(val, 64)
				case i.cAttr.colMap["BASESCHEMEFREEQUANTITY"]:
					tempInvoice.BaseSchemeFreeQuantity, _ = strconv.ParseFloat(val, 64)
				case i.cAttr.colMap["PAYMENTDUEDATE"]:
					tempInvoice.PaymentDueDate, _ = utils.ConvertDate(val)
				case i.cAttr.colMap["REMARKS"]:
					tempInvoice.Remarks = val
				case i.cAttr.colMap["COMPANYNAME"]:
					tempInvoice.CompanyName = val
				case i.cAttr.colMap["NETINVOICEAMOUNT"]:
					tempInvoice.NetInvoiceAmount, _ = strconv.ParseFloat(val, 64)
				case i.cAttr.colMap["LASTTRANSACTIONDATE"]:
					tempInvoice.LastTransactionDate, _ = utils.ConvertDate(val)
				case i.cAttr.colMap["SGSTAMOUNT"]:
					tempInvoice.SGSTAmount, _ = strconv.ParseFloat(val, 64)
				case i.cAttr.colMap["CGSTAMOUNT"]:
					tempInvoice.CGSTAmount, _ = strconv.ParseFloat(val, 64)
				case i.cAttr.colMap["IGSTAMOUNT"]:
					tempInvoice.IGSTAmount, _ = strconv.ParseFloat(val, 64)
				case i.cAttr.colMap["CESS"]:
					tempInvoice.Cess, _ = strconv.ParseFloat(val, 64)
				case i.cAttr.colMap["CESSAMOUNT"]:
					tempInvoice.CessAmount, _ = strconv.ParseFloat(val, 64)
				case i.cAttr.colMap["TAXABLEAMOUNT"]:
					tempInvoice.TaxableAmount, _ = strconv.ParseFloat(val, 64)
				case i.cAttr.colMap["HSN"]:
					tempInvoice.HSN = val
				case i.cAttr.colMap["ORDERNUMBER"]:
					tempInvoice.OrderNumber = val
				case i.cAttr.colMap["ORDERDATE"]:
					tempInvoice.OrderDate, _ = utils.ConvertDate(val)
				case i.cAttr.colMap["BARCODE"]:
					tempInvoice.Barcode = val

				}

			}

		}
		if flag == 0 {
			tempInvoice.DeveloperId = i.developerID
			tempInvoice.File_Path = g.FilePath
			tempInvoice.File_Received_Dttm, _ = utils.ConvertDate(g.ProcessingTime)
			tempInvoice.SupplierId = g.DistributorCode
			tempInvoice.Inv_File_Id = -1
			Invoice = append(Invoice, tempInvoice)
		}
		flag = 0
	}

	//Got final record to write
	recordCount := len(Invoice)
	if recordCount < 0 {

		jsonValue, _ := json.Marshal(Invoice)
		err := utils.WriteToSyncService(URLPath, jsonValue)
		if err != nil {
			//Try to write directly to db
			var d db.DbObj
			dbPtr, err := d.GetConnection("awacs_smart", cfg)
			if err != nil {
				log.Print(err)
				g.GcsClient.MoveObject(g.FileName, "error_Files/"+g.FileName, "balatestawacs")
				log.Println("Porting Error :" + g.FileName)
				g.ErrorMsg = "Error while connecting to db"
				g.LogFileDetails(false)
				return err
			}

			dbPtr.AutoMigrate(&models.Invoice{})

			//Insert records to temp table
			totalRecordCount := recordCount
			batchSize := bt.GetBatchSize(Invoice[0])
			if totalRecordCount <= batchSize {
				err = dbPtr.Save(Invoice).Error
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
					updateStockBatch := Invoice[startIndex:updateRecordLastIndex]

					err = dbPtr.Save(updateStockBatch).Error
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
	g.GcsClient.MoveObject(g.FileName, "ported/"+g.FileName, "balatestawacs")
	log.Println("Porting Done :" + g.FileName)
	g.Records = recordCount
	g.LogFileDetails(true)
	return nil

}

func getReplaceStrings(distributorCode string) (replace []models.ReplaceStrings, replaceDistributorCode map[string]bool) {
	bucketName, fileName, err := gc.GetBucketAndFileName("gs://awacs_config/Search_replace_config.csv")
	if err != nil {
		log.Print("Error while parsing filepath :gs://awacs_config/Search_replace_config.csv ")
	}
	var gcsObj gc.GcsBucketClient
	gClient := gcsObj.InitClient().SetBucketName(bucketName).SetNewReader(fileName)
	if !gClient.GetLastStatus() {
		log.Print("Error while reading filepath :gs://awacs_config/Search_replace_config.csv ")
	}

	rc := csv.NewReader(gClient.GetReader())
	rc.Comma = ','
	for {
		line, err := rc.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Print(err)
		}
		var tempReplace models.ReplaceStrings
		fmt.Println(line)

		tempReplace.DistributorCode, tempReplace.Search_String, tempReplace.Replace_String = line[0], line[1], line[2]
		fmt.Printf("Distributor code %v", tempReplace.DistributorCode)
		fmt.Println(replaceDistributorCode)
		replaceDistributorCode = make(map[string]bool)
		replaceDistributorCode[tempReplace.DistributorCode] = true
		replace = append(replace, tempReplace)
	}
	return
}

func getDistributorForMultiLineFile(cfg cr.Config) (distributorDetail map[string]bool) {
	requestURL := "http://" + cfg.Server.Host + ":" + strconv.Itoa(cfg.Server.Port) + "/api/distributors"
	res, err := http.Get(requestURL)
	if err != nil {
		log.Println("Error while getting distributor details " + err.Error())
		return
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Println("Error while parsing request body for distributor responce " + err.Error())
		return
	}
	var distributors []models.Distributors
	json.Unmarshal(body, &distributors)
	for _, val := range distributors {
		distributorDetail[val.User_StockistCode_cd] = true
	}
	return
}
