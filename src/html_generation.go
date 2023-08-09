// This file deals with the generation of the HTML document based on the processed data
package main

import (
	"html/template"
	"os"
)

// Info about the board (outer board used for holding all the nodes)
type BoardConfigFields struct {
	Width  int
	Height int
}

// All the data required for generating HTML page from template is stored in this struct
type TemplateData struct {
	// Input GDF Data
	GdfData *GdfDataStruct
	// With computed fields
	Nodes []NodeData
	// Board configuration
	BoardConfig BoardConfigFields
}

func computeBoardConfig(gdfData *GdfDataStruct, nodes []NodeData) BoardConfigFields {
	extraWidth := 10
	nodeBoxWidthPx := gdfData.DisplayConfig.NodeBoxWidthPx
	// Not configurable for now (but it shouldn't matter much anyway)
	nodeBoxHeightPx := 150
	maxWidth := 0
	maxHeight := 0
	// Initially compute max of left and top. Then add node width and height.
	for _, node := range nodes {
		if node.ElemFields.LeftPx > maxWidth {
			maxWidth = node.ElemFields.LeftPx
		}
		if node.ElemFields.TopPx > maxHeight {
			maxHeight = node.ElemFields.TopPx
		}
	}

	return BoardConfigFields{maxWidth + nodeBoxWidthPx + extraWidth, maxHeight + nodeBoxHeightPx}
}

// Constructor for TemplateData. There are some fields like BoardConfig that needs to be
// calculated
func newTemplateData(gdfData *GdfDataStruct, nodes []NodeData) TemplateData {
	boardConfig := computeBoardConfig(gdfData, nodes)
	return TemplateData{gdfData, nodes, boardConfig}
}

// The function responsible for generating the final HTML from template
func fillTemplateWriteOutput(templateFile string, data TemplateData, outputFile string) error {
	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		return err
	}

	writer, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer writer.Close()

	err = tmpl.Execute(writer, data)
	if err != nil {
		return err
	}

	return nil
}
