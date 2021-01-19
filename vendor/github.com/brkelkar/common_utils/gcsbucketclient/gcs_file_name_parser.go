package gcsbucketclient

import (
	"errors"
	"strings"

	"github.com/brkelkar/common_utils/logger"
)

//GetBucketAndFileName Parse GCS file to seperate bucket name and file name
func GetBucketAndFileName(filePath string) (bucketName string, fileName string, err error) {
	if len(filePath) == 0 {
		err = errors.New("Filepath incorrent")
		logger.Error("File Path is empty", err)
		return "", "", err
	}
	if !strings.Contains(filePath, "gs://") {

		err = errors.New("Filepath is not full path")
		logger.Error("Filepath is not full path", err)
		return "", "", err
	}
	tempString := strings.Split(filePath, "/")
	bucketName = tempString[2]
	fileName = tempString[3]
	return
}
