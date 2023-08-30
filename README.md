# linkitall

❗The project is still in beta phase. Expect breaking changes to graph definitions in the future.

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

If executed successfully, it will generate  an `index.html` file at `targetdir`. Additionally, asset files (CSS, JS, etc) required for the generated HTML will be copied to the `targetdir` with name `linkitall_assets`. These are initially located in the same directory of the `linkitall` tool. However, it is a one-time action. Subsequent invocation of the tool will skip this copy-assets step. Same is true for the vendor files (3rd party libraries) used by the project. They are stored in `linkitall_vendor` directory.

One can open the generated HTML file in a browser and see the result.

### CLI

```
Usage: linkitall [--serve] [--release] [--listen LISTEN] --indir INDIR [--graph GRAPH] [--out OUT] [--overwrite]

Options:
  --serve, -s            run in edit-update-serve mode
  --release, -r          run in release mode
  --listen LISTEN, -l LISTEN
                         listen address in serve mode [default: :8101]
  --indir INDIR, -i INDIR
                         path to the input directory
  --graph GRAPH, -g GRAPH
                         input graph base filename [default: graph.yaml]
  --out OUT, -o OUT      output html base filename [default: index.html]
  --overwrite            overwrite asset files
  --help, -h             display this help and exit
```

1. `serve` - to run in server mode. See below.
2. `release` - run in release mode. This includes:
    - use CDN for links, instead of local vendor files.
3. `listen` - the address to listen to (eg: ":8101") in the server mode.
4. `indir` - input (or target) directory containing the graph file.
5. `graph` - base-name of the graph file (eg: "main.yaml") inside `indir`.
6. `out` - base-name of the output file to be created inside `indir`.


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

All files in this repo (except those in `src/linkitall_vendor/`) follow the MIT License.
See LICENSE for details. The files in the vendor directory have their own licenses.
