package models

import (
	"fmt"
	"time"
)

//Stocks Table model
type Stocks struct {
	UserId      string  `gorm:"column:UserId"`
	ProductCode string  `gorm:"column:ProductCode"`
	Closing     float64 `gorm:"column:Closing"`
}

func getTableName() (tableName string) {
	tableNamePreFix := "TMP_SMART_STOCKS_"
	t := time.Now()
	tableName = fmt.Sprintf("%s%d%02d%02d%d%d%s", tableNamePreFix, t.Year(), t.Month(), t.Day(), t.Hour(), (t.Minute() / 5), "0")
	return
}

//TableName retruns temp table name for Stock table
func (Stocks) TableName() string {
	return "dbo." + getTableName()
}
