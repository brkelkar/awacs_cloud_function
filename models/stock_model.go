package models

import (
	"strconv"
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
	timeformat := t.Format("200601021504")
	timeFormatLength := 12
	min, _ := strconv.Atoi(timeformat[timeFormatLength-1 : timeFormatLength])
	offset := "0"
	if min > 5 {
		offset = "5"
	}

	tableName = tableNamePreFix + timeformat[0:timeFormatLength-1] + offset
	return
}

//TableName retruns temp table name for Stock table
func (s Stocks) TableName() string {
	return "dbo." + s.getTableName()
}
