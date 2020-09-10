package main

import (
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"sync"
	"time"
)

var FILES = []string{
	"0cf50f1c99234954b00340471538ce9d.MOV",
	"0db9a58b669048dc999eb8f11f7ba424.MOV",
	"0d38ceda70b14ccfaf6960514615757f.MOV",
	"0CB55372-0173-49F7-9EAF-6CF1A40382C5.MOV",
	"0f132134b2474cbd858559ed979835a3.MOV",
}

func init() {
	log.SetOutput(os.Stdout)
	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
}

//downloadContent downloads a content and sends to a channel.
//The function only accepts a channel for sending values.
func downloadContent(bodyCh chan<- []byte, url string) {
	client := &http.Client{Timeout: time.Second * 5}
	resp, err := client.Get(url)
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
	bodyCh <- body
}

//saveContent saves a content to a file
func saveContent(bodyCh <-chan []byte, tmpDir string, workerID int) {
	for data := range bodyCh {
		contentFile, err := ioutil.TempFile(tmpDir, "go_*.mov")
		if err != nil {
			log.WithError(err).Error("cannot create a tmp file")
			continue
		}
		_, err = contentFile.Write(data)
		if err != nil {
			log.WithError(err).Error("cannot write to a file")
			continue
		}
		if err := contentFile.Close(); err != nil {
			log.WithError(err).Error("cannot close a file")
		}
		log.Infof("[WORKER %d]Finished writing %s", workerID, contentFile.Name())
	}
}

// This is the simple web scraper.
// The scraper gets data from sources and saves them to our local machine
// This scraper has only three simultaneous consumers to prevent the storage overload
func main() {
	hostURL := os.Getenv("API_HOST_URL")
	if hostURL == "" {
		log.Error("cannot find the API_HOST_URL environment variable")
		return
	}

	const NumReceivers = 3

	start := time.Now()

	tmpDir, err := ioutil.TempDir("", "demo")
	if err != nil {
		log.WithError(err).Error("cannot create tmp directory")
		return
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	bodyCh := make(chan []byte)

	wgReceivers := sync.WaitGroup{}
	wgReceivers.Add(NumReceivers)

	wg := sync.WaitGroup{}
	wg.Add(len(FILES))
	// the sender
	for _, file := range FILES {
		go func(fileName string) {
			defer wg.Done()
			log.Infof("Begin downloading %s", fileName)

			u, err := url.Parse(hostURL)
			if err != nil {
				log.WithError(err).Error("cannot create an url path")
			}
			u.Path = path.Join(u.Path, fileName)
			downloadContent(bodyCh, u.String())
		}(file)
	}

	// receivers
	for i := 0; i < NumReceivers; i++ {
		go func(workerID int) {
			defer wgReceivers.Done()
			// Receive values until bodyCh is
			// closed and the value buffer queue
			// of bodyCh becomes empty.
			saveContent(bodyCh, tmpDir, workerID)
		}(i)
	}

	// waiting until senders finished a work
	wg.Wait()
	// close the data channel
	close(bodyCh)

	// waiting until receivers process a data
	wgReceivers.Wait()

	duration := time.Since(start)
	log.Infof("Execution time: %s seconds.", duration)
}
