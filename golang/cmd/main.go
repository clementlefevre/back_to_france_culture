package main

import (
	"backtofranceculture/downloader/services"
	"flag"
	"log"
	"strings"
)

var ROOT_FOLDER = "./podcasts"
var ROOT_DEVICE = "e:/Music"

// https://www.franceculture.fr/emissions/la-compagnie-des-auteurs/homere-14-qui-est-homere
func main() {

	url := flag.String("url", "", "urls separated by coma")
	sync_folder := flag.String("sync", "", "root folder of your external device")

	file_with_urls := flag.String("file", "", "file with urls separated by coma")
	flag.Parse()

	log.Printf("Using %s as destination fodler", *sync_folder)

	services.CreateFolder(ROOT_FOLDER)

	allUrls := strings.Split(*url, ",")

	if len(*file_with_urls) > 0 {
		allUrls = nil
		records := services.ReadCsvFile(*file_with_urls)

		for _, record := range records {
			allUrls = append(allUrls, record[0])
		}
	}

	//services.DownloadAndSaveToFileAll(allUrls, ROOT_FOLDER)

	if len(*sync_folder) > 0 {
		m_local := services.RetrieveMp3FilesPaths(ROOT_FOLDER)
		m_device := services.RetrieveMp3FilesPaths(ROOT_DEVICE)
		m_diff := services.GetDiffBetweenLocalAndDevice(m_local, m_device, *sync_folder)
		services.SyncFolders(m_diff, ROOT_FOLDER, *sync_folder)
	}

}
