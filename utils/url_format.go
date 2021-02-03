package utils

import (
	cr "github.com/brkelkar/common_utils/configreader"
)

//GetHostURL return Host URL using config
func GetHostURL(cfg cr.Config) (URL string) {
	URL = "http://" + cfg.Server.Host
	return
}
