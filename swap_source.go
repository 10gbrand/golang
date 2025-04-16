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

func SwapSources(config Config) error {
	sourceMapData, err := os.ReadFile(config.SourceConfig)
	if err != nil {
		return fmt.Errorf("kunde inte läsa source map: %w", err)
	}

	var sourceMap SourceMap
	if err := json.Unmarshal(sourceMapData, &sourceMap); err != nil {
		return fmt.Errorf("kunde inte parsa source map: %w", err)
	}

	// Samla alla miljöer utom "localhost"
	environments := make(map[string]struct{})
	for _, source := range sourceMap {
		for env := range source {
			if env != "localhost" {
				environments[env] = struct{}{}
			}
		}
	}

	// Bearbeta varje miljö
	for env := range environments {
		outputDir := filepath.Join(config.ResultDir, env)
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("kunde inte skapa mapp för %s: %w", env, err)
		}

		err := filepath.WalkDir(config.ResultDir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if !d.IsDir() && filepath.Ext(path) == ".json" {
				if err := swapSourcesInFile(path, sourceMap, outputDir, env); err != nil {
					fmt.Fprintf(os.Stderr, "Fel i %s (%s): %v\n", path, env, err)
				}
			}
			return nil
		})

		if err != nil {
			return fmt.Errorf("bearbetning misslyckades för %s: %w", env, err)
		}
	}
	return nil
}

func swapSourcesInFile(inputPath string, sourceMap SourceMap, outputDir string, targetEnv string) error {
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("läsfel: %w", err)
	}

	var style StyleFile
	if err := json.Unmarshal(data, &style); err != nil {
		return fmt.Errorf("json-fel: %w", err)
	}

	newSources := make(map[string]interface{})
	for _, layer := range style.Layers {
		sourceName, ok := layer["source"].(string)
		if !ok {
			continue
		}

		if envs, exists := sourceMap[sourceName]; exists {
			if cfg, found := envs[targetEnv]; found {
				newSources[sourceName] = cfg
			}
		}
	}

	style.Sources = newSources
	outputData, err := json.MarshalIndent(style, "", "  ")
	if err != nil {
		return fmt.Errorf("json-serialiseringsfel: %w", err)
	}

	outputPath := filepath.Join(outputDir, filepath.Base(inputPath))
	return os.WriteFile(outputPath, outputData, 0644)
}
