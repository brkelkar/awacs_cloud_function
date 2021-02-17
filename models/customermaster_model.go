package models

import (
	"time"
)

//CustomerMaster Table model
type CustomerMaster struct {
	UserId          string `gorm:"column:UserId"`
	Code            string `gorm:"column:Code"`
	CompanionCode   string `gorm:"column:CompanionCode"`
	Name            string `gorm:"column:Name"`
	Address1        string `gorm:"column:Address1"`
	Address2        string `gorm:"column:Address2"`
	Address3        string `gorm:"column:Address3"`
	City            string `gorm:"column:City"`
	State           string `gorm:"column:State"`
	Area            string `gorm:"column:Area"`
	Pincode         string `gorm:"column:Pincode"`
	KeyPerson       string `gorm:"column:KeyPerson"`
	Cell            string `gorm:"column:Cell"`
	Phone           string `gorm:"column:Phone"`
	Email           string `gorm:"column:Email"`
	DrugLic1        string `gorm:"column:DrugLic1"`
	DrugLic2        string `gorm:"column:DrugLic2"`
	DrugLic3        string `gorm:"column:DrugLic3"`
	DrugLic4        string `gorm:"column:DrugLic4"`
	DrugLic5        string `gorm:"column:DrugLic5"`
	DrugLic6        string `gorm:"column:DrugLic6"`
	GSTIN           string `gorm:"column:GSTIN"`
	PAN             string `gorm:"column:PAN"`
	SalesmanCode    string `gorm:"column:SalesmanCode"`
	IsLocked        string `gorm:"column:IsLocked"`
	IsLockedBilling string `gorm:"column:IsLockedBilling"`
	AllowDelivery   string `gorm:"column:AllowDelivery"`
}

func (c *CustomerMaster) getTableName() (tableName string) {
	tableNameprefix := "TMP_SMART_CUSTOMERMASTER_"
	t := time.Now().UTC()
	timeformat := t.Format("200601021504")
	tableName = tableNameprefix + timeformat[0:len(timeformat)-1]
	return
}

//TableName retruns temp table name for Outstanding table
func (c CustomerMaster) TableName() string {
	return "dbo." + c.getTableName()
}
