package lunarlauncher

import (
	"fmt"
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
	err := downloadFile(artifactPath, artifact.Url)
	if err != nil {
		return fmt.Errorf("download file: %w", err)
	}
	if isArtifactZip(artifact) {
		unzipArtifact(artifact, artifactPath)
	}
	return err
}
