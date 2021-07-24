package main

import (
	"backtofranceculture/downloader/services"
	"flag"
	"strings"
)

var ROOT_FOLDER = "./podcasts"
var ROOT_DEVICE = "e:/Music"

// https://www.franceculture.fr/emissions/la-compagnie-des-auteurs/homere-14-qui-est-homere
func main() {

	url := flag.String("url", "", "urls separated by coma")
	file_with_urls := flag.String("file", "", "file with urls separated by coma")
	flag.Parse()

	services.CreateFolder(ROOT_FOLDER)

	allUrls := strings.Split(*url, ",")

	if len(*file_with_urls) > 0 {
		allUrls = nil
		records := services.ReadCsvFile(*file_with_urls)

		for _, record := range records {
			allUrls = append(allUrls, record[0])
		}
	}

	services.DownloadAndSaveToFileAll(allUrls, ROOT_FOLDER)

	m_local := services.RetrieveMp3FilesPaths(ROOT_FOLDER)
	m_device := services.RetrieveMp3FilesPaths(ROOT_DEVICE)
	m_diff := services.GetDiffBetweenLocalAndDevice(m_local, m_device)
	services.SyncFolders(m_diff, ROOT_FOLDER, ROOT_DEVICE)
}
