// This file handles the computation of various fields related to nodes.
// There are some more computations required for the final HTML generation. Those will be handled
// by html_generation.go.
package main

import (
	"fmt"
)

const defaultCapacity = 5
const defaultInvalidLevel = -1

// All the fields related to defining IDs for node and connections
type NodeIntIdFields struct {
	// Assign unique integer ID to every node. This is just the index
	Uid int
	// UIDs of all dependencies that this node depends on
	DependsOnIds []int
	// UIDs of all nodes that depend on this node
	UsedByIds []int
}

// Fields related to position (level, shift)
type NodePositionFields struct {
	// Level is for the vertical position. Level 0 is at the bottom (for axioms and such)
	Level int
	// Shift is for the horizontal position
	Shift int
}

// In the generated HTML file, every node has a set of connection dots.
// Each connection dot (a div in html) has an element ID of its own.
// Also we should keep track of the ID of the other node which is connected to this node.
type DotElemFields struct {
	// HTML ID of the Dot element
	DotElemId string
	// HTML ID of the partner node (not the node holding the dot)
	PartnerNodeId string
}

// For every node, we keep track of all the HTML data required.
// IDs in this struct are strings which will be mapped to HTML element IDs.
type NodeElemFields struct {
	// Element ID of the node
	NodeElemId string
	// Dots used for "depends-on" connections.
	// By default, these appear at the bottom of every node as a node depends on other nodes that
	// are of lower level (more fundamental).
	DependsOnDots []DotElemFields
	// Dots for "used-by" connections.
	// By default, these appear at the top of every node.
	UsedByDots []DotElemFields
	// Classes used by the node (HTML). This will be used to handle parameters like Importance.
	Classes string
	// Left edge position (px)
	LeftPx int
	// Top edge position (px)
	TopPx int
}

// All the data corresponding to a node
type NodeData struct {
	// Fields coming from input
	InputFields NodeInputFields
	// Unique number based IDs (computed)
	IntIdFields NodeIntIdFields
	// Position related fields (computed). Does not handle HTML related positions.
	Position NodePositionFields
	// HTML related fields (computed). Also handles positions on the HTML page.
	ElemFields NodeElemFields
}

// Create a list of NodeData based on GDF data
// Not handling the error.
func createNodeDataList(gdfData *GdfDataStruct) []NodeData {
	// This will contain the final result
	nodeDataSeq := make([]NodeData, 0, len(gdfData.Nodes))
	for _, node := range gdfData.Nodes {
		var nodeData NodeData
		nodeData.InputFields = node
		pushBack(&nodeDataSeq, nodeData)
	}
	return nodeDataSeq
}

// Fill all the fields related to integer IDs.
func fillIntIdFields(nodeDataSeq []NodeData) error {
	// A mapping of name -> integer unique ID
	nodeName2Id := map[string]int{}

	// Fill Uid. Also report if node names are repeated. This may not be really required
	// since there is already some error checking when reading the input data.
	for idx := range nodeDataSeq {
		node := &nodeDataSeq[idx]
		if _, found := nodeName2Id[node.InputFields.Name]; found {
			return fmt.Errorf("node name repeated '%v'", node.InputFields.Name)
		}
		nodeName2Id[node.InputFields.Name] = idx
		node.IntIdFields.Uid = idx
		// Initialize the arrays inside the struct
		node.IntIdFields.DependsOnIds = make([]int, 0, defaultCapacity)
		node.IntIdFields.UsedByIds = make([]int, 0, defaultCapacity)
	}

	// Fill DependsOnIds and UsedByIds using nodeName2Id
	for idx := range nodeDataSeq {
		node := &nodeDataSeq[idx]
		for _, depNodeName := range node.InputFields.DependsOn {
			depNodeId, ok := nodeName2Id[depNodeName]
			if !ok {
				return fmt.Errorf("dependency not found '%v'", depNodeName)
			}
			pushBack(&node.IntIdFields.DependsOnIds, depNodeId)
			depNode := &nodeDataSeq[depNodeId]
			pushBack(&depNode.IntIdFields.UsedByIds, idx)
		}
	}

	return nil
}

// Initialization step of computeLevels algorithm
// Assign level 0 to all nodes without dependencies. Everyone else gets an invalid level.
// Keep the level 0 nodes in nextLevelNodeIds.
func initializeForComputeLevels(nodes []NodeData, level0NodeIds *[]int) {
	for idx := range nodes {
		node := &nodes[idx]
		// We set an invalid value here. This will be useful when checking if all nodes received
		// a valid value.
		node.Position.Level = defaultInvalidLevel
		if len(node.IntIdFields.DependsOnIds) == 0 {
			// No dependencies -> level 0 (absolute bottom)
			node.Position.Level = 0
			pushBack(level0NodeIds, idx)
		}
	}
}

