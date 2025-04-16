---
title: Promtp
tags:
    - prompt
    - ai
    - golang
---

# Prompt

## Fråga 1

Bygg en golang app som slår samman delar av json-filer som specificeras i ./layer_styles/def/stylefiles.csv

stylefiles.csv:

```csv
aktiv,src_file,target_file,order
1,geonote_basvag_style,springfield_geonote_style,2
1,geonote_hansyn_style,springfield_geonote_style,1
1,geonote_avdelning_style,springfield_geonote_style,0
1,geonote_ovrigt_style,springfield_geonote_style,3
1,arter_style,springfield_geodatapackage_style,50
1,bakgrundskarta_style,springfield_geodatapackage_style,0
1,skifte_style,springfield_geodatapackage_style,1
1,avdelning_style,springfield_geodatapackage_style,2
1,atgard_style,springfield_geodatapackage_style,3
1,anmalan_style,springfield_geodatapackage_style,3.5
1,restriktion_style,springfield_geodatapackage_style,4
1,fornlamning_style,springfield_geodatapackage_style,4.2
1,kulturminne_style,springfield_geodatapackage_style,4.5
1,hansyn_style,springfield_geodatapackage_style,5
0,vag_style,springfield_geodatapackage_style,6
1,drivning_style,springfield_geodatapackage_style,7
1,text_style,springfield_geodatapackage_style,99
```

Exempel json geonote_basvag_style.json:

```json
{
  "$schema": "./spring_style_schema.json",
  "id": "0b11a362d0e44047a2ef65a3850fc867",
  "name": "Mall",
  "version": 8,
  "metadata": {
    "maputnik:renderer": "mbgljs",
    "springGroups": [
      {
        "id": "geonote_basvag",
        "name": "Geonote basväg",
        "collapsed": true,
        "visible": true,
        "springLayers": [
          {
            "id": "geonote_basvag",
            "name": "Basväg",
            "visible": true
          },
          {
            "id": "geonote_avlagg_line",
            "name": "Avlägg",
            "visible": true
          }
        ]
      }
    ],
    "springStyleOrder": 1,
    "springStyleVersion": "0.1"
  },
  "zoom": 0.861983335785597,
  "pitch": 0,
  "center": [
    17.3145660472,
    62.91542224
  ],
  "sprite": "http://localhost:4200/assets/springfield/geonote/sprite/spring_sprites",
  "glyphs": "http://localhost:4200/assets/springfield/geonote/font-glyphs/glyphs/{fontstack}/{range}.pbf",
  "bearing": 0,
  "sources": {
    "geonotes": {
      "type": "geojson",
      "data": "http://localhost:4200/assets/springfield/geonote/geonotes.geojson"
    }
  },
  "layers": [
    {
      "id": "geonote_basvag_background_l",
      "type": "line",
      "metadata": {
        "maputnik:comment": "geonote_basvag urn:sveaskog:atgplan:drivningsinfotyp:basvag",
        "springLayer": "geonote_basvag"
      },
      "source": "geonotes",
      "minzoom": 13,
      "layout": {
        "visibility": "visible"
      },
      "paint": {
        "line-color": "rgba(250, 250, 0, 1)",
        "line-opacity": 1,
        "line-width": 2,
        "line-offset": 4
      },
      "filter": [
        "all",
        [
          "match",
          [
            "geometry-type"
          ],
          [
            "LineString",
            "MultiLineString"
          ],
          true,
          false
        ],
        [
          "==",
          [
            "get",
            "group"
          ],
          "urn:sveaskog:atgplan:drivningsinfotyp:basvag"
        ]
      ]
    },
    {
        "id": "geonote_basvag_background_r",
        "type": "line",
        "metadata": {
          "maputnik:comment": "geonote_basvag urn:sveaskog:atgplan:drivningsinfotyp:basvag",
          "springLayer": "geonote_basvag"
        },
        "source": "geonotes",
        "minzoom": 13,
        "layout": {
          "visibility": "visible"
        },
        "paint": {
          "line-color": "rgba(250, 250, 0, 1)",
          "line-opacity": 1,
          "line-width": 2,
          "line-offset": -4
        },
        "filter": [
          "all",
          [
            "match",
            [
              "geometry-type"
            ],
            [
              "LineString",
              "MultiLineString"
            ],
            true,
            false
          ],
          [
            "==",
            [
              "get",
              "group"
            ],
            "urn:sveaskog:atgplan:drivningsinfotyp:basvag"
          ]
        ]
      }
  ]
}
```

