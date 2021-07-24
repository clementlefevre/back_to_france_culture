package services

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"sync"

	"github.com/antchfx/htmlquery"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func DownloadAndSaveToFileAll(allUrls []string, ROOT_FOLDER string) {
	// create a channel for work "tasks"
	ch := make(chan string)
	var wg sync.WaitGroup

	// start the workers, we limit the workers to 5
	for t := 1; t < 5; t++ {
		wg.Add(1)
		go GetPodcast(ch, &wg)
	}

	// push the lines to the queue channel for processing
	for _, url := range allUrls {
		ch <- url
	}
	// this will cause the workers to stop and exit their receive loop
	close(ch)
	// make sure they all exit
	wg.Wait()
	log.Printf("Finished to download %d files", len(allUrls))

}

func GetPodcast(ch chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	for url := range ch {
		// do work
		DownloadAndSaveToFile(url, ROOT_FOLDER)
	}

}

func DownloadAndSaveToFile(url string, ROOT_FOLDER string) {

	log.Println("Start Downloading : " + url)

	doc, err := htmlquery.LoadURL(url)
	if err != nil {
		panic(err)
	}
	// Find podcast mp3 url
	list, err := htmlquery.QueryAll(doc, "//button[@class='replay-button paused aod playable blue textualized']")
	if err != nil {
		panic(err)
	}
	result := list[0]
	url_broadcast := htmlquery.SelectAttr(result, "data-emission-title")
	url_title := htmlquery.SelectAttr(result, "data-diffusion-title")
	url_dl := htmlquery.SelectAttr(result, "data-url")

	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	// replace diacritics
	t := transform.Chain(norm.NFD, transform.RemoveFunc(IsMn), norm.NFC)
	broadcast_name, _, _ := transform.String(t, url_broadcast)
	filename, _, _ := transform.String(t, url_title)

	// keep only alphanumerical
	broadcast_name = reg.ReplaceAllString(broadcast_name, "_")
	filename = reg.ReplaceAllString(filename, "_")
	filename = fmt.Sprintf("%s.mp3", filename)

	CreateFolder(fmt.Sprintf("%s/%s", ROOT_FOLDER, broadcast_name))

	fileUrl := url_dl
	err = DownloadFile(fmt.Sprintf("%s/%s/%s", ROOT_FOLDER, broadcast_name, filename), fileUrl)
	if err != nil {
		panic(err)
	}
	log.Println("Downloaded: " + fileUrl)

}

func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