// Process a single node in the computeLevels function.
func processNodeComputeLevels(node *NodeData, nodes []NodeData, nextLevelNodeIds *[]int) {
	childLevel := node.Position.Level + 1
	for _, childNodeId := range node.IntIdFields.UsedByIds {
		childNode := &nodes[childNodeId]
		childNode.Position.Level = childLevel
		pushBack(nextLevelNodeIds, childNodeId)
	}
}

// Go over all the current nodes and process them
func processCurrentNodesComputeLevels(nodes []NodeData, currentLevelNodeIds []int,
	nextLevelNodeIds *[]int) {
	for _, nodeId := range currentLevelNodeIds {
		processNodeComputeLevels(&nodes[nodeId], nodes, nextLevelNodeIds)
	}
}

// Validate the result of computeLevels
func validateComputeLevels(nodes []NodeData) error {
	for _, node := range nodes {
		if node.Position.Level == defaultInvalidLevel {
			return fmt.Errorf("unreachable node '%v'", node.InputFields.Name)
		}

		// Every child must be at least 1 level above the current node
		expectedMinLevelForChildren := node.Position.Level + 1
		for _, childNodeId := range node.IntIdFields.UsedByIds {
			childNode := &nodes[childNodeId]
			if childNode.Position.Level < expectedMinLevelForChildren {
				// Shows a bug in the code
				panic(fmt.Sprintf("Child node level %v < expected level %v (child %v, parent %v)",
					childNode.Position.Level, expectedMinLevelForChildren,
					childNode.InputFields.Name, node.InputFields.Name))
			}
		}

		// A nodes level = max(level of all parents) + 1
		maxParentLevel := -1
		for _, parentNodeId := range node.IntIdFields.DependsOnIds {
			parentNode := &nodes[parentNodeId]
			if parentNode.Position.Level > maxParentLevel {
				maxParentLevel = parentNode.Position.Level
			}
		}
		if node.Position.Level != maxParentLevel+1 {
			// Shows a bug in the code
			panic(fmt.Sprintf("For node %v, mismatch in level. Got %v, expected %v",
				node.InputFields.Name, node.Position.Level, maxParentLevel+1))
		}
	}

	return nil
}

// Compute level for all the nodes.
// This is a tricky algorithm.  Basic idea is as follows:
// Find all nodes that do not have any dependencies and assign level 0
// Now, iteratively run for level = 0, 1, ... so on:
//
//	  for all nodes in current level,
//			find nodes that use them
//			assign them current level + 1
//			keep a list of these and use as the starting point of next iteration
//
// Maximum number of iterations = number of nodes.
// At the end, perform sanity checks on the code (and coder)
func computeLevels(nodes []NodeData) error {
	var currentLevelNodeIds []int
	nextLevelNodeIds := make([]int, 0, defaultCapacity)

	initializeForComputeLevels(nodes, &nextLevelNodeIds)
	if len(nextLevelNodeIds) == 0 {
		return fmt.Errorf("found no level 0 nodes")
	}

	// In the worse case, every node get a unique level.
	// We already assigned 0 for the initialization for at least 1 node.
	maxIterationCount := len(nodes) - 1
	// Heart of the algorithm
	for step := 0; step < maxIterationCount; step++ {
		// Check if there is any node to process
		if len(nextLevelNodeIds) == 0 {
			break
		}
		currentLevelNodeIds = nextLevelNodeIds
		nextLevelNodeIds = make([]int, 0, defaultCapacity)

		processCurrentNodesComputeLevels(nodes, currentLevelNodeIds, &nextLevelNodeIds)
		nextLevelNodeIds = getUnique(nextLevelNodeIds)
	}

	err := validateComputeLevels(nodes)
	if err != nil {
		return err
	}

	return nil
}

// Compute shifts - This is straightforward. For every level, we go from left to right.
// We can also compute levelMap with this function.
func computeShiftsAndGetLevelMap(nodes []NodeData) [][]int {
	levelMap := make([][]int, 0)
	if len(nodes) == 0 {
		return levelMap
	}

	// find maxLevel
	maxLevel := -1
	for _, node := range nodes {
		if node.Position.Level > maxLevel {
			maxLevel = node.Position.Level
		}
	}

	if maxLevel < 0 {
		// shows a bug in code
		panic("Unable to find maxLevel")
	}

	// Initialize levelMap for each level
	levelMap = make([][]int, maxLevel+1)
	for level := 0; level <= maxLevel; level++ {
		levelMap[level] = make([]int, 0)
	}

	// Fill levelMap and shift based on each node level
	for idx := range nodes {
		node := &nodes[idx]
		level := node.Position.Level
		node.Position.Shift = len(levelMap[level])
		levelMap[level] = append(levelMap[level], idx)
	}

	// Sanity check:
	for level := 0; level <= maxLevel; level++ {
		if len(levelMap[level]) == 0 {
			panic(fmt.Sprintf("Level %v has 0 nodes", level))
		}
	}

	return levelMap
}

