package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestReadCSV(t *testing.T) {
	// Skapa en temporär CSV-fil
	content := `Aktiv,SrcFile,TargetFile,Order
1,file1,target1,2.0
1,file2,target1,1.0
0,file3,target2,3.0
`
	tmpDir := t.TempDir()
	csvPath := filepath.Join(tmpDir, "test.csv")
	if err := os.WriteFile(csvPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	entries, err := readCSV(csvPath)
	if err != nil {
		t.Fatalf("readCSV failed: %v", err)
	}

	if len(entries) != 3 {
		t.Errorf("Expected 3 entries, got %d", len(entries))
	}

	if entries[0].SrcFile != "file1" || entries[1].Order != 1.0 {
		t.Error("CSV parsing failed or order incorrect")
	}
}

func TestGroupEntries(t *testing.T) {
	entries := []CSVEntry{
		{Aktiv: 1, TargetFile: "A"},
		{Aktiv: 1, TargetFile: "B"},
		{Aktiv: 0, TargetFile: "A"},
		{Aktiv: 1, TargetFile: "A"},
	}
	grouped := groupEntries(entries)

	if len(grouped) != 2 {
		t.Errorf("Expected 2 groups, got %d", len(grouped))
	}

	if len(grouped["A"]) != 2 {
		t.Errorf("Expected 2 entries for group A, got %d", len(grouped["A"]))
	}
}

func TestSortEntries(t *testing.T) {
	entries := []CSVEntry{
		{Order: 5},
		{Order: 2},
		{Order: 9},
	}

	asc := sortEntries(entries, true)
	if asc[0].Order != 2 || asc[2].Order != 9 {
		t.Error("Ascending sort failed")
	}

	desc := sortEntries(entries, false)
	if desc[0].Order != 9 || desc[2].Order != 2 {
		t.Error("Descending sort failed")
	}
}

func TestUpdateJSONStructure(t *testing.T) {
	jsonStr := `{
		"layers": [],
		"sources": {},
		"metadata": {
			"title": "Example"
		}
	}`

	var base map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &base); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	newLayers := []interface{}{"layer1", "layer2"}
	newSources := map[string]interface{}{"source1": map[string]string{"type": "geojson"}}
	newGroups := []interface{}{"group1", "group2"}

	updateJSONStructure(base, newLayers, newSources, newGroups)

	if !reflect.DeepEqual(base["layers"], newLayers) {
		t.Errorf("Layers not updated correctly")
	}
	if !reflect.DeepEqual(base["sources"], newSources) {
		t.Errorf("Sources not updated correctly")
	}
	meta := base["metadata"].(map[string]interface{})
	if !reflect.DeepEqual(meta["springGroups"], newGroups) {
		t.Errorf("springGroups not updated correctly")
	}
}
