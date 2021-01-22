package models

import (
	"time"
)

//Outstanding Table model
type Outstanding struct {
	CustomerCode   string     `gorm:"column:CustomerCode"`
	DocumentNumber string     `gorm:"column:DocumentNumber"`
	DocumentDate   *time.Time `gorm:"type:datetime;column:DocumentDate"`
	Amount         float64    `gorm:"column:Amount"`
	AdjustedAmount float64    `gorm:"column:AdjustedAmount"`
	PendingAmount  float64    `gorm:"column:PendingAmount"`
	DueDate        *time.Time `gorm:"type:datetime;column:DueDate"`
}

func (o *Outstanding) getTableName() (tableName string) {
	tableNameprefix := "TMP_SMART_OUTSTANDING_"
	t := time.Now().UTC()
	timeformat := t.Format("200601021504")
	tableName = tableNameprefix + timeformat[0:len(timeformat)-1]
	return
}

//TableName retruns temp table name for Outstanding table
func (o Outstanding) TableName() string {
	return "dbo." + o.getTableName()
}
