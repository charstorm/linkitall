// This file handles the conversion of Graph in YAML format to HTML

// Important uncommon shortforms used:
// GDF - Graph Definition File (usually in YAML)
package main

import (
	// "html/template"
	"log"
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

	log.Printf("%+v", gdfData)
}
