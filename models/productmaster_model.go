package models

import (
	"time"
)

//ProductMaster Table model
type ProductMaster struct {
	UserId          string     `gorm:"column:UserId"`
	UPC             string     `gorm:"column:UPC"`
	Code            string     `gorm:"column:Code"`
	FavCode         string     `gorm:"column:FavCode"`
	Name            string     `gorm:"column:Name"`
	Pack            string     `gorm:"column:Pack"`
	BoxPack         string     `gorm:"column:BoxPack"`
	CasePack        string     `gorm:"column:CasePack"`
	CompanyCode     string     `gorm:"column:CompanyCode"`
	CompanyName     string     `gorm:"column:CompanyName"`
	DivisionCode    string     `gorm:"column:DivisionCode"`
	DivisionName    string     `gorm:"column:DivisionName"`
	Category        string     `gorm:"column:Category"`
	PTS             float64    `gorm:"column:PTS"`
	PTR             float64    `gorm:"column:PTR"`
	MRP             float64    `gorm:"column:MRP"`
	HSN             string     `gorm:"column:HSN"`
	Content         string     `gorm:"column:Content"`
	IsActive        bool       `gorm:"column:IsActive"`
	LastUpdatedTime *time.Time `gorm:"type:datetime;column:LastUpdatedTime"`
	Closing         float64    `gorm:"column:Closing"`
	MinDiscount     float64    `gorm:"column:MinDiscount"`
	MaxDiscount     float64    `gorm:"column:MaxDiscount"`
	IsLocked        bool       `gorm:"column:IsLocked"`
}

func (p *ProductMaster) getTableName() (tableName string) {
	tableNameprefix := "TMP_SMART_PRODUCTMASTER_"
	t := time.Now().UTC()
	timeformat := t.Format("200601021504")
	tableName = tableNameprefix + timeformat[0:len(timeformat)-1]
	return
}

//TableName retruns temp table name for Outstanding table
func (p ProductMaster) TableName() string {
	return "dbo." + p.getTableName()
}