// Used to convert numeric IDs to string IDs used by HTML elements
func formatIntId(id int) string {
	// TODO: do we need this many digits?!
	return fmt.Sprintf("%05d", id)
}

// Build HTML field IDs used by the connection dots
// D connections: D_OWNER_PARTNER
// U connections: U_PARTNER_OWNER
func buildDotElemFields(prefix string, ownerId int, partnerId int) DotElemFields {
	ownerIdStr := formatIntId(ownerId)
	partnerIdStr := formatIntId(partnerId)
	first := ownerIdStr
	second := partnerIdStr
	if prefix == "U" {
		first = partnerIdStr
		second = ownerIdStr
	}
	dotElemId := fmt.Sprintf("%s_%s_%s", prefix, first, second)
	return DotElemFields{dotElemId, partnerIdStr}
}

// Fill HTML element IDs in string form
// The ids are formatted as follows (examples):
// Node => N00001 (N prefix)
// DependsOn Connection Dot: D_N00008_N00005 (dot is carried by the first node N00008.
// It depends on the second node N00005)
// UsedBy Connection Dot: U_N00003_N00002 (dot is carried by the first node N00002.
// It is used by the second node N00003). Note that the order of nodes is reversed for
// the UsedBy connection dot. This makes it easy to map one dot to another when making
// connections. It is always the node on the top that comes first.
func fillElemIds(node *NodeData) {
	nodeElemId := formatIntId(node.IntIdFields.Uid)
	node.ElemFields.NodeElemId = nodeElemId

	node.ElemFields.DependsOnDots = make([]DotElemFields, 0)
	for _, dependsOnId := range node.IntIdFields.DependsOnIds {
		dotElemFields := buildDotElemFields("D", node.IntIdFields.Uid, dependsOnId)
		pushBack(&node.ElemFields.DependsOnDots, dotElemFields)
	}

	node.ElemFields.UsedByDots = make([]DotElemFields, 0)
	for _, usedById := range node.IntIdFields.UsedByIds {
		dotElemFields := buildDotElemFields("U", node.IntIdFields.Uid, usedById)
		pushBack(&node.ElemFields.UsedByDots, dotElemFields)
	}
}

// Fill HTML element IDs for all the nodes
func fillElemIdsForAllNodes(nodes []NodeData) {
	for idx := range nodes {
		fillElemIds(&nodes[idx])
	}
}

// Each node gets a position, which will be set based on inline CSS.
// It is a bit tricky since we want to center the alignment.
func computeNodePositionsAndUpdate(displayConfig *DisplayConfigFields,
	levelMap [][]int, nodes []NodeData) {

	// To be used to calculate max shift and center aligning
	maxNodesPerLevel := 0
	for _, nodeIdsForLevel := range levelMap {
		if len(nodeIdsForLevel) > maxNodesPerLevel {
			maxNodesPerLevel = len(nodeIdsForLevel)
		}
	}

	if maxNodesPerLevel == 0 {
		return
	}

	maxLevel := len(levelMap) - 1
	hscale := displayConfig.HorizontalStepPx
	// We can use levelMap to initialize the positions of nodes
	for level, nodeIdsForLevel := range levelMap {
		for shift, nodeId := range nodeIdsForLevel {
			node := &nodes[nodeId]
			// Horizontal centering shift
			centering := ((maxNodesPerLevel - len(nodeIdsForLevel)) * hscale) / 2
			node.ElemFields.LeftPx = shift*hscale + centering
			node.ElemFields.TopPx = (maxLevel - level) * displayConfig.VerticalStepPx
		}
	}
}

// Do all the steps related to creating list of NodeData and filling all the fields.
// This is the top level function which handles everything.
func createComputeAndFillNodeDataList(gdfData *GdfDataStruct) ([]NodeData, error) {
	nodeDataSeq := createNodeDataList(gdfData)

	err := fillIntIdFields(nodeDataSeq)
	if err != nil {
		return nodeDataSeq, err
	}

	err = computeLevels(nodeDataSeq)
	if err != nil {
		return nodeDataSeq, err
	}

	levelMap := computeShiftsAndGetLevelMap(nodeDataSeq)
	fillElemIdsForAllNodes(nodeDataSeq)
	computeNodePositionsAndUpdate(&gdfData.DisplayConfig, levelMap, nodeDataSeq)

	return nodeDataSeq, nil
}
