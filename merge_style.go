package main

import (
	"path/filepath"
)

func RunMerge(config Config) error {
	entries, err := readCSV(config.CSVPath)
	if err != nil {
		return err
	}

	targets := groupEntries(entries)

	for target, entries := range targets {
		processTarget(config, target, entries)
	}
	return nil
}

func processTarget(config Config, target string, entries []CSVEntry) {
	if len(entries) == 0 {
		return
	}

	sortedForLayers := sortEntries(entries, true)
	mergedLayers, mergedSources := mergeLayersAndSources(config, sortedForLayers)

	sortedForGroups := sortEntries(entries, false)
	mergedSpringGroups := mergeSpringGroups(config, sortedForGroups)

	// Ta bort ".json" från slutet av target om det finns
	targetBase := target
	if filepath.Ext(targetBase) == ".json" {
		targetBase = targetBase[:len(targetBase)-len(".json")]
	}

	baseFile := filepath.Join("def", "mall_"+targetBase)
	baseData := getBaseJSON(config, baseFile)

	updateJSONStructure(baseData, mergedLayers, mergedSources, mergedSpringGroups)
	writeOutput(config, target, baseData)
}
