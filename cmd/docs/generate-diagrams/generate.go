package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/volovikariel/IdentityManager/internal/paths"
)

func convertToOutputPath(inputPath string) string {
	base := filepath.Base(inputPath)
	filename := strings.TrimSuffix(base, filepath.Ext(base))
	return fmt.Sprintf("%s.svg", filename)
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

type Config struct {
	InputPathDir  string
	OutputPathDir string
	Recursive     bool
}

func parseFlags() Config {
	var inputPathDir string
	var outputPathDir string
	var recursive bool
	flag.StringVar(&inputPathDir, "i", "", "Path containing .d2 files to be converted to .svg")
	flag.StringVar(&outputPathDir, "o", "", "Path in which .svg files will be put")
	flag.BoolVar(&recursive, "r", false, "Should recursively walk subdirs")
	flag.Parse()
	if !flag.Parsed() {
		log.Fatal("Failed to parse flags")
	}
	if inputPathDir == "" {
		log.Fatal("No input path specified")
	}
	if outputPathDir == "" {
		log.Fatal("No output path specified")
	}
	return Config{inputPathDir, outputPathDir, recursive}
}

func processFiles(inputPaths []string, config Config) {
	var wg sync.WaitGroup
	for _, d2FilePath := range inputPaths {
		wg.Add(1)
		go func(filePath string) {
			defer wg.Done()
			processD2File(filePath, config)
		}(d2FilePath)
	}

	wg.Wait()
}

func processD2File(d2FilePath string, config Config) {
	fullInputPath := filepath.Join(config.InputPathDir, d2FilePath)
	fullOutputPath := filepath.Join(config.OutputPathDir, paths.ReplaceExtension(d2FilePath, ".svg"))
	generateDiagram(fullInputPath, fullOutputPath)
}
func main() {
	_, err := exec.LookPath("d2")
	if err != nil {
		log.Fatal("Need to have d2lang (available here: https://d2lang.com/) installed to generate the diagrams.")
	}

	config := parseFlags()

	cwd, err := os.Getwd()
	if err != nil {
		log.Println("Error while getting current working directory:", err)
	}
	// If we passed in a input or output path, make it relative to the dir from which the command was run
	config.InputPathDir = filepath.Join(cwd, config.InputPathDir)
	config.OutputPathDir = filepath.Join(cwd, config.OutputPathDir)

	inputPaths := paths.GetInputPaths(os.DirFS(config.InputPathDir), ".d2", config.Recursive)
	processFiles(inputPaths, config)
}
