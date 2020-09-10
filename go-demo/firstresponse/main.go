package main

import (
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
)

var StorageMapping = map[string]string{
	"storage0": "0cf50f1c99234954b00340471538ce9d.MOV",
	"storage1": "0db9a58b669048dc999eb8f11f7ba424.MOV",
	"storage2": "0d38ceda70b14ccfaf6960514615757f.MOV",
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

//getHostURL returns the url path
func getHostURL(hostURL, storageID string) (string, error) {
	u, err := url.Parse(hostURL)
	if err != nil {
		return "", errors.Wrap(err, "cannot parse an url path")
	}
	fileName, ok := StorageMapping[storageID]
	if !ok {
		return "", errors.New("cannot find a storage")
	}
	u.Path = path.Join(u.Path, fileName)
	return u.String(), nil
}

//replicaStorage downloads a content and sends to a channel.
//The function only accepts a channel for sending values.
func replicaStorage(respCh chan<- apiResponse, hostURL, storageID string) {
	urlPath, err := getHostURL(hostURL, storageID)
	if err != nil {
		log.WithError(err).Error("cannot get an url")
		return
	}
	//get a content
	client := &http.Client{Timeout: time.Second * 5}
	resp, err := client.Get(urlPath)
	if err != nil {
		log.WithError(err).Error("cannot connect to a host")
		return
	}
	defer func() { _ = resp.Body.Close() }()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithError(err).Error("cannot read a content body")
		return
	}

	respCh <- apiResponse{
		name: storageID,
		body: body,
	}
}

func main() {
	hostURL := os.Getenv("API_HOST_URL")
	if hostURL == "" {
		log.Error("cannot find the API_HOST_URL environment variable")
		return
	}

	start := time.Now()

	// Create a channel on which to send the result.
	respCh := make(chan apiResponse)
	// Send requests to multiple replicas, and use the first response.
	for i := 0; i < 3; i++ {
		go replicaStorage(respCh, hostURL, fmt.Sprintf("storage%d", i))
	}
	resp := <-respCh
	log.Infof("%s returns the first result", strings.Title(resp.name))

	duration := time.Since(start)
	log.Infof("Execution time: %s seconds.", duration)
}
