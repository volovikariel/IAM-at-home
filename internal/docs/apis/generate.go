package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/volovikariel/IdentityManager/internal/cli"
	"github.com/volovikariel/IdentityManager/internal/paths"
)

type Endpoint struct {
	Path        string     `json:"path"`
	Method      string     `json:"method"`
	Description string     `json:"description"`
	Requests    []Request  `json:"requests"`
	Responses   []Response `json:"responses"`
}

type Request struct {
	ContentType string      `json:"contentType"`
	Parameters  []Parameter `json:"parameters"`
}

type Parameter struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Response struct {
	Status string `json:"status"`
	Body   string `json:"body"`
}

func generateMarkdown(endpoints []Endpoint) string {
	var sb strings.Builder
	sb.WriteString("# API Documentation\n\n")
	for _, endpoint := range endpoints {
		sb.WriteString("## " + endpoint.Method + " " + endpoint.Path + "\n")
		sb.WriteString("### Description\n")
		sb.WriteString(endpoint.Description + "\n\n")
		sb.WriteString("### Request\n")
		for _, request := range endpoint.Requests {
			sb.WriteString("#### Content Type\n")
			sb.WriteString(request.ContentType + "\n\n")
			sb.WriteString("#### Parameters\n")
			for _, param := range request.Parameters {
				sb.WriteString("- " + param.Name + " (`" + param.Type + "`): " + param.Value + "\n")
			}
		}
		sb.WriteString("\n### Responses\n")
		for _, response := range endpoint.Responses {
			sb.WriteString("- " + response.Status + "\n")
			if response.Body != "" {
				sb.WriteString("```json\n" + response.Body + "\n```\n")
			}
		}
		sb.WriteString("\n---\n\n")
	}
	return sb.String()
}

func main() {
	_, currentFile, _, _ := runtime.Caller(0)
	internalApisDir := filepath.Dir(currentFile)

	defaultInputPathDir := filepath.Join(internalApisDir, "input")
	defaultOutputPathDir := filepath.Join(internalApisDir, "../../../docs/apis")
	recursiveDefault := false

	var inputPathDir string
	var outputPathDir string
	var recursive bool
	flag.StringVar(&inputPathDir, "i", defaultInputPathDir, "Path containing .json files to be converted to .md")
	flag.StringVar(&outputPathDir, "o", defaultOutputPathDir, "Path in which .md files will be put")
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
	if cli.IsFlagPassed("i") {
		inputPathDir = filepath.Join(cwd, inputPathDir)
	}
	if cli.IsFlagPassed("o") {
		outputPathDir = filepath.Join(cwd, outputPathDir)
	}

	inputPaths := paths.GetInputPaths(os.DirFS(inputPathDir), ".json", recursive)
	for _, apiFilePath := range inputPaths {
		apiFileFullPath := filepath.Join(inputPathDir, apiFilePath)

		jsonFile, err := os.Open(apiFileFullPath)
		if err != nil {
			log.Fatalln("Error opening JSON file:", err)
		}
		defer jsonFile.Close()

		byteValue, _ := io.ReadAll(jsonFile)

		var endpoints []Endpoint
		json.Unmarshal(byteValue, &endpoints)

		markdown := generateMarkdown(endpoints)

		markdownFullFilePath := filepath.Join(outputPathDir, paths.ReplaceExtension(apiFilePath, ".md"))
		log.Printf("Generating Markdown file: %s from %s\n", markdownFullFilePath, apiFileFullPath)

		if err := os.MkdirAll(filepath.Dir(markdownFullFilePath), os.ModePerm); err != nil {
			log.Fatalln("Error creating directory:", err)
		}
		if err = os.WriteFile(markdownFullFilePath, []byte(markdown), 0644); err != nil {
			log.Fatalln("Error writing Markdown file:", err)
		}
	}
}
