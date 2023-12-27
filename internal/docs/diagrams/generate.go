package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/volovikariel/IdentityManager/internal/cli"
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

func main() {
	_, err := exec.LookPath("d2")
	if err != nil {
		log.Fatal("Need to have d2lang (available here: https://d2lang.com/) installed to generate the diagrams.")
	}

	_, filename, _, _ := runtime.Caller(0)
	internalDiagramsDir := filepath.Dir(filename)

	defaultInputPathDir := filepath.Join(internalDiagramsDir, "input")
	defaultOutputPathDir := filepath.Join(internalDiagramsDir, "../../../docs/diagrams")
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
		log.Fatalln("Error while getting current working directory:", err)
	}
	if cli.IsFlagPassed("i") {
		inputPathDir = filepath.Join(cwd, inputPathDir)
	}
	if cli.IsFlagPassed("o") {
		outputPathDir = filepath.Join(cwd, outputPathDir)
	}

	inputPaths := paths.GetInputPaths(os.DirFS(inputPathDir), ".d2", recursive)
	for _, inputPath := range inputPaths {
		fullInputPath := filepath.Join(inputPathDir, inputPath)
		fullOutputPath := filepath.Join(outputPathDir, paths.ReplaceExtension(inputPath, ".svg"))
		generateDiagram(fullInputPath, fullOutputPath)
	}
}
