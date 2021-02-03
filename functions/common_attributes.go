package functions

var (
	apiPath string
	URLPath string
)

//CommonAttr used for parsing files and holding header index and list
type CommonAttr struct {
	colMap  map[string]int
	colName []string
}
