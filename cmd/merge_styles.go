package cmd

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// StyleData representerar strukturen för JSON-filerna
type StyleData struct {
	Metadata struct {
		SpringGroups []interface{} `json:"springGroups"`
	} `json:"metadata"`
	Sources map[string]interface{} `json:"sources"`
	Layers  []interface{}          `json:"layers"`
}

// StyleFile representerar en rad från stylefiles.csv
type StyleFile struct {
	Active     int     `csv:"aktiv"`
	SrcFile    string  `csv:"src_file"`
	TargetFile string  `csv:"target_file"`
	Order      float64 `csv:"order"`
}

func main() {
	csvFile := "def/stylefiles.csv"
	baseDir := "."
	styleFiles, err := readCSV(csvFile)
	if err != nil {
		fmt.Printf("Error reading CSV file: %v\n", err)
		return
	}

	// Filtrera endast aktiva rader
	activeFiles := filterActiveFiles(styleFiles)

	// Gruppera efter target_file
	groupedFiles := groupByTargetFile(activeFiles)

	for targetFile, srcFiles := range groupedFiles {
		err := mergeStyles(baseDir, targetFile, srcFiles)
		if err != nil {
			fmt.Printf("Error merging styles for %s: %v\n", targetFile, err)
			return
		}
	}

	fmt.Println("Merging completed successfully!")
}

// Läser CSV-filen och returnerar en lista med StyleFile-objekt
func readCSV(filePath string) ([]StyleFile, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1 // Tillåt varierande antal fält

	rows, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var styleFiles []StyleFile
	for i, row := range rows {
		if i == 0 { // Hoppa över header-raden
			continue
		}
		active := 0
		fmt.Sscanf(row[0], "%d", &active)
		order := 0.0
		fmt.Sscanf(row[3], "%f", &order)

		styleFiles = append(styleFiles, StyleFile{
			Active:     active,
			SrcFile:    row[1],
			TargetFile: row[2],
			Order:      order,
		})
	}
	return styleFiles, nil
}

// Filtrerar endast aktiva filer (aktiv == 1)
func filterActiveFiles(files []StyleFile) []StyleFile {
	var activeFiles []StyleFile
	for _, file := range files {
		if file.Active == 1 {
			activeFiles = append(activeFiles, file)
		}
	}
	return activeFiles
}

// Grupperar src_files efter target_file
func groupByTargetFile(files []StyleFile) map[string][]StyleFile {
	grouped := make(map[string][]StyleFile)
	for _, file := range files {
		grouped[file.TargetFile] = append(grouped[file.TargetFile], file)
	}

	for _, group := range grouped {
		sort.Slice(group, func(i, j int) bool {
			return group[i].Order < group[j].Order
		})
	}

	return grouped
}

// Slår samman JSON-filer och skriver till target_file
func mergeStyles(baseDir string, targetFile string, srcFiles []StyleFile) error {
	var mergedData StyleData

	for _, src := range srcFiles {
		srcPath := filepath.Join(baseDir, src.SrcFile+".json")
		data, err := readJSON(srcPath)
		if err != nil {
			return fmt.Errorf("error reading %s: %v", srcPath, err)
		}

		// Slå samman metadata.springGroups[]
		if len(data.Metadata.SpringGroups) > 0 {
			mergedData.Metadata.SpringGroups = append(mergedData.Metadata.SpringGroups, data.Metadata.SpringGroups...)
		}

		// Slå samman sources{}
		for key, value := range data.Sources {
			if mergedData.Sources == nil {
				mergedData.Sources = make(map[string]interface{})
			}
			if _, exists := mergedData.Sources[key]; !exists {
				mergedData.Sources[key] = value
			}
		}

		// Slå samman layers[]
		if len(data.Layers) > 0 {
			mergedData.Layers = append(mergedData.Layers, data.Layers...)
		}
	}

	targetPath := filepath.Join(baseDir, "result", targetFile+".json")
	return writeJSON(targetPath, mergedData)
}

// Läser JSON-fil och returnerar StyleData-objektet
func readJSON(filePath string) (*StyleData, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not read file %s: %v", filePath, err)
	}

	var style StyleData
	err = json.Unmarshal(data, &style)
	if err != nil {
		return nil, fmt.Errorf("could not parse JSON in file %s: %v", filePath, err)
	}

	return &style, nil
}
