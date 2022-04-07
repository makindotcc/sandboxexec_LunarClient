package lunarlauncher

import (
	"bufio"
	"errors"
	"fmt"
	"net/http"
	"path"
	"strings"
)

const texturesDir = "textures"

type TextureMeta struct {
	Path string
	Hash string
}

func DownloadTexture(texturesBaseUrl string, meta TextureMeta) error {
	return downloadFileToWD(path.Join(texturesDir, meta.Path), texturesBaseUrl+meta.Hash)
}

func parseRawTextureMeta(rawMeta string) (TextureMeta, error) {
	splitted := strings.Split(rawMeta, " ")
	if len(splitted) != 2 {
		return TextureMeta{}, errors.New("invalid raw meta format")
	}
	return TextureMeta{
		Path: splitted[0],
		Hash: splitted[1],
	}, nil
}

// Fetch used textures meta list.
func FetchTextures(fileIndexUrl string) ([]TextureMeta, error) {
	r, err := http.Get(fileIndexUrl)
	if err != nil {
		return nil, fmt.Errorf("http get (url=%s): %w", fileIndexUrl, err)
	}
	if r.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid status code (url=%s) code='%d': %w",
				fileIndexUrl, r.StatusCode, err)
	}
	defer r.Body.Close()

	textures := make([]TextureMeta, 0, 4096)
	scanner := bufio.NewScanner(r.Body)
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
