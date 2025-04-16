package main

type Config struct {
	CSVPath      string `json:"csvPath"`
	BaseJSONPath string `json:"baseJSONPath"`
	ResultDir    string `json:"resultDir"`
	SourceConfig string `json:"sourceConfig"`
}

type CSVEntry struct {
	Aktiv      int
	SrcFile    string
	TargetFile string
	Order      float64
}
