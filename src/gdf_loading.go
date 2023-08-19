// This file handles the loading of Graph Definition File (GDF)
// Contains loading of data to struct, checking input values, and filling optional fields.
// Some of the checks may be redundant as similar checks may be present in modules that actually
// use the data.
package main

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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

type LinkToFields struct {
	// Resource name to be linked to
	ResourceName string `yaml:"resource"`
	// A target for the final resource (page/section/div-id) etc.
	Target string `yaml:"target"`
}

type ResourceConfigMap map[string]string

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
	// Link to the resource
	LinkTo LinkToFields `yaml:"linkto,omitempty"`
}

type AlgoConfigFields struct {
	LevelStrategy string `yaml:"level-strategy,omitempty"`
}

type GdfDataStruct struct {
	Nodes          []NodeInputFields   `yaml:"nodes"`
	HeadConfig     HeadConfigFields    `yaml:"head-config"`
	DisplayConfig  DisplayConfigFields `yaml:"display-config,omitempty"`
	ResourceConfig ResourceConfigMap   `yaml:"resources"`
	AlgoConfig     AlgoConfigFields    `yaml:"algo-config,omitempty"`
}

func validateAndUpdateAlgoConfig(algoConfig *AlgoConfigFields) error {
	if len(algoConfig.LevelStrategy) == 0 {
		algoConfig.LevelStrategy = "bottom2top"
	}
	if algoConfig.LevelStrategy != "bottom2top" && algoConfig.LevelStrategy != "top2bottom" {
		return fmt.Errorf("invalid level strategy: '%v'", algoConfig.LevelStrategy)
	}
	return nil
}

func validateAndUpdateDisplayConfig(displayConfig *DisplayConfigFields) error {
	if displayConfig.HorizontalStepPx == 0 {
		displayConfig.HorizontalStepPx = 400
	}
	if displayConfig.VerticalStepPx == 0 {
		displayConfig.VerticalStepPx = 300
	}
	if displayConfig.NodeBoxWidthPx == 0 {
		displayConfig.NodeBoxWidthPx = 300
	}
	return nil
}

// Replace underscores with spaces and capitalize first letter of every word.
func convertNameToTitle(name string) string {
	spacedStr := strings.ReplaceAll(name, "_", " ")
	// We have to do some extra-steps to do since strings.Title is deprecated. :eyeroll:
	caser := cases.Title(language.English)
	result := caser.String(spacedStr)
	return result
}

// Validate data related to nodes in GDF
// This function changes blank ("") value for node.Importance to "normal".
func validateAndUpdateNodes(nodes []NodeInputFields) error {
	// Number of nodes without any dependencies (level 0 nodes)
	numLevel0Nodes := 0
	// Unique nodes (names)
	uniqueNames := map[string]bool{}
	for idx := range nodes {
		node := &nodes[idx]

		// CHECK: node name must be [a-zA-Z0-9_]
		if !name_pattern.MatchString(node.Name) {
			return fmt.Errorf("invalid node name (only letters, numbers, _) '%v'", node.Name)
		}

		// CHECK: node name must be unique
		if _, ok := uniqueNames[node.Name]; ok {
			return fmt.Errorf("node name repeated '%v'", node.Name)
		}
		uniqueNames[node.Name] = true

		if node.Importance == "" {
			node.Importance = "normal"
		}
		// CHECK: importance must be one of the 7 options
		if !importance_pattern.MatchString(node.Importance) {
			return fmt.Errorf("unknown importance pattern for node '%v': '%v'",
				node.Name, node.Importance)
		}

		if len(node.DependsOn) == 0 {
			// These nodes do not depend on any nodes
			numLevel0Nodes += 1
		}

		// If the title is not given, fill it using the name.
		if len(node.Title) == 0 {
			node.Title = convertNameToTitle(node.Name)
		}
	}

	if numLevel0Nodes == 0 {
		return fmt.Errorf("these must be atleast 1 node without any dependency")
	}

	for _, node := range nodes {
		for _, dep := range node.DependsOn {
			// CHECK: dependency must be one of the node names
			if _, ok := uniqueNames[dep]; !ok {
				return fmt.Errorf("unknown dependency for node '%v': '%v'", node.Name, dep)
			}
		}
	}

	return nil
}

// Goes over each source in resources and makes sure the input is proper.
// Also iterates over the nodes and makes sure all the resources are available.
// TODO: check if the specified resource file actually exists!
func validateAndUpdateResources(resources ResourceConfigMap, nodes []NodeInputFields) error {
	// Check all nodes are using resources actually present in the GDF
	for idx := range nodes {
		node := &nodes[idx]
		if len(node.LinkTo.ResourceName) == 0 {
			continue
		}
		_, ok := resources[node.LinkTo.ResourceName]
		if !ok {
			return fmt.Errorf("error in node %s: linkto resource %s not found",
				node.Name, node.LinkTo.ResourceName)
		}
	}

	return nil
}

// Validate graph data loaded from YAML
// Input must not be nil.
func validateAndUpdateGraphData(data *GdfDataStruct) error {
	err := validateAndUpdateNodes(data.Nodes)
	if err != nil {
		return err
	}

	err = validateAndUpdateDisplayConfig(&data.DisplayConfig)
	if err != nil {
		return err
	}

	err = validateAndUpdateAlgoConfig(&data.AlgoConfig)
	if err != nil {
		return err
	}

	err = validateAndUpdateResources(data.ResourceConfig, data.Nodes)
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
func loadGdf(filename string) (*GdfDataStruct, bool, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, false, err
	}
	defer file.Close()

	fileData, err := io.ReadAll(file)
	if err != nil {
		return nil, false, err
	}

	var data GdfDataStruct
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
