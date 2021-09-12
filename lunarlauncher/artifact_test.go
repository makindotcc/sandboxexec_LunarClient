package lunarlauncher

import "testing"

func TestIsArtifactZip(t *testing.T) {
	zipArtifact := Artifact{
		Name: "natives.zip",
	}
	if !isArtifactZip(zipArtifact) {
		t.Errorf("artifact %s is a zip file!", zipArtifact.Name)
	}

	jarArtifact := Artifact{
		Name: "client.jar",
	}
	if isArtifactZip(jarArtifact) {
		t.Errorf("artifact %s is not a zip file!", zipArtifact.Name)
	}
}
