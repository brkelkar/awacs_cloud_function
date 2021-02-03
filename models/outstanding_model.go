package models

import (
	"time"
)

//Outstanding Table model
type Outstanding struct {
	UserId         string  `gorm:"column:userId"`
	CustomerCode   string  `gorm:"column:CustomerCode"`
	DocumentNumber string  `gorm:"column:DocumentNumber"`
	DocumentDate   string  `gorm:"type:datetime;column:DocumentDate"`
	Amount         float64 `gorm:"column:Amount"`
	AdjustedAmount float64 `gorm:"column:AdjustedAmount"`
	PendingAmount  float64 `gorm:"column:PendingAmount"`
	DueDate        string  `gorm:"type:datetime;column:DueDate"`
}

//Outstanding Table model
type CustomerOutstanding struct {
	UserId          string     `gorm:"column:UserId"`
	CustomerCode    string     `gorm:"column:CustomerCode"`
	Outstanding     float64    `gorm:"column:Outstanding"`
	OutstandingJson string     `gorm:"column:OutstandingJson"`
	LastUpdated     time.Time `gorm:"column:LastUpdated"`
}

func (o *CustomerOutstanding) getTableName() (tableName string) {
	tableNameprefix := "TMP_SMART_CUSTOMEROUTSTANDING_"
	t := time.Now().UTC()
	timeformat := t.Format("200601021504")
	tableName = tableNameprefix + timeformat[0:len(timeformat)-1]
	return
}

//TableName retruns temp table name for Outstanding table
func (o CustomerOutstanding) TableName() string {
	return "dbo." + o.getTableName()
}
