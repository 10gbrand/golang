package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

type StyleFile struct {
	Sources map[string]interface{}   `json:"sources"`
	Layers  []map[string]interface{} `json:"layers"`
}

type SourceMap map[string]map[string]map[string]interface{}

func SwapSources(config Config, sourceMapPath string) error {
	sourceMapData, err := os.ReadFile(sourceMapPath)
	if err != nil {
		return fmt.Errorf("kunde inte läsa source map: %w", err)
	}

	var sourceMap SourceMap
	if err := json.Unmarshal(sourceMapData, &sourceMap); err != nil {
		return fmt.Errorf("kunde inte parsa source map: %w", err)
	}

	// Skapa app-mappen om den inte finns
	outputDir := filepath.Join(config.ResultDir, "app")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("kunde inte skapa app-mappen: %w", err)
	}

	// Loopa igenom alla JSON-filer i resultDir
	err = filepath.Walk(config.ResultDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(path) == ".json" && !info.IsDir() {
			if err := swapSourcesInFile(path, sourceMap, outputDir); err != nil {
				fmt.Fprintf(os.Stderr, "Kunde inte byta sources i filen %s: %v\n", path, err)
			} else {
				fmt.Printf("Bytte sources i: %s\n", path)
			}
		}

		return nil
	})

	return err
}

func swapSourcesInFile(inputPath string, sourceMap SourceMap, outputDir string) error {
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("kunde inte läsa fil: %w", err)
	}

	var style StyleFile
	if err := json.Unmarshal(data, &style); err != nil {
		return fmt.Errorf("kunde inte parsa JSON: %w", err)
	}

	newSources := make(map[string]interface{})
	for _, layer := range style.Layers {
		sourceName, ok := layer["source"].(string)
		if !ok {
			continue
		}

		if source, found := sourceMap[sourceName]; found {
			if appSrc, exists := source["app"]; exists {
				newSources[sourceName] = appSrc
			}
		}
	}
	style.Sources = newSources

	updatedData, err := json.MarshalIndent(style, "", "  ")
	if err != nil {
		return fmt.Errorf("kunde inte skapa JSON: %w", err)
	}

	// Samma filnamn, men skriv till "app"-mappen
	filename := filepath.Base(inputPath)
	outputPath := filepath.Join(outputDir, filename)

	if err := os.WriteFile(outputPath, updatedData, 0644); err != nil {
		return fmt.Errorf("kunde inte skriva fil: %w", err)
	}

	return nil
}
