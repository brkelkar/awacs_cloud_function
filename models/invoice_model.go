package models

import (
	"time"
)

//Invoice model for temp tables
type Invoice struct {
	Id                       int64      `gorm:"primaryKey;column:Id"`
	DeveloperId              string     `gorm:"column:DeveloperId"`
	SupplierId               string     `gorm:"column:SupplierId"`
	BillNumber               string     `gorm:"column:BillNumber"`
	BillDate                 *time.Time `gorm:"type:datetime;column:BillDate"`
	ChallanNumber            string     `gorm:"column:ChallanNumber"`
	ChallanDate              *time.Time `gorm:"type:datetime;column:ChallanDate"`
	BuyerId                  string     `gorm:"column:BuyerId"`
	Reason                   string     `gorm:"column:Reason"`
	BuyerProductCode         string     `gorm:"column:BuyerProductCode"`
	UPC                      string     `gorm:"column:UPC"`
	SupplierProductCode      string     `gorm:"column:SupplierProductCode"`
	SupplierProductName      string     `gorm:"column:SupplierProductName"`
	SupplierProductPack      string     `gorm:"column:SupplierProductPack"`
	MRP                      float64    `gorm:"column:MRP"`
	Batch                    string     `gorm:"column:Batch"`
	Expiry                   *time.Time `gorm:"type:datetime;column:Expiry"`
	Quantity                 float64    `gorm:"column:Quantity"`
	FreeQuantity             float64    `gorm:"column:FreeQuantity"`
	Rate                     float64    `gorm:"column:Rate"`
	Amount                   float64    `gorm:"column:Amount"`
	Discount                 float64    `gorm:"column:Discount"`
	DiscountAmount           float64    `gorm:"column:DiscountAmount"`
	AddlScheme               float64    `gorm:"column:AddlScheme"`
	AddlSchemeAmount         float64    `gorm:"column:AddlSchemeAmount"`
	AddlDiscount             float64    `gorm:"column:AddlDiscount"`
	AddlDiscountAmount       float64    `gorm:"column:AddlDiscountAmount"`
	DeductableBeforeDiscount float64    `gorm:"column:DeductableBeforeDiscount"`
	MRPInclusiveTax          int        `gorm:"column:MRPInclusiveTax"`
	VATApplication           string     `gorm:"column:VATApplication"`
	VAT                      float64    `gorm:"column:VAT"`
	AddlTax                  float64    `gorm:"column:AddlTax"`
	CST                      float64    `gorm:"column:CST"`
	SGST                     float64    `gorm:"column:SGST"`
	CGST                     float64    `gorm:"column:CGST"`
	IGST                     float64    `gorm:"column:IGST"`
	BaseSchemeQuantity       float64    `gorm:"column:BaseSchemeQuantity"`
	BaseSchemeFreeQuantity   float64    `gorm:"column:BaseSchemeFreeQuantity"`
	PaymentDueDate           *time.Time `gorm:"type:datetime;column:PaymentDueDate"`
	Remarks                  string     `gorm:"column:Remarks"`
	CompanyName              string     `gorm:"column:CompanyName"`
	NetInvoiceAmount         float64    `gorm:"column:NetInvoiceAmount"`
	NetPurchaseAmount        float64    `gorm:"column:NetPurchaseAmount"`
	LastTransactionDate      *time.Time `gorm:"type:datetime;column:LastTransactionDate"`
	ReceivedOn               *time.Time `gorm:"type:datetime;column:ReceivedOn"`
	BatchId                  *string    `gorm:"type:uniqueidentifier;column:BatchId"`
	EmailPrepared            bool       `gorm:"column:EmailPrepared"`
	EmailPreparedOn          *time.Time `gorm:"type:datetime;column:EmailPreparedOn"`
	Taken                    bool       `gorm:"column:Taken"`
	TakenOn                  *time.Time `gorm:"type:datetime;column:TakenOn"`
	LastHostAddress          string     `gorm:"column:LastHostAddress"`
	WorkspaceId              int64      `gorm:"column:WorkspaceId"`
	SGSTAmount               float64    `gorm:"column:SGSTAmount"`
	CGSTAmount               float64    `gorm:"column:CGSTAmount"`
	IGSTAmount               float64    `gorm:"column:IGSTAmount"`
	Cess                     float64    `gorm:"column:Cess"`
	CessAmount               float64    `gorm:"column:CessAmount"`
	TaxableAmount            float64    `gorm:"column:TaxableAmount"`
	HSN                      string     `gorm:"column:HSN"`
	OrderNumber              string     `gorm:"column:OrderNumber"`
	OrderDate                *time.Time `gorm:"type:datetime;column:OrderDate"`
	Barcode                  string     `gorm:"column:Barcode"`
	CreatedDate              *time.Time `gorm:"type:datetime;column:CreatedDate"`
	File_Path                string     `gorm:"column:File_Path"`
	File_Received_Dttm       *time.Time `gorm:"type:datetime;column:File_Received_Dttm"`
	Inv_File_Id              float64    `gorm:"column:Inv_File_Id"`
}

func (i *Invoice) getTableName() (tableName string) {
	tableNamePreFix := "TMP_SMART_INVOICE_"
	t := time.Now().UTC()
	timeformat := t.Format("200601021504")
	tableName = tableNamePreFix + timeformat[0:len(timeformat)-1]
	return
}

//TableName retruns temp table name for Invoice details
func (i Invoice) TableName() string {
	return "dbo." + i.getTableName()
}