De nycklar i json som skall sammanfogas är "layers", "sources" och "springGroups".
vis sammanslagningen skall objekten innom "layers", "sources" sorteras stigande enligt order i stylefiles.csv medan "springGroups" skall sorteras fallande.

## Fråga 1

Använd ./layer_styles/def/mall4layer.json istället för att använda "sortedForLayers[0].SrcFile" i 

```golang
baseData := getBaseJSON(sortedForLayers[0].SrcFile)
```

i denna app:

```golang
package main

import (
	"encoding/csv"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

type CSVEntry struct {
	Aktiv      int
	SrcFile    string
	TargetFile string
	Order      float64
}

func main() {
	entries, err := readCSV("layer_styles/def/stylefiles.csv")
	if err != nil {
		log.Fatal(err)
	}

	targets := groupEntries(entries)

	for target, entries := range targets {
		processTarget(target, entries)
	}
}

func processTarget(target string, entries []CSVEntry) {
	if len(entries) == 0 {
		return
	}

	sortedForLayers := sortEntries(entries, true)
	mergedLayers, mergedSources := mergeLayersAndSources(sortedForLayers)

	sortedForGroups := sortEntries(entries, false)
	mergedSpringGroups := mergeSpringGroups(sortedForGroups)

	baseData := getBaseJSON(sortedForLayers[0].SrcFile)
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

// Implementerade hjälpfunktioner

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
		if i == 0 { // Skip header
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

		// Lägg till layers
		if l, ok := data["layers"].([]interface{}); ok {
			layers = append(layers, l...)
		}

		// Merge sources
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

	// Reverse order for springGroups
	for i, j := 0, len(groups)-1; i < j; i, j = i+1, j-1 {
		groups[i], groups[j] = groups[j], groups[i]
	}

	return groups
}

func getBaseJSON(srcFile string) map[string]interface{} {
	path := filepath.Join("layer_styles/def", srcFile+".json")
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
	outputPath := filepath.Join("layer_styles/def/result", target+".json")
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

```

## fråga 3

hur gör jag för att ge varje resultatfil en slumpad guid i "id": "j3f36e14-e3f5-43c1-84c0-50a9c80dc5c7" i

```golang
package main

import (
	"encoding/csv"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

type CSVEntry struct {
	Aktiv      int
	SrcFile    string
	TargetFile string
	Order      float64
}

func main() {
	entries, err := readCSV("layer_styles/def/stylefiles.csv")
	if err != nil {
		log.Fatal(err)
	}

	targets := groupEntries(entries)

	for target, entries := range targets {
		processTarget(target, entries)
	}
}

func processTarget(target string, entries []CSVEntry) {
	if len(entries) == 0 {
		return
	}

	sortedForLayers := sortEntries(entries, true)
	mergedLayers, mergedSources := mergeLayersAndSources(sortedForLayers)

	sortedForGroups := sortEntries(entries, false)
	mergedSpringGroups := mergeSpringGroups(sortedForGroups)

	baseData := getBaseJSON(sortedForLayers[0].SrcFile)
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

// Implementerade hjälpfunktioner

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
		if i == 0 { // Skip header
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

		// Lägg till layers
		if l, ok := data["layers"].([]interface{}); ok {
			layers = append(layers, l...)
		}

		// Merge sources
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

	// Reverse order for springGroups
	for i, j := 0, len(groups)-1; i < j; i, j = i+1, j-1 {
		groups[i], groups[j] = groups[j], groups[i]
	}

	return groups
}

func getBaseJSON(srcFile string) map[string]interface{} {
	path := filepath.Join("layer_styles/def", srcFile+".json")
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
	outputPath := filepath.Join("layer_styles/def/result", target+".json")
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

```