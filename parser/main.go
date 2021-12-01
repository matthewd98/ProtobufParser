package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const url = "https://github.com/matthew-daddario/ProtobufParser"

func main() {
	mode := flag.Int("mode", 0, "0=crawl directory; 1=use embedded example file")
	dirToCrawl := flag.String("dir", "", "Directory containing .proto files")
	outputFileName := flag.String("output", "output.json", "Output JSON file name")
	flag.Parse()

	var doc *Documentation = nil

	switch *mode {
	case 0:
		doc = directoryMode(*dirToCrawl)
	case 1:
		doc = exampleMode()
	default:
		fmt.Println("Invalid mode")
		os.Exit(1)
	}

	outputFile, err := os.Create(*outputFileName)
	if err != nil {
		fmt.Printf("Failed to create output file. Reason: %s\n", err.Error())
		os.Exit(1)
	}
	defer outputFile.Close()

	bytes, _ := json.MarshalIndent(doc, "", "  ")
	outputFile.Write(bytes)

	wd, _ := os.Getwd()
	fmt.Printf("\nSuccessfully created file %s/%s\n", wd, *outputFileName)
}

func directoryMode(dirToCrawl string) *Documentation {
	fileInfo, err := os.Stat(dirToCrawl)
	if err != nil {
		fmt.Println("Invalid path specified")
		os.Exit(1)
	}

	if !fileInfo.IsDir() {
		fmt.Println("Specified path is not a directory")
		os.Exit(1)
	}

	doc := &Documentation{}

	var protoFiles []string
	err = filepath.Walk(dirToCrawl, func(filePath string, fileInfo os.FileInfo, err error) error {
		if !fileInfo.IsDir() && filepath.Ext(filePath) == ".proto" {
			protoFiles = append(protoFiles, filePath)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("File walk error : %s\n", err.Error())
		os.Exit(1)
	}
	for _, filePath := range protoFiles {
		data, readError := os.ReadFile(filePath)
		if readError != nil {
			fmt.Printf("File read error : %s\n", filePath)
			continue
		}
		fmt.Printf("Parsing file : %s\n", filePath)
		fileData := string(data)
		schema := parse(fileData)
		schema.FilePath, schema.FileName = path.Split(filePath)
		schema.Url = url + strings.Split(filePath, "/schema")[1]
		doc.Schema = append(doc.Schema, *schema)
	}

	return doc
}

func exampleMode() *Documentation {
	doc := &Documentation{}

	schema := parse(ProtoExample)
	schema.FileName = "Example file"
	doc.Schema = append(doc.Schema, *schema)

	return doc
}
