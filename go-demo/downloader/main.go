package main

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"sync"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
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
	hostURL := os.Getenv("API_HOST_URL")
	if hostURL == "" {
		log.Error("cannot find the API_HOST_URL environment variable")
		return
	}

	start := time.Now()

	tmpDir, err := ioutil.TempDir("", "demo")
	if err != nil {
		log.WithError(err).Error("cannot create tmp directory")
		return
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	wg := sync.WaitGroup{}
	wg.Add(len(FILES))

	for _, file := range FILES {
		go func(fileName string) {
			defer wg.Done()
			log.Debugf("Begin downloading %s", fileName)

			u, err := url.Parse(hostURL)
			if err != nil {
				log.WithError(err).Error("cannot create an url path")
			}
			u.Path = path.Join(u.Path, fileName)
			content, err := downloadContent(u.String())
			if err != nil {
				log.WithError(err).Error("cannot download a content")
				return
			}

			if err := saveContent(content, tmpDir); err != nil {
				log.WithError(err).Error("cannot save a content")
				return
			}
			log.Debugf("Finished writing %s", fileName)
		}(file)
	}
	// waiting until all tasks are completed
	wg.Wait()

	duration := time.Since(start)
	log.Infof("Execution time: %s seconds.", duration)
}
