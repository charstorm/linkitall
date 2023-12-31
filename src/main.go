// This file handles the conversion of Graph in YAML format to HTML

// Important uncommon shortforms used:
// GDF - Graph Definition File (usually in YAML)

// TODO: Add an output directory option for this tool. As of now, input directory is the target
// directory where we will keep all the files. When used in serve mode, we serve files from the
// input directory.

package main

import (
	// "html/template"
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	argparse "github.com/alexflint/go-arg"
	copylib "github.com/otiai10/copy"
)

// Some fields (GraphFile, OutFile) are basepaths (just the filename without dir).
// For these, full path is attached by the getInputsForProcessing() function.
type CliArgs struct {
	ServerMode bool   `arg:"-s,--serve" help:"run in edit-update-serve mode"`
	Release    bool   `arg:"-r,--release" help:"run in release mode"`
	ServerAddr string `arg:"-l,--listen" default:":8101" help:"listen address in serve mode"`
	InputDir   string `arg:"-i,--indir,required" help:"path to the input directory"`
	GraphFile  string `arg:"-g,--graph" default:"graph.yaml" help:"input graph base filename"`
	OutFile    string `arg:"-o,--out" default:"index.html" help:"output html base filename"`
	Overwrite  bool   `arg:"--overwrite" help:"overwrite asset files"`
}

var bufferedStdin *bufio.Reader = bufio.NewReader(os.Stdin)

// Return path to the assets directory inside indir
func getPathToAssetDir(indir string) string {
	return filepath.Join(indir, "linkitall_assets")
}

// Similar to the above, but for vendor directory
func getPathToVendorDir(indir string) string {
	return filepath.Join(indir, "linkitall_vendor")
}

// Check if the give `path` is accessible.
// kind can be "file" or "dir"
func isPathAccessible(path string, kind string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}

	result := false
	if kind == "dir" {
		result = stat.IsDir()
	} else if kind == "file" {
		file, err := os.Open(path)
		if err == nil {
			result = true
			file.Close()
		}
	}
	return result
}

// Check if file can be written
func canFileWrite(path string) bool {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return false
	}
	file.Close()
	return true
}

// Parse arguments and perform steps to prepare input for processing.
// If --indir is specified "?", get the input path from the user via stdin.
// Final InputDir path is converted to absolute path.
// Check for existence of indir and graph file.
func getInputsForProcessing() (CliArgs, error) {
	var args CliArgs
	if len(os.Args) == 1 {
		fmt.Printf("No args. Use --help\n")
		os.Exit(1)
	}
	argparse.MustParse(&args)

	if args.InputDir == "?" {
		fmt.Printf("Enter input directory => ")
		line, err := bufferedStdin.ReadString('\n')
		if err != nil {
			return args, err
		}
		line = strings.TrimSpace(line)
		args.InputDir = line
	}

	if !isPathAccessible(args.InputDir, "dir") {
		return args, fmt.Errorf("input dir not accessible: %s", args.InputDir)
	}

	absInputDir, err := filepath.Abs(args.InputDir)
	if err != nil {
		return args, err
	}

	args.InputDir = absInputDir
	args.GraphFile = filepath.Join(args.InputDir, args.GraphFile)
	if !isPathAccessible(args.GraphFile, "file") {
		return args, fmt.Errorf("unable to find graph file: %s", args.GraphFile)
	}
	// Fill full path to input and output
	args.OutFile = filepath.Join(args.InputDir, args.OutFile)
	if !canFileWrite(args.OutFile) {
		return args, fmt.Errorf("unable to open file for writing: %s", args.OutFile)
	}

	return args, nil
}

// Return path to the parent dir where the executable is.
func getExecutableDir() (string, error) {
	execFile, err := os.Executable()
	if err != nil {
		return "", err
	}

	execDir := filepath.Dir(execFile)
	absExecDir, err := filepath.Abs(execDir)
	if err != nil {
		return "", err
	}

	return absExecDir, nil
}

