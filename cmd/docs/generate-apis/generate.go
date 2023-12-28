package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/volovikariel/IdentityManager/internal/cli"
	"github.com/volovikariel/IdentityManager/internal/paths"
)

type API struct {
	Name      string     `json:"name"`
	Endpoints []Endpoint `json:"endpoints"`
}

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

func generateMarkdown(api API) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# %s\n\n", api.Name))
	for _, endpoint := range api.Endpoints {
		sb.WriteString("## " + endpoint.Method + " " + endpoint.Path + "\n")
		if endpoint.Description != "" {
			sb.WriteString("### Description\n")
			sb.WriteString(endpoint.Description + "\n\n")
		}
		if len(endpoint.Requests) > 0 {
			sb.WriteString("### Request\n")
			for _, request := range endpoint.Requests {
				sb.WriteString("#### Content Type\n")
				sb.WriteString(request.ContentType + "\n\n")
				sb.WriteString("#### Parameters\n")
				for _, param := range request.Parameters {
					sb.WriteString("- " + param.Name + " (`" + param.Type + "`): " + param.Value + "\n")
				}
			}

		}
		if len(endpoint.Responses) > 0 {
			sb.WriteString("\n### Responses\n")
			for _, response := range endpoint.Responses {
				sb.WriteString("- " + response.Status + "\n")
				if response.Body != "" {
					sb.WriteString("```json\n" + response.Body + "\n```\n")
				}
			}
		}
		sb.WriteString("\n---\n\n")
	}
	return sb.String()
}

type Config struct {
	InputPathDir  string
	OutputPathDir string
	Recursive     bool
}

func parseFlags() Config {
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
	return Config{inputPathDir, outputPathDir, recursive}
}

func makeRelative(dir1 string, dir2 string) string {
	return filepath.Join(dir1, dir2)
}

func processFiles(inputPaths []string, config Config) {
	var wg sync.WaitGroup
	for _, apiFilePath := range inputPaths {
		wg.Add(1)
		go func(filePath string) {
			defer wg.Done()
			processJSONFile(filePath, config)
		}(apiFilePath)
	}

	wg.Wait()
}

func processJSONFile(apiFilePath string, config Config) {
	apiFileFullPath := filepath.Join(config.InputPathDir, apiFilePath)

	jsonFile, err := os.Open(apiFileFullPath)
	defer jsonFile.Close()
	if err != nil {
		log.Fatalln("Error opening JSON file:", err)
	}

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Fatalln("Error reading JSON file:", err)
	}

	var api API
	err = json.Unmarshal(byteValue, &api)
	if err != nil {
		log.Fatalln("Error unmarshalling JSON file:", err)
	}

	markdown := generateMarkdown(api)

	markdownFullFilePath := filepath.Join(config.OutputPathDir, paths.ReplaceExtension(apiFilePath, ".md"))
	log.Printf("Generating Markdown file: %s from %s\n", markdownFullFilePath, apiFileFullPath)

	if err := os.MkdirAll(filepath.Dir(markdownFullFilePath), os.ModePerm); err != nil {
		log.Fatalln("Error creating directory:", err)
	}
	if err = os.WriteFile(markdownFullFilePath, []byte(markdown), 0644); err != nil {
		log.Fatalln("Error writing Markdown file:", err)
	}
}

func main() {
	config := parseFlags()
	cwd, err := os.Getwd()
	if err != nil {
		log.Println("Error while getting current working directory:", err)
	}
	// If we passed in a input or output path, make it relative to the dir from which the command was run
	if cli.IsFlagPassed("i") {
		fmt.Println(cwd, config.InputPathDir)
		config.InputPathDir = makeRelative(cwd, config.InputPathDir)
	}
	if cli.IsFlagPassed("o") {
		config.OutputPathDir = makeRelative(cwd, config.OutputPathDir)
	}
	inputPaths := paths.GetInputPaths(os.DirFS(config.InputPathDir), ".json", config.Recursive)
	processFiles(inputPaths, config)
}
