package paths

import (
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"strings"
)

func GetInputPaths(infs fs.FS, fileExtension string, recurse bool) []string {
	var inputPaths []string
	err := fs.WalkDir(infs, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filepath.Ext(path) == fileExtension {
			inputPaths = append(inputPaths, path)
		}

		// We skip all directories that aren't at the root of the specified FS.
		// In other words, we only want to skip subdirectories.
		if !recurse && d.IsDir() && path != "." {
			return fs.SkipDir
		}
		return nil
	})
	if err != nil {
		log.Println("Error while walking input directory:", err)
	}
	return inputPaths
}

func ReplaceExtension(inputPath string, extension string) string {
	filePath := strings.TrimSuffix(inputPath, filepath.Ext(inputPath))
	return fmt.Sprintf("%s%s", filePath, extension)
}
