// This file handles the loading of Graph Definition File (GDF)
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"gopkg.in/yaml.v2"
)

var name_pattern = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
var importance_pattern = regexp.MustCompile(`^(lowest|lower|low|normal|high|higher|highest)$`)

// Used in the <head> of the final HTML
type HeadConfigFields struct {
	Title       string
	Description string
	Author      string
}

// Configuration related to positioning
type DisplayConfigFields struct {
	// Size of horizontal grid step
	HorizontalStepPx int `yaml:"horizontal-step-px,omitempty"`
	// Size of vertical grid step
	VerticalStepPx int `yaml:"vertical-step-px,omitempty"`
	// Width of the node box
	NodeBoxWidthPx int `yaml:"node-box-width-px,omitempty"`
}

// Defines the node definition by the user in the
type NodeInputFields struct {
	// A unique name for the node (no spaces, all small letters)
	Name string
	// Title of the node (shown in big font)
	Title string
	// Subtitle (shown in smaller font or sometimes omitted)
	Subtitle string `yaml:"subtitle,omitempty"`
	// Importance to be assigned to this node. It is a 7 point scale:
	// lowest, lower, low, normal, high, higher, highest
	Importance string `yaml:"importance,omitempty"`
	// List of node names (current node depends on these nodes)
	DependsOn []string `yaml:"depends-on,omitempty"`
}

type GdfData struct {
	Nodes         []NodeInputFields   `yaml:"nodes"`
	HeadConfig    HeadConfigFields    `yaml:"head-config"`
	DisplayConfig DisplayConfigFields `yaml:"display-config,omitempty"`
}

func validateAndUpdateDisplayConfig(displayConfig *DisplayConfigFields) error {
	defaultGridStepPx := 400
	if displayConfig.HorizontalStepPx == 0 {
		displayConfig.HorizontalStepPx = defaultGridStepPx
	}
	if displayConfig.VerticalStepPx == 0 {
		displayConfig.VerticalStepPx = defaultGridStepPx
	}
	if displayConfig.NodeBoxWidthPx == 0 {
		displayConfig.NodeBoxWidthPx = 300
	}
	return nil
}

// Validate graph data loaded from YAML
// Input must not be nil.
// This function changes blank ("") value for node.Importance to "normal".
func validateAndUpdateGraphData(data *GdfData) error {
	// Number of nodes without any dependencies (level 0 nodes)
	numLevel0Nodes := 0
	// Unique nodes (names)
	uniqueNames := map[string]bool{}
	for idx, _ := range data.Nodes {
		node := &data.Nodes[idx]

		// CHECK: node name must be [a-zA-Z0-9_]
		if !name_pattern.MatchString(node.Name) {
			return fmt.Errorf("Invalid node name (only letters, numbers, _) '%v'", node.Name)
		}

		// CHECK: node name must be unique
		if _, ok := uniqueNames[node.Name]; ok {
			return fmt.Errorf("Node name repeated '%v'", node.Name)
		}
		uniqueNames[node.Name] = true

		if node.Importance == "" {
			node.Importance = "normal"
		}
		// CHECK: importance must be one of the 7 options
		if !importance_pattern.MatchString(node.Importance) {
			return fmt.Errorf("Unknown importance pattern for node '%v': '%v'",
				node.Name, node.Importance)
		}

		if len(node.DependsOn) == 0 {
			// These nodes do not depend on any nodes
			numLevel0Nodes += 1
		}
	}

	if numLevel0Nodes == 0 {
		return fmt.Errorf("These must be atleast 1 node without any dependency")
	}

	for _, node := range data.Nodes {
		for _, dep := range node.DependsOn {
			// CHECK: dependency must be one of the node names
			if _, ok := uniqueNames[dep]; !ok {
				return fmt.Errorf("Unknown dependency for node '%v': '%v'", node.Name, dep)
			}
		}
	}

	err := validateAndUpdateDisplayConfig(&data.DisplayConfig)
	if err != nil {
		return err
	}

	return nil
}

// Load Graph Definition File
//
// Inputs:
// filename - input filename (YAML file for GDF)
//
// Returns: (data, readable, error)
// data - loaded data (if everything goes fine)
// readable - true if file is readable
// error - error if any
func loadGdf(filename string) (*GdfData, bool, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, false, err
	}
	defer file.Close()

	fileData, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, false, err
	}

	var data GdfData
	err = yaml.UnmarshalStrict(fileData, &data)
	if err != nil {
		return nil, true, err
	}

	err = validateAndUpdateGraphData(&data)
	if err != nil {
		return nil, true, err
	}

	return &data, true, nil
}
