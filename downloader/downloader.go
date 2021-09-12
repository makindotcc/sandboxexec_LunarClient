package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const texturesDir = "textures"
const texturesBaseUrl = "https://textures.lunarclientcdn.com/file/"

var errAlreadyDownloaded = errors.New("already downloaded")

type textureMeta struct {
	path string
	hash string
}

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

func downloadTexture(rootDir string, meta textureMeta) error {
	downloadPath := path.Join(rootDir, meta.path)
	isSubDir, err := IsSubDir(rootDir, downloadPath)
	if err != nil {
		return fmt.Errorf("sub dir validate: %w", err)
	}
	if !isSubDir {
		return fmt.Errorf("texture directory path traversal '%s'", downloadPath)
	}
	
	alreadyDownloaded, err := isAlreadyDownloaded(downloadPath)
	if err != nil {
		return fmt.Errorf("is already downloaded: %w", err)
	}
	if alreadyDownloaded {
		return errAlreadyDownloaded
	}

	response, err := http.Get(texturesBaseUrl + meta.hash)
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

func downloadTextures(rootDir string, textures []textureMeta) {
	allTexturesCount := len(textures)
	for processedCount, texture := range textures {
		log.Printf("Downloading texture (%d/%d): %s\n", (processedCount + 1), allTexturesCount, texture)
		err := downloadTexture(rootDir, texture)
		if errors.Is(err, errAlreadyDownloaded) {
			log.Printf("Texture %v is already downloaded. Skipping...\n", texture)
		} else if err != nil {
			log.Printf("Texture %v download error: %s\n", texture, err)
		}
	}
}

func parseRawTextureMeta(rawMeta string) (textureMeta, error) {
	splitted := strings.Split(rawMeta, " ")
	if len(splitted) != 2 {
		return textureMeta{}, errors.New("invalid raw meta format")
	}
	return textureMeta{
		path: splitted[0],
		hash: splitted[1],
	}, nil
}

func readTextures() ([]textureMeta, error) {
	file, err := os.Open("textures.txt")
	if err != nil {
		return nil, fmt.Errorf("open textures.txt: %w", err)
	}

	textures := make([]textureMeta, 0, 512)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		rawMeta := scanner.Text()
		meta, err := parseRawTextureMeta(rawMeta)
		if err != nil {
			return nil, fmt.Errorf("parse raw meta '%s': %w", rawMeta, err)
		}
		textures = append(textures, meta)
	}
	return textures, nil
}

func main() {
	log.Println("Welcome witam in downloader")
	texturesMetas, err := readTextures()
	if err != nil {
		panic(err)
	}
	downloadTextures(texturesDir, texturesMetas)
}
