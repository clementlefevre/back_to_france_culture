package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

var ROOT_FOLDER = "./podcasts"
var ROOT_DEVICE = "e:/Music"

/* // https://www.franceculture.fr/emissions/la-compagnie-des-auteurs/homere-14-qui-est-homere
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
} */

func main() {
	m_local := RetrieveMp3FilesPaths(ROOT_FOLDER)
	m_device := RetrieveMp3FilesPaths(ROOT_DEVICE)
	m_diff := getDiffBetweenLocalAndDevice(m_local, m_device)
	SyncFolders(m_diff)
}

func RetrieveMp3FilesPaths(path string) map[string][]string {

	m := make(map[string][]string)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	for _, dir := range files {
		if dir.Mode().IsDir() {

			files, err := ioutil.ReadDir(filepath.Join(path, dir.Name()))
			if err != nil {
				log.Fatal(err)
			}
			for _, f := range files {
				if filepath.Ext(f.Name()) == ".mp3" {
					//full_path, err := filepath.Abs(filepath.Join(ROOT_FOLDER, dir.Name(), f.Name()))
					if err != nil {
						log.Fatal(err)

					}
					m[dir.Name()] = append(m[dir.Name()], f.Name())
				}
			}
		}
	}
	return m
}

func getDiffBetweenLocalAndDevice(m_local map[string][]string, m_device map[string][]string) map[string][]string {
	m_diff := make(map[string][]string) // initialize some storage for the diff

	for dir := range m_local {
		if _, ok := m_device[dir]; !ok { // check if the key from the first map exists in the second
			m_diff[dir] = m_local[dir]
		} else {
			m_diff[dir] = difference(m_local[dir], m_device[dir])
		}

	}
	return m_diff
}

func difference(slice1 []string, slice2 []string) []string {
	var diff []string

	for _, s1 := range slice1 {
		found := false
		for _, s2 := range slice2 {
			if s1 == s2 {
				found = true
				break
			}
		}
		// String not found. We add it to return slice
		if !found {
			diff = append(diff, s1)
		}
	}
	// Swap the slices, only if it was the first loop

	return diff
}

func Copy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_EXCL, os.ModeDir|0666)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

func SyncFolders(m_diff map[string][]string) {
	for dir, files := range m_diff {
		for _, file := range files {
			src := filepath.Join(ROOT_FOLDER, dir, file)
			dst := filepath.Join(ROOT_DEVICE, dir, file)
			log.Printf("Copying %s\n into %s\n", src, dst)
			Copy(src, dst)
		}

	}
}
