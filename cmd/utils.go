package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

type Config struct {
	RepoName      string `json:"repoName"`
	CommiterName  string `json:"commiterName"`
	CommiterEmail string `json:"commiterEmail"`
	GithubApiKey  string `json:"githubApiKey"`
	SaveLocation  string `json:"saveLocation"`
}

func LoadConfig() (Config, error) {
	contents, err := ReadWholeFile("./", "config.json")

	if err != nil {
		return Config{}, err
	}

	var config Config

	err = json.Unmarshal(contents, &config)

	if err != nil {
		return Config{}, err
	}

	return config, nil
}

func IsDir(path string) (bool, error) {
	file, err := os.Open(path)

	if err != nil {
		return false, err
	}

	fileInfo, err := file.Stat()

	if err != nil {
		return false, err
	}

	if fileInfo.IsDir() {
		return true, nil
	}

	return false, nil
}

func ReadWholeFile(filePath string, fileName string) ([]byte, error) {
	contents, err := os.ReadFile(filePath + fileName)

	if err != nil {
		fmt.Printf("Unable to find the file %s in path %s", fileName, filePath)
		return nil, err
	}

	return contents, nil
}

func GetNewestFiles(dir string) ([]string, error) {
	files, _ := os.ReadDir(dir)
	var modTime time.Time
	var names []string
	for _, fi := range files {
		info, _ := fi.Info()
		if info.Mode().IsRegular() {
			if !info.ModTime().Before(modTime) {
				if info.ModTime().After(modTime) {
					modTime = info.ModTime()
					names = names[:0]
				}
				names = append(names, fi.Name())
			}
		}
	}

	return names, nil
}

func Copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
