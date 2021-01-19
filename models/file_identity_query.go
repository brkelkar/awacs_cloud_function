package models

import (
	"strings"
	"text/template"
)

//FileIndetityQuery used for query formating
type FileIndetityQuery string

func (s FileIndetityQuery) format(data map[string]interface{}) (out string, err error) {
	t := template.Must(template.New("").Parse(string(s)))
	builder := &strings.Builder{}
	if err = t.Execute(builder, data); err != nil {
		return
	}
	out = builder.String()
	return
}

// GetFileIndexQuery Returns get file index query to execute
func (FileIndetityQuery) GetFileIndexQuery(data map[string]interface{}) (out string, err error) {
	var query FileIndetityQuery

	query = `DECLARE @ID NUMERIC(18);
	BEGIN
	IF EXISTS (SELECT TOP 1 * FROM SMART_Invoice_File_Hdr WITH(NoLOCK)
	where Stk_Stockist_code = '{{.DistributorID}}' AND File_Path = '{{.FilePath}}')
	SELECT Inv_File_Id FROM SMART_Invoice_File_Hdr where File_Path = '{{.FilePath}}'
	ELSE
	BEGIN
	INSERT INTO SMART_Invoice_File_Hdr
	(Stk_Stockist_code, File_Path, Table_Name, File_Received_dttm, CreatedBy, Inv_File_Status_Id, Error_Msg)
	VALUES ('{{.DistributorID}}', '{{.FilePath}}', '', '{{.CurrentDate}}', 'Lambda', 2, NULL)
	SELECT @ID = @@IDENTITY INSERT INTO SMART_Invoice_File_Dtl(Inv_File_Id, Inv_File_Status_Id, CreatedBy) values (@ID,2,'Lambda')
	SELECT @ID
	END
	END `

	return query.format(data)

}

// GetUpdateFileIndexQuery Returns Update query
func (FileIndetityQuery) GetUpdateFileIndexQuery(data map[string]interface{}) (out string, err error) {
	var query FileIndetityQuery
	query = `BEGIN
	IF EXISTS (SELECT TOP 1 * FROM SMART_Invoice_File_Hdr WITH(NoLOCK) WHERE Inv_File_Id = {{.FileID}} AND Inv_File_Status_Id = 2)
	UPDATE SMART_Invoice_File_Hdr SET  Inv_File_Status_Id  = 1 , CreatedBy = 'Lambda', File_RecordCnt = {{.RecordCount}}, Table_Name = '{{.TableName}}'
	, Error_Msg = ''  WHERE Inv_File_Id = {{.FileID}}
	INSERT INTO SMART_Invoice_File_Dtl(Inv_File_Id, Inv_File_Status_Id, CreatedBy, Inv_File_Status_Start_dttm, Inv_File_Status_End_dttm)
	values ({{.FileID}},1,'Lambda', GETDATE(), GETDATE())
	END `

	return query.format(data)
}
