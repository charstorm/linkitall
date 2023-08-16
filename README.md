# linkitall

❗The project is still in beta phase. Expect breaking changes to graph definitions in future.

Linkitall is a tool to create a dependency graph of ideas. The generated graph gives a strict level-based graphical structure to knowledge. Linkitall takes a graph definition file in YAML format and generates an HTML file containing nodes and connections between them. This makes Linkitall a powerful tool to build teaching materials as well as reference materials.

This tool is part of the **graphitout** project.  Checkout our [YouTube channel](https://www.youtube.com/channel/UCSYUPhmh-x85NSslUvz-PUQ) to see this tool in action.

Our motto is **untangle, refactor, rename**

## Build

The steps for building the tool from the source is given here.

Clone the source files locally:

```bash
git clone https://github.com/charstorm/linkitall.git
```

The source files are located in `linkitall/src` directory. The project is written in Golang and is tested on version `go1.20.7`.

CD to the `src` directory, build the tool, and check the result:

```bash
cd linkitall/src
go build
./linkitall --help
```

To make the tool available globally, it is required to add the directory path to PATH environment variable.

```bash
export PATH="$PATH:/path/to/linkitall/src/"
```

## Graph Generation

The tool takes the path to a directory containing the graph file `graph.yaml` and it's dependent resources as the input. See `examples/simple` for a quick reference. More details about the graph file is explained in the section "Graph File",

To run the graph generation, execute the following command (assuming `targetdir` as the directory of graph file and its resources):

```bash
linkitall -i targetdir
```

If executed successfully, it will generate  an `index.html` file at `targetdir`. Additionally, asset files (CSS, JS, etc) required for the generated HTML will be copied to the `targetdir` with name `linkitall_assets`. These are initially located in the same directory of the `linkitall` tool. However, it is a one-time action. Subsequent invocation of the tool will skip this copy-assets step.

One can open the generated HTML file in a browser and see the result.

### Server Mode

The default behavior of the tool is to run the generation process only once. This is not ideal for development. For that, we have added a server mode, which can be enabled by the -s flag.

```bash
linkitall -s -i targetdir
```

This will start an HTTP development server at default port 8101. One can see the results by visiting http://127.0.0.1:8101 . 

The tool will wait for user input to update the generated graph. Enter will trigger a graph generation, q will quit the tool.

With this development process for the graph will be as follows:

1. Make changes to graph file

2. Press Enter to trigger a graph generation
   
   1. If there are errors, fix them and continue

3. Refresh the webpage to see the updated graph

## Graph File

To be filled later.

For now see `examples/simple` for reference.

## License

All files in this repo, except those in `src/linkitall_assets/vendor/`, follow the MIT License.
See LICENSE for details. The files in the vendor directory have their own licenses.
