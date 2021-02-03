package models

import (
	"time"
)

//Stocks Table model
type Stocks struct {
	UserId      string  `gorm:"column:UserId"`
	ProductCode string  `gorm:"column:ProductCode"`
	Closing     float64 `gorm:"column:Closing"`
}

func (s *Stocks) getTableName() (tableName string) {
	tableNamePreFix := "TMP_SMART_STOCKS_"
	t := time.Now().UTC()
	timeformat:=t.Format("200601021504")
	tableName=tableNamePreFix+timeformat[0:len(timeformat)-1]
	return 
}

//TableName retruns temp table name for Stock table
func (s Stocks) TableName() string {
	return "dbo." + s.getTableName()
}
