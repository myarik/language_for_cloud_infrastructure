package main

import (
	"bufio"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const DefaultNumConsumer = 3

func init() {
	log.SetOutput(os.Stdout)
	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
}

//contentProducer downloads a content and sends to a channel.
func contentProducer(sourceFile string) <-chan []byte {
	bodyCh := make(chan []byte)
	go func() {
		// Open a source file
		file, err := os.Open(sourceFile)
		if err != nil {
			log.Fatal(err)
		}
		defer func() { _ = file.Close() }()
		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			sourceUrl := scanner.Text()
			client := &http.Client{Timeout: time.Second * 5}
			resp, err := client.Get(sourceUrl)
			if err != nil {
				log.WithError(err).Error("cannot connect to a host")
				return
			}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.WithError(err).Error("cannot read a content body")
				return
			}
			log.WithField("url", sourceUrl).Debug("Downloaded")
			bodyCh <- body
			_ = resp.Body.Close()
		}
		close(bodyCh)
	}()
	return bodyCh

}

//contentConsumer saves a content to a file
func contentConsumer(done chan<- struct{}, bodyCh <-chan []byte, tmpDir string, workerID int) {
	for data := range bodyCh {
		// Add timeout to see how it works
		if log.IsLevelEnabled(log.DebugLevel) {
			time.Sleep(3 * time.Second)
		}

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
		log.Debugf("[WORKER %d]Finished writing %s", workerID, contentFile.Name())
	}
	done <- struct{}{}
}

// This is the simple web scraper.
// The scraper gets data from sources and saves them to our local machine
// This scraper has only three simultaneous consumers to prevent the storage overload
func main() {
	start := time.Now()

	sourceFile := os.Getenv("CONTENT_FILE")
	if sourceFile == "" {
		log.Fatal("cannot find the CONTENT_FILE environment variable")
	}

	isDebuglevel := os.Getenv("DEBUG")
	if isDebuglevel == "True" || isDebuglevel == "true" {
		log.SetLevel(log.DebugLevel)
	}

	tmpDir, err := ioutil.TempDir("", "demo")
	if err != nil {
		log.WithError(err).Error("cannot create tmp directory")
		return
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Set up a done channel
	done := make(chan struct{})
	defer close(done)

	//Stat producer
	bodyCh := contentProducer(sourceFile)
	log.Debug("Start producer")

	// Start consumers
	for i := 0; i < DefaultNumConsumer; i++ {
		go func(workerID int) {
			contentConsumer(done, bodyCh, tmpDir, workerID)
		}(i)
	}
	log.Debug("Start consumers")

	// Wait until all producers finished work
	for i := 0; i < DefaultNumConsumer; i++ {
		<-done
	}
	log.Infof("Execution time: %s seconds.", time.Since(start))
}
