package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

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
	var inputPathDir string
	var outputPathDir string
	var recursive bool
	flag.StringVar(&inputPathDir, "i", "", "Path containing .json files to be converted to .md")
	flag.StringVar(&outputPathDir, "o", "", "Path in which .md files will be put")
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
	config.InputPathDir = filepath.Join(cwd, config.InputPathDir)
	config.OutputPathDir = filepath.Join(cwd, config.OutputPathDir)
	inputPaths := paths.GetInputPaths(os.DirFS(config.InputPathDir), ".json", config.Recursive)
	processFiles(inputPaths, config)
}
