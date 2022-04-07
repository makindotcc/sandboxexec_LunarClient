package lunarlauncher

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/artdarek/go-unzip"
)

type ArtifactType string

const (
	ClassPath ArtifactType = "CLASS_PATH"
	Natives   ArtifactType = "NATIVES"
)

type Artifact struct {
	Name string `json:"name"`
	Sha1 string `json:"sha1"`
	Url  string `json:"url"`
	Type ArtifactType `json:"type"`
}

func isArtifactZip(artifact Artifact) bool {
	return filepath.Ext(artifact.Name) == ".zip"
}

func unzipArtifact(artifact Artifact, artifactPath string) error {
	var targetDirName string
	if strings.HasPrefix(artifact.Name, "natives") {
		targetDirName = "natives"
	} else {
		targetDirName = artifact.Name
	}
	targetFullDirPath := path.Join(filepath.Dir(artifactPath), targetDirName)
	uz := unzip.New(artifactPath, targetFullDirPath)
	return uz.Extract()
}

func DownloadArtifact(mcVersion McVersion, artifact Artifact) error {
	artifactPath := path.Join("offline", string(mcVersion), artifact.Name)
	err := downloadFileToWD(artifactPath, artifact.Url)
	if err != nil {
		return fmt.Errorf("download file: %w", err)
	}
	if isArtifactZip(artifact) {
		unzipArtifact(artifact, artifactPath)
	}
	return err
}

func DownloadLog4j() error {
	// for whatever reason lc downloads log4j 2.17.1 and tries to put this 
	// library in minecraft's home directory.
	// we can't allow that because vanilla minecraft could read it,
	// and we don't sandbox vanilla minecraft, so it would be a sandbox escape.
	// So we need to download log4j from a trusted source.
	//
	// Client log:
	// Failed to write jar bytes: /Users/makin/Library/Application Support/minecraft/libraries/org/apache/logging/log4j/log4j-api/2.17.1: Operation not permitted
	// Failed to download /Users/makin/Library/Application Support/minecraft/libraries/org/apache/logging/log4j/log4j-api/2.17.1/log4j-api-2.17.1.jar!
	// Failed to write jar bytes: /Users/makin/Library/Application Support/minecraft/libraries/org/apache/logging/log4j/log4j-core/2.17.1: Operation not permitted
	// Failed to download /Users/makin/Library/Application Support/minecraft/libraries/org/apache/logging/log4j/log4j-core/2.17.1/log4j-core-2.17.1.jar!

	artifacts := []struct{
		path string
		url string
	}{
		{"org/apache/logging/log4j/log4j-api/2.17.1/log4j-api-2.17.1.jar","https://repo1.maven.org/maven2/org/apache/logging/log4j/log4j-core/2.17.1/log4j-core-2.17.1.jar"},
		{"org/apache/logging/log4j/log4j-core/2.17.1/log4j-core-2.17.1.jar","https://repo1.maven.org/maven2/org/apache/logging/log4j/log4j-api/2.17.1/log4j-api-2.17.1.jar"},
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("get user home dir: %s", err)
	}
	librariesPath := path.Join(home, "Library", "Application Support", "minecraft", "libraries")
	for _, artifact := range artifacts {
		path := path.Join(librariesPath, artifact.path)
		log.Printf("Downloading %s from %s to %s\n", artifact.path, artifact.url, path)
		err := downloadFile(path, artifact.url)
		if err != nil && !errors.Is(err, ErrAlreadyDownloaded) {
			return fmt.Errorf("download artifact %s: %s", artifact.url, err)
		}
	}
	return nil
}
