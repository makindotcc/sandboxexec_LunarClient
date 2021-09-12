package main

import (
	"errors"
	"log"

	"github.com/makindotcc/lunarlauncher"
)

func downloadTextures(texturesBaseUrl string, textures []lunarlauncher.TextureMeta) {
	for processedCount, texture := range textures {
		log.Printf("Downloading texture (%d/%d): %s\n", (processedCount + 1),
			len(textures), texture)
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

func downloadLaunchArtifacts(mcVersion lunarlauncher.McVersion,
	artifacts []lunarlauncher.Artifact) {
	for processedCount, artifact := range artifacts {
		log.Printf("Downloading artifact (%d/%d): %s\n", (processedCount + 1),
			len(artifacts), artifact)
		lunarlauncher.DownloadArtifact(mcVersion, artifact)
	}
}

func main() {
	const mcVersion = lunarlauncher.Mc1_16

	log.Println("Preparing lunar assets...")
	log.Println("Fetching launch meta...")
	launchMeta, err := lunarlauncher.FetchLaunchMeta(mcVersion, "master")
	if err != nil {
		panic(err)
	}
	log.Println("Got launch meta:", launchMeta)

	downloadLaunchTextures(launchMeta)

	downloadLaunchArtifacts(mcVersion, launchMeta.LaunchTypeData.Artifacts)
}
