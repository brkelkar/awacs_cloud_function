package utils

import (
	"bytes"
	"errors"
	"net/http"
)

//WriteToSyncService sends records to syncUpload service
func WriteToSyncService(URLPath string, payload []byte) (err error) {

	resp, err := http.Post(URLPath, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return
	}
	if resp.Status != "200 OK" {
		return errors.New("Failed to write by Sync service")
	}
	return nil
}
