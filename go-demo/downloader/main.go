package main

import (
	"bufio"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetOutput(os.Stdout)
	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
}

//downloadContent downloads a content
func downloadContent(url string) (body []byte, err error) {
	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Get(url)
	if err != nil {
		return body, err
	}
	defer func() { _ = resp.Body.Close() }()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return body, err
	}
	return body, nil
}

//saveContent saves a content to a file
func saveContent(data []byte, tmpDir string) error {
	contentFile, err := ioutil.TempFile(tmpDir, "go_*.mov")
	if err != nil {
		return errors.Wrap(err, "cannot create a tmp file")
	}
	defer func() { _ = contentFile.Close() }()

	_, err = contentFile.Write(data)
	if err != nil {
		return errors.Wrap(err, "cannot write to a file")
	}
	return nil
}

// This is the simple web scraper.
// The scraper gets data from sources and saves them to our local machine
// for further analysis
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

	file, err := os.Open(sourceFile)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = file.Close() }()
	scanner := bufio.NewScanner(file)

	tmpDir, err := ioutil.TempDir("", "demo")
	if err != nil {
		log.WithError(err).Fatal("cannot create tmp directory")
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	wg := sync.WaitGroup{}

	for scanner.Scan() {
		wg.Add(1)
		// Running concurrently
		go func(fileUrl string) {
			defer wg.Done()
			log.Debugf("Begin downloading %s", fileUrl)

			content, err := downloadContent(fileUrl)
			if err != nil {
				log.WithError(err).Error("cannot download a content")
				return
			}

			if err := saveContent(content, tmpDir); err != nil {
				log.WithError(err).Error("cannot save a content")
				return
			}
			log.Debugf("Finished writing %s", fileUrl)
		}(scanner.Text())
	}
	// waiting until all tasks are completed
	wg.Wait()

	if err := scanner.Err(); err != nil {
		log.WithError(err).Fatal("cannot read source file")
	}

	duration := time.Since(start)
	log.Infof("Execution time: %s seconds.", duration)
}
