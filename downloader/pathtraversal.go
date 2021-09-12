package lunarlauncher

import (
	"fmt"
	"path/filepath"
	"strings"
)

func IsSubDir(rootDir string, path string) (bool, error) {
	rootDirAbs, err := filepath.Abs(rootDir)
	if err != nil {
		return false, fmt.Errorf("root dir get path abs: %w", err)
	}
	pathAbs, err := filepath.Abs(path)
	if err != nil {
		return false, fmt.Errorf("sub dir get path abs: %w", err)
	}
	return strings.HasPrefix(pathAbs, rootDirAbs), nil
}