// Copy all the assets files to the target directory where the output will be generated.
// The asset files (source) are located in the same directory of the executable.
func copyAssetsAndVendorFilesToDir(targetDir string, overwrite bool, release bool) error {
	execDir, err := getExecutableDir()
	if err != nil {
		return err
	}

	assetPath := getPathToAssetDir(execDir)
	targetAssetPath := getPathToAssetDir(targetDir)
	if !overwrite && isPathAccessible(targetAssetPath, "dir") {
		log.Printf("Asset dir %s already exists. Skipping copying assets\n", targetAssetPath)
	} else {
		log.Printf("Copy %s -> %s\n", assetPath, targetAssetPath)
		copylib.Copy(assetPath, targetAssetPath)
	}

	vendorPath := getPathToVendorDir(execDir)
	targetVendorPath := getPathToVendorDir(targetDir)
	if release {
		log.Printf("In release mode. Not copying vendor dir")
	} else if !overwrite && isPathAccessible(targetVendorPath, "dir") {
		log.Printf("Vendor dir %s already exists. Skipping copying vendor\n", targetVendorPath)
	} else {
		log.Printf("Copy %s -> %s\n", vendorPath, targetVendorPath)
		copylib.Copy(vendorPath, targetVendorPath)
	}

	return nil
}

// ** This is the core function which does all the processing **
// Process Graph Data File (GDF) and writes the HTML output.
// The `indir` is also the target dir. Output is generated at the same location.
// Copy the required asset dir to the `indir` before calling this function.
func processGraphWriteOutput(args *CliArgs) error {
	log.Printf("Reading graph: %s\n", args.GraphFile)
	gdfData, readable, err := loadGdf(args.GraphFile)
	if !readable {
		log.Fatalf("graph file %s not readable: %s\n", args.GraphFile, err)
	}

	if err != nil {
		return err
	}

	log.Printf("Preparing nodes\n")
	nodes, err := createComputeAndFillNodeDataList(gdfData)
	if err != nil {
		return err
	}
	log.Printf("Number of nodes: %d\n", len(nodes))

	log.Printf("Generating template data\n")
	controlConfig := ControlConfigFields{
		Release: args.Release,
	}
	templateData := newTemplateData(gdfData, nodes, controlConfig)

	targetAssetDir := getPathToAssetDir(args.InputDir)
	templateFile := filepath.Join(targetAssetDir, "template.html")

	log.Printf("Filling template and writing output\n")
	err = fillTemplateWriteOutput(templateFile, templateData, args.OutFile)
	if err != nil {
		return err
	}

	log.Printf("Done\n")
	return nil
}

// Process the graph file. Print error if any.
func processAndLogError(args *CliArgs) {
	err := processGraphWriteOutput(args)

	if err != nil {
		log.Printf("Error: %s", err)
	}
}

// In server mode, we run a http server on the target directory.
// We also run a read-update cycle to update the output file.
func runInServerMode(args *CliArgs) {
	// Run processing once before starting server
	processAndLogError(args)

	// Start server on the target dir
	fileServer := http.FileServer(http.Dir(args.InputDir))
	go func() {
		log.Printf("Starting server for dir %s. Listening at %s\n",
			args.InputDir, args.ServerAddr)
		http.ListenAndServe(args.ServerAddr, fileServer)
	}()

	time.Sleep(time.Second)

	// Run read-update cycle
	for {
		fmt.Printf("\nq: quit, enter: update output => ")
		line, err := bufferedStdin.ReadString('\n')
		if err != nil {
			log.Fatalf("unable to read from stdin. %s", err)
		}

		line = strings.TrimSpace(line)
		if line == "q" {
			break
		} else if len(line) > 0 {
			log.Printf("Warning: Ignoring input: '%s'", line)
			continue
		}
		processAndLogError(args)
	}
}

func main() {
	args, err := getInputsForProcessing()
	if err != nil {
		log.Fatalf("unable to read args. %s", err)
	}

	// Check if the graph file exists. It is the input file. User only specifies the dir.

	// It doesn't matter whether we are running in server mode or not. We always copy the asset
	// files to the target dir (input dir in this case).
	err = copyAssetsAndVendorFilesToDir(args.InputDir, args.Overwrite, args.Release)
	if err != nil {
		log.Fatalf("unable to copy asset files to %s", args.InputDir)
	}

	if args.ServerMode {
		runInServerMode(&args)
	} else {
		err = processGraphWriteOutput(&args)
	}

	if err != nil {
		log.Fatalf("error while processing %s", err)
	}
}
