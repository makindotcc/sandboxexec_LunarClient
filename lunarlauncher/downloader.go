package lunarlauncher

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
)

var ErrAlreadyDownloaded = errors.New("already downloaded")

func isAlreadyDownloaded(downloadPath string) (bool, error) {
	_, err := os.Stat(downloadPath)
	if err == nil {
		// todo: compare hash
		return true, nil
	}
	if !errors.Is(err, os.ErrNotExist) {
		return false, fmt.Errorf("file stat (%s): %w", downloadPath, err)
	}
	return false, nil
}

func downloadFile(filePathRelative string, fileUrl string) error {
	fileUrlParsed, err := url.Parse(fileUrl)
	if err != nil {
		return fmt.Errorf("parse file url: %w", err)
	}
	// prevent downloading file from our mac, for example by using "file" schema
	// like http.Get for sure won't handle it anyway, but maybe someday
	// go will get a handler for "file" schema so what then???
	// we don't wanna put our private files to game "asset" dir that
	// sandbox has access to. 
	if fileUrlParsed.Scheme != "https" {
		return fmt.Errorf("invalid url scheme (%s)", fileUrlParsed.Scheme)
	}

	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get wd: %w", err)
	}

	downloadPath := path.Join(wd, filePathRelative)
	isSubDir, err := IsSubDir(wd, downloadPath)
	if err != nil {
		return fmt.Errorf("sub dir validate: %w", err)
	}
	if !isSubDir {
		return fmt.Errorf("directory path traversal '%s'", downloadPath)
	}
	
	alreadyDownloaded, err := isAlreadyDownloaded(downloadPath)
	if err != nil {
		return fmt.Errorf("is already downloaded: %w", err)
	}
	if alreadyDownloaded {
		return ErrAlreadyDownloaded
	}

	response, err := http.Get(fileUrl)
	if err != nil {
		return nil
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid response status code (%d)", response.StatusCode)
	}

	err = os.MkdirAll(filepath.Dir(downloadPath), os.ModePerm)
	if err != nil {
		return fmt.Errorf("mkdir (%s)", downloadPath)
	}
	file, err := os.Create(downloadPath)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return fmt.Errorf("write to file from response: %w", err)
	}
	return nil
}
