package utils

import (
	"bytes"
	"errors"
	"log"
	"net/http"
)

//WriteToSyncService sends records to syncUpload service
func WriteToSyncService(URLPath string, payload []byte) (err error) {

	resp, err := http.Post(URLPath, "application/json", bytes.NewBuffer(payload))
	log.Println(resp)
	if err != nil {
		log.Println(err)
		return
	}
	if resp.Status != "200 OK" {
		return errors.New("Failed to write by Sync service status := " + resp.Status)
	}
	return nil
}
