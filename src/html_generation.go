// This file deals with the generation of the HTML document based on the processed data
package main

import (
	"html/template"
	"os"
)

// All the data required for generating HTML page from template is stored in this struct
type TemplateData struct {
	HeadConfig *HeadConfigFields
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
