package models

//ReplaceStrings structs holds value pair of string to replace and new string
type ReplaceStrings struct {
	DistributorCode string
	Search_String   string `gorm:"column:Search_String"`
	Replace_String  string `gorm:"column:Replace_String"`
}

//TableName return table name for Porting Validation
func (ReplaceStrings) TableName() string {
	return "dbo.SMART_Invoice_Porting_Validation"
}
