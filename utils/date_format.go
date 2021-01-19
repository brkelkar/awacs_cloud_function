package utils

import "time"

//DateFormatMap provides mapper for date format
var (
	DateFormatMap map[string]string
	err           error
)

func init() {
	DateFormatMap = make(map[string]string)
	DateFormatMap["d-M-yyyy"] = "2-1-2006"
	DateFormatMap["d/M/yyyy"] = "2/1/2006"
	DateFormatMap["dd-MM-yyyy"] = "02-01-2006"
	DateFormatMap["dd-MM-yyyy h:mm:ss"] = "02-01-2006 3:04:05"
	DateFormatMap["dd-MM-yyyy h:mm:ss tt"] = "02-01-2006 3:04:05 AM"
	DateFormatMap["dd-MM-yyyy HH:mm"] = "02-01-2006 15:04"
	DateFormatMap["dd-MM-yyyy HH:mm:ss"] = "02-01-2006 15:04:05"
	DateFormatMap["dd-MM-yyyy hh:mm:ss tt"] = "02-01-2006 15:04:05 AM"
	DateFormatMap["dd-MMM-yy"] = "02-Jan-06"
	DateFormatMap["dd-MMM-yy HH:mm:ss"] = "02-Jan-06 15:04:05"
	DateFormatMap["dd-MMM-yy hh:mm:ss tt"] = "02-Jan-06 15:04:05 AM"
	DateFormatMap["dd-MMM-yyyy"] = "02-Jan-2006"
	DateFormatMap["dd-MMM-yyyy HH:mm"] = "02-Jan-2006 15:04"
	DateFormatMap["dd-MMM-yyyy HH:mm:ss"] = "02-Jan-2006 15:04:05"
	DateFormatMap["dd-MMM-yyyy hh:mm:ss tt"] = "02-Jan-2006 15:04:05 AM"
	DateFormatMap["dd/MM/yyyy"] = "02/01/2006"
	DateFormatMap["dd/MM/yyyy h:mm:ss"] = "02/01/2006 3:04:05"
	DateFormatMap["dd/MM/yyyy h:mm:ss tt"] = "02/01/2006 3:04:05 AM"
	DateFormatMap["dd/MM/yyyy HH:mm"] = "02/01/2006 15:04"
	DateFormatMap["dd/MM/yyyy HH:mm:ss"] = "02/01/2006 15:04:05"
	DateFormatMap["dd/MM/yyyy hh:mm:ss tt"] = "02/01/2006 15:04:05 AM"
	DateFormatMap["dd/MMM/yy"] = "02/Jan/06"
	DateFormatMap["dd/MMM/yy HH:mm:ss"] = "02/Jan/06 15:04:05"
	DateFormatMap["dd/MMM/yy hh:mm:ss tt"] = "02/Jan/06 15:04:05 AM"
	DateFormatMap["dd/MMM/yyyy"] = "02/Jan/2006"
	DateFormatMap["dd/MMM/yyyy HH:mm"] = "02/Jan/2006 15:04"
	DateFormatMap["dd/MMM/yyyy HH:mm:ss"] = "02/Jan/2006 15:04:05"
	DateFormatMap["dd/MMM/yyyy hh:mm:ss tt"] = "02/Jan/2006 15:04:05 AM"
	DateFormatMap["MM-yyyy"] = "01-2006"
	DateFormatMap["MM/yyyy"] = "01/2006"
	DateFormatMap["MMM-yy"] = "Jan-06"
	DateFormatMap["MMM-yyyy"] = "Jan-2006"
	DateFormatMap["MMM/yy"] = "Jan/06"
	DateFormatMap["MMM/yyyy"] = "Jan/2006"
	DateFormatMap["yyyy MM dd"] = "2006 01 02"
	DateFormatMap["yyyy-MM-dd"] = "2006-01-02"
	DateFormatMap["yyyy/MM/dd"] = "2006/01/02"
	DateFormatMap["MM-yy"] = "01-06"
	DateFormatMap["MM/yy"] = "01/06"
	DateFormatMap["dd MM yyyy HH:mm:ss"] = "02 01 2006 15:04:05"
	DateFormatMap["dd-MM-yy"] = "02-01-06"
	DateFormatMap["dd/MM/yy"] = "02/01/06"
	DateFormatMap["d-M-yyyy hh:mm:ss tt"] = "2-1-2006 15:04:05 AM"
	DateFormatMap["d/M/yyyy hh:mm:ss tt"] = "2/1/2006 15:04:05 AM"
	DateFormatMap["dd-MM-yy tt hh:mm:ss"] = "02-01-06 AM 15:04:05"
	DateFormatMap["dd MMM yy hh:mm:ss tt"] = "02 Jan 06 15:04:05 AM"
	DateFormatMap["dd-MM-yy HH:mm:ss"] = "02-01-06 15:04:05"
	DateFormatMap["dd/MM/yy HH:mm:ss"] = "02/01/06 15:04:05"

}

//ConvertDate takes string returns time.time pointer
func ConvertDate(dateString string) (*time.Time, error) {
	for _, val := range DateFormatMap {
		t, err := time.Parse(val, dateString)
		if err != nil {
			continue
		} else {
			return &t, nil
		}

	}
	return nil, err

}
