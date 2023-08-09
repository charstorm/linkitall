// This file handles the conversion of Graph in YAML format to HTML

// Important uncommon shortforms used:
// GDF - Graph Definition File (usually in YAML)
package main

import (
	// "html/template"
	"log"

	"gopkg.in/yaml.v2"
	"os"
)

func main() {
	inputFile := "graph.yaml"

	gdfData, readable, err := loadGdf(inputFile)
	if !readable {
		log.Fatalf("Fatal, cannot proceed: %v\n", err)
	}

	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}

	nodes, err := createComputeAndFillNodeDataList(gdfData)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}

	templateData := newTemplateData(gdfData, nodes)

	packedData, err := yaml.Marshal(templateData)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}

	err = os.WriteFile("tmp/template_data.yaml", packedData, 0644)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}

	err = fillTemplateWriteOutput("assets/template.html", templateData, "tmp/output.html")
	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}

}
