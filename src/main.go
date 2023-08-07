// This file handles the conversion of Graph in YAML format to HTML

// Important uncommon shortforms used:
// GDF - Graph Definition File (usually in YAML)
package main

import (
	// "html/template"
	"log"

	"gopkg.in/yaml.v2"
	"io/ioutil"
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

	// log.Printf("%+v", gdfData)
	nodes, err := createComputeAndFillNodeDataList(gdfData)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}

	packedData, err := yaml.Marshal(nodes)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}

	err = ioutil.WriteFile("/tmp/output.yaml", packedData, 0644)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}

}
