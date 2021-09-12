package main

import (
	"errors"
	"log"

	"github.com/makindotcc/lunarlauncher"
)

func downloadTextures(texturesBaseUrl string, textures []lunarlauncher.TextureMeta) {
	allTexturesCount := len(textures)
	for processedCount, texture := range textures {
		log.Printf("Downloading texture (%d/%d): %s\n", (processedCount + 1), allTexturesCount, texture)
		err := lunarlauncher.DownloadTexture(texturesBaseUrl, texture)
		if errors.Is(err, lunarlauncher.ErrAlreadyDownloaded) {
			log.Printf("Texture %v is already downloaded. Skipping...\n", texture)
		} else if err != nil {
			log.Printf("Texture %v download error: %s\n", texture, err)
		}
	}
}

func downloadLaunchTextures(launchMeta lunarlauncher.LaunchMeta) {
	textures, err := lunarlauncher.FetchTextures(launchMeta.Textures.IndexURL)
	if err != nil {
		panic(err)
	}
	downloadTextures(launchMeta.Textures.BaseURL, textures)
}

func main() {
	log.Println("Preparing lunar assets...")
	log.Println("Fetching launch meta...")
	launchMeta, err := lunarlauncher.FetchLaunchMeta(lunarlauncher.Mc1_16, "master")
	if err != nil {
		panic(err)
	}
	log.Println("Got launch meta:", launchMeta)
	
	downloadLaunchTextures(launchMeta)
}
