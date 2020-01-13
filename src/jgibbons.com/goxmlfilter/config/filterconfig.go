package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type FilterConfig struct {
	RowsInFile         int            `json:"rows_in_file"`
	NumFiles           int            `json:"num_files"`
	RestQuery          string         `json:"query"`
	IgnoreCertAuthoriy bool           `json:"ignore_cert_authority"`
	DelimTag           string         `json:"first_tag"`
	DebugOutput        bool           `json:"debug_output"`
	FiltersEquals      []FilterField  `json:"filter_equals"`
	FiltersNotEquals   []FilterField  `json:"filter_not_equals"`
	ExtractColumns     []ExtractField `json:"filter_extract"`
}

type FilterField struct {
	Name  string `json:"element"`
	Value string `json:"value"`
}
type ExtractField struct {
	InputName string `json:"element"`
	OuputName string `json:"as"`
}

func New(fileName string) (*FilterConfig, error) {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Printf("Failed to read config file [%s]: %v", fileName, err)
		return nil, err
	}

	data := FilterConfig{}
	if err = json.Unmarshal([]byte(file), &data); err != nil {
		log.Printf("Failed to decode file [%s] into json : %v", fileName, err)
		return nil, err
	}

	return &data, nil
}

func (data *FilterConfig) DumpConfig() {
	fmt.Printf("URL [%s]\n", data.RestQuery)
	fmt.Printf("IgnoreCertAuthority [%t]\n", data.IgnoreCertAuthoriy)
	fmt.Printf("RowsInFile [%d]\n", data.RowsInFile)
	fmt.Printf("NumFiles [%d]\n", data.NumFiles)
	fmt.Printf("DelimTag [%s]\n", data.DelimTag)
	fmt.Printf("DebugOutput [%t]\n", data.DebugOutput)

	fmt.Printf("Filters:\n")
	for i := 0; i < len(data.FiltersEquals); i++ {
		f := data.FiltersEquals[i]
		fmt.Printf("  %s = %s\n", f.Name, f.Value)
	}
	for i := 0; i < len(data.FiltersNotEquals); i++ {
		f := data.FiltersNotEquals[i]
		fmt.Printf("  %s != %s\n", f.Name, f.Value)
	}
	fmt.Println("Extract:")
	for i := 0; i < len(data.ExtractColumns); i++ {
		f := data.ExtractColumns[i]
		fmt.Printf("  %s AS %s\n", f.InputName, f.OuputName)
	}
}
