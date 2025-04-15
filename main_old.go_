package main

import (
	"encoding/csv"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"

	"github.com/google/uuid"
)

type Config struct {
	CSVPath      string `json:"csvPath"`
	BaseJSONPath string `json:"baseJSONPath"`
	ResultDir    string `json:"resultDir"`
}

type CSVEntry struct {
	Aktiv      int
	SrcFile    string
	TargetFile string
	Order      float64
}

var config Config

func main() {
	var err error
	config, err = loadConfig("config.json")
	if err != nil {
		log.Fatalf("Kunde inte läsa konfigurationsfil: %v", err)
	}

	entries, err := readCSV(config.CSVPath)
	if err != nil {
		log.Fatal(err)
	}

	targets := groupEntries(entries)

	for target, entries := range targets {
		processTarget(target, entries)
	}
}

func loadConfig(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func processTarget(target string, entries []CSVEntry) {
	if len(entries) == 0 {
		return
	}

	sortedForLayers := sortEntries(entries, true)
	mergedLayers, mergedSources := mergeLayersAndSources(sortedForLayers)

	sortedForGroups := sortEntries(entries, false)
	mergedSpringGroups := mergeSpringGroups(sortedForGroups)

	//	baseData := getBaseJSON(config.DefaultBaseJSON)
	//
	// Strip ".json" om det redan finns
	targetBase := target
	if filepath.Ext(targetBase) == ".json" {
		targetBase = targetBase[:len(targetBase)-len(".json")]
	}

	// T.ex. mall_springfield_geodatapackage_style
	baseFile := filepath.Join("def", "mall_"+targetBase)
	baseData := getBaseJSON(baseFile)

	updateJSONStructure(baseData, mergedLayers, mergedSources, mergedSpringGroups)

	writeOutput(target, baseData)
}

func updateJSONStructure(base map[string]interface{}, layers []interface{}, sources map[string]interface{}, groups []interface{}) {
	base["layers"] = layers
	base["sources"] = sources

	if metadata, ok := base["metadata"].(map[string]interface{}); ok {
		metadata["springGroups"] = groups
	}
}

func readCSV(path string) ([]CSVEntry, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ','
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var entries []CSVEntry
	for i, record := range records {
		if i == 0 {
			continue
		}
		if len(record) != 4 {
			continue
		}

		aktiv, _ := strconv.Atoi(record[0])
		order, _ := strconv.ParseFloat(record[3], 64)

		entries = append(entries, CSVEntry{
			Aktiv:      aktiv,
			SrcFile:    record[1],
			TargetFile: record[2],
			Order:      order,
		})
	}
	return entries, nil
}

func groupEntries(entries []CSVEntry) map[string][]CSVEntry {
	groups := make(map[string][]CSVEntry)
	for _, entry := range entries {
		if entry.Aktiv == 1 {
			groups[entry.TargetFile] = append(groups[entry.TargetFile], entry)
		}
	}
	return groups
}

func sortEntries(entries []CSVEntry, ascending bool) []CSVEntry {
	sorted := make([]CSVEntry, len(entries))
	copy(sorted, entries)

	sort.Slice(sorted, func(i, j int) bool {
		if ascending {
			return sorted[i].Order < sorted[j].Order
		}
		return sorted[i].Order > sorted[j].Order
	})
	return sorted
}

func mergeLayersAndSources(entries []CSVEntry) ([]interface{}, map[string]interface{}) {
	layers := make([]interface{}, 0)
	sources := make(map[string]interface{})

	for _, entry := range entries {
		data := getBaseJSON(entry.SrcFile)

		if l, ok := data["layers"].([]interface{}); ok {
			layers = append(layers, l...)
		}

		if s, ok := data["sources"].(map[string]interface{}); ok {
			for k, v := range s {
				sources[k] = v
			}
		}
	}

	return layers, sources
}

func mergeSpringGroups(entries []CSVEntry) []interface{} {
	groups := make([]interface{}, 0)

	for _, entry := range entries {
		data := getBaseJSON(entry.SrcFile)

		if metadata, ok := data["metadata"].(map[string]interface{}); ok {
			if g, ok := metadata["springGroups"].([]interface{}); ok {
				groups = append(groups, g...)
			}
		}
	}

	for i, j := 0, len(groups)-1; i < j; i, j = i+1, j-1 {
		groups[i], groups[j] = groups[j], groups[i]
	}

	return groups
}

func getBaseJSON(srcFile string) map[string]interface{} {
	path := filepath.Join(config.BaseJSONPath, srcFile+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		log.Fatal(err)
	}
	return result
}

func writeOutput(target string, data map[string]interface{}) {
	randomID := uuid.New().String()
	data["id"] = randomID

	outputPath := filepath.Join(config.ResultDir, target+".json")
	err := os.MkdirAll(filepath.Dir(outputPath), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	if err := os.WriteFile(outputPath, jsonData, 0644); err != nil {
		log.Fatal(err)
	}
}
