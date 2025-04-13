package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestReadCSV(t *testing.T) {
	// Skapa tillfällig CSV-fil
	tempDir := t.TempDir()
	testCSVPath := filepath.Join(tempDir, "test.csv")
	content := `aktiv,src,target,order
1,test1,test_target,1.0
0,test2,ignored,2.0
1,test3,test_target,0.5
`

	err := os.WriteFile(testCSVPath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to write temp CSV file: %v", err)
	}

	entries, err := readCSV(testCSVPath)
	if err != nil {
		t.Fatalf("readCSV failed: %v", err)
	}

	if len(entries) != 3 {
		t.Errorf("Expected 3 entries, got %d", len(entries))
	}
}

func TestGroupEntries(t *testing.T) {
	entries := []CSVEntry{
		{Aktiv: 1, TargetFile: "a"},
		{Aktiv: 1, TargetFile: "b"},
		{Aktiv: 0, TargetFile: "a"},
		{Aktiv: 1, TargetFile: "a"},
	}

	grouped := groupEntries(entries)

	if len(grouped["a"]) != 2 {
		t.Errorf("Expected 2 entries for target 'a', got %d", len(grouped["a"]))
	}
	if len(grouped["b"]) != 1 {
		t.Errorf("Expected 1 entry for target 'b', got %d", len(grouped["b"]))
	}
}

func TestSortEntries(t *testing.T) {
	entries := []CSVEntry{
		{Order: 5.0},
		{Order: 1.0},
		{Order: 3.0},
	}

	sortedAsc := sortEntries(entries, true)
	if sortedAsc[0].Order != 1.0 || sortedAsc[2].Order != 5.0 {
		t.Errorf("Ascending sort failed: %+v", sortedAsc)
	}

	sortedDesc := sortEntries(entries, false)
	if sortedDesc[0].Order != 5.0 || sortedDesc[2].Order != 1.0 {
		t.Errorf("Descending sort failed: %+v", sortedDesc)
	}
}

func TestUpdateJSONStructure(t *testing.T) {
	base := map[string]interface{}{
		"layers":   []interface{}{},
		"sources":  map[string]interface{}{},
		"metadata": map[string]interface{}{},
	}

	layers := []interface{}{"layer1"}
	sources := map[string]interface{}{"src1": "data"}
	groups := []interface{}{"group1"}

	updateJSONStructure(base, layers, sources, groups)

	if !reflect.DeepEqual(base["layers"], layers) {
		t.Errorf("Expected layers to be updated")
	}
	if !reflect.DeepEqual(base["sources"], sources) {
		t.Errorf("Expected sources to be updated")
	}

	meta := base["metadata"].(map[string]interface{})
	if !reflect.DeepEqual(meta["springGroups"], groups) {
		t.Errorf("Expected springGroups to be updated")
	}
}

func TestGetBaseJSON(t *testing.T) {
	// Sätt global path variabel
	baseJSONPath = t.TempDir()

	jsonData := `{
		"layers": ["l1"],
		"sources": {"s1": "data"},
		"metadata": {
			"springGroups": ["g1"]
		}
	}`
	jsonFile := filepath.Join(baseJSONPath, "test.json")
	err := os.WriteFile(jsonFile, []byte(jsonData), 0644)
	if err != nil {
		t.Fatalf("Failed to write JSON file: %v", err)
	}

	data := getBaseJSON("test")

	if data["layers"] == nil || data["sources"] == nil {
		t.Errorf("Expected valid JSON data, got: %v", data)
	}
}

func TestWriteOutput(t *testing.T) {
	resultDir = t.TempDir()

	data := map[string]interface{}{
		"test": "value",
	}

	writeOutput("output_test", data)

	outputPath := filepath.Join(resultDir, "output_test.json")
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read written file: %v", err)
	}

	var out map[string]interface{}
	err = json.Unmarshal(content, &out)
	if err != nil {
		t.Fatalf("Failed to unmarshal output: %v", err)
	}

	if out["id"] == nil {
		t.Errorf("Expected output JSON to include 'id'")
	}
	if out["test"] != "value" {
		t.Errorf("Expected 'test' field to be 'value', got %v", out["test"])
	}
}
