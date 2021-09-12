package main

import (
	"os"
	"testing"
)

func TestInvalidIsSubDir(t *testing.T) {
	invalidPaths := []string{
		"../../essa/../essa",
		"../essa",
		"..",
		"/Users",
		"/essa",
		"/root",
	}

	workingDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	for _, path := range invalidPaths {
		subDir, err := IsSubDir(workingDir, path)
		if err != nil {
			panic(err)
		}
		if subDir {
			t.Errorf("IsSubDir returned true for path: %s\n", path)
		} else {
			t.Logf("Validated path '%s' successfully\n", path)
		}
	}
}

func TestValidIsSubDir(t *testing.T) {
	validPaths := []string{
		"",
		"essa/../das/fdg/ads/fgd",
		"textures",
		"textures/asset/../skin/../cape",
	}

	workingDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	for _, path := range validPaths {
		subDir, err := IsSubDir(workingDir, path)
		if err != nil {
			panic(err)
		}
		if subDir {
			t.Logf("Validated path '%s' successfully\n", path)
		} else {
			t.Errorf("IsSubDir returned false for path: %s\n", path)
		}
	}
}
