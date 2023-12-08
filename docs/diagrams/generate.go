package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func convertToOutputPath(inputPath string) string {
	base := filepath.Base(inputPath)
	filename := strings.TrimSuffix(base, filepath.Ext(base))
	return fmt.Sprintf("%s.svg", filename)
}

func getInputPaths(infs fs.FS, recurse bool) []string {
	var inputPaths []string
	err := fs.WalkDir(infs, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filepath.Ext(path) == ".d2" {
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
func generateDiagram(inputPath, outputPath string) {
	cmd := exec.Command("d2", inputPath, outputPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Printf("Failed to generate diagram: %v\n", err)
	}
}

func main() {
	_, err := exec.LookPath("d2")
	if err != nil {
		log.Fatal("Need to have d2lang (available here: https://d2lang.com/) installed to generate the diagrams.")
	}

	_, filename, _, _ := runtime.Caller(0)
	diagramsDir := filepath.Dir(filename)
	defaultInputPathDir := filepath.Join(diagramsDir, "input")
	defaultOutputPathDir := filepath.Join(diagramsDir, "output")
	recursiveDefault := false

	var inputPathDir string
	var outputPathDir string
	var recursive bool
	flag.StringVar(&inputPathDir, "i", defaultInputPathDir, "Path containing .d2 files to be converted to .svg")
	flag.StringVar(&outputPathDir, "o", defaultOutputPathDir, "Path in which .svg files will be put")
	flag.BoolVar(&recursive, "r", recursiveDefault, "Should recursively walk subdirs")
	flag.Parse()
	if !flag.Parsed() {
		log.Fatal("Failed to parse flags")
	}

	// If we passed in input or output flags, make them relative to the directory from which the go run command was run.
	cwd, err := os.Getwd()
	if err != nil {
		log.Println("Error while getting current working directory:", err)
	}
	if isFlagPassed("i") {
		inputPathDir = filepath.Join(cwd, inputPathDir)
	}
	if isFlagPassed("o") {
		outputPathDir = filepath.Join(cwd, outputPathDir)
	}

	inputPaths := getInputPaths(os.DirFS(inputPathDir), recursive)
	for _, inputPath := range inputPaths {
		fullInputPath := filepath.Join(inputPathDir, inputPath)
		fullOutputPath := filepath.Join(outputPathDir, convertToOutputPath(inputPath))
		generateDiagram(fullInputPath, fullOutputPath)
	}
}
