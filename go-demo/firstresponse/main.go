package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var StorageMapping = map[string]string{
	"storage0": "http://212.183.159.230/10MB.zip",
	"storage1": "http://ipv4.download.thinkbroadband.com/10MB.zip",
	"storage2": "http://speedtest.tele2.net/10MB.zip",
}

//apiResponse represents a response from the worker
type apiResponse struct {
	name string
	body []byte
}

func init() {
	log.SetOutput(os.Stdout)
	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
}

//replicaStorage downloads a content and sends to a channel.
//The function only accepts a channel for sending values.
func replicaStorage(respCh chan<- apiResponse, storageID string) {
	urlPath, ok := StorageMapping[storageID]
	if !ok {
		log.Errorf("cannot get an url: %s", storageID)
		return
	}
	//get a content
	client := &http.Client{Timeout: time.Second * 5}
	resp, err := client.Get(urlPath)
	if err != nil {
		log.WithError(err).Error("cannot connect to a host")
		return
	}
	if resp.StatusCode != http.StatusOK {
		log.WithFields(log.Fields{
			"storage":       storageID,
			"response code": resp.StatusCode,
		}).Errorf("returns error")
		return
	}
	defer func() { _ = resp.Body.Close() }()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithError(err).Error("cannot read a content body")
		return
	}
	log.WithField("storage", storageID).Debug("downloaded file")
	respCh <- apiResponse{
		name: storageID,
		body: body,
	}
}

func main() {
	start := time.Now()

	isDebuglevel := os.Getenv("DEBUG")
	if isDebuglevel == "True" || isDebuglevel == "true" {
		log.SetLevel(log.DebugLevel)
	}

	// Create a channel on which to send the result.
	respCh := make(chan apiResponse)

	// Set up a done channel
	wg := sync.WaitGroup{}
	done := make(chan struct{})
	// Send requests to multiple replicas, and use the first response.
	for i := 0; i < 3; i++ {
		storageID := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			replicaStorage(respCh, fmt.Sprintf("storage%d", storageID))
		}()
	}
	go func() {
		defer close(done)
		// waiting until all tasks are completed
		wg.Wait()
	}()

	for {
		select {
		case <-done:
			log.Error("Replicas doesn't respond")
			return
		case resp := <-respCh:
			log.Infof("%s returns the first result", strings.Title(resp.name))
			duration := time.Since(start)
			log.Infof("Execution time: %s seconds.", duration)
			return
		}
	}
}
