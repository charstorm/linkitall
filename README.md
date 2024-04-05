## Introduction

Linkitall offers the capability to generate a comprehensive visual representation illustrating the interrelation between ideas. This graphical depiction presents a structured hierarchy and connections among various concepts through a dependency graph format. Utilizing Linkitall entails the creation of a YAML file specifying the graph's structure, following which the tool seamlessly produces an HTML file featuring interactive nodes and links. Leveraging its functionalities, Linkitall serves as an invaluable resource for developing educational and reference materials, providing a holistic overview of diverse topics.

Moreover, Linkitall encourages a deeper understanding of concepts by facilitating the exploration of dependencies. As we question "why" or "how," we inherently traverse these dependencies, uncovering the intricate relationships that underpin our knowledge. This iterative process not only enhances comprehension but also fosters critical thinking and analysis.

An example: [Trigonometric Relations](https://charstorm.github.io/class-11-12-india/class11/maths/trigonometry/relations/)

## Usage

Download the latest release zipfile and follow the instructions in the README.md inside.

If the release files fail due to some reason, please build the tool from the source as explained below.

## Build

The steps for building the tool from the source are given here.

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

To make the tool available globally, it is required to add the directory path to the PATH environment variable.

```bash
export PATH="$PATH:/path/to/linkitall/src/"
```

## Graph Generation

The tool takes the path to a directory containing the graph file `graph.yaml` and its dependent resources as the input. See `examples/simple` for a quick reference.
More details about the graph file are explained in the section "Graph File",

To run the graph generation, execute the following command (assuming `targetdir` is the directory of graph file and its resources):

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

1. Make changes to the graph file

2. Press Enter to trigger a graph generation

   1. If there are errors, fix them and continue

3. Refresh the webpage to see the updated graph

## Graph File

The Graph Definition File (GDF) is a YAML file with different sections.
The default filename expected for the file is `graph.yaml`.
Different sections of the graph file are explained below.

### head-config
These fields will be forwarded to the `head` section of the output html file.
Example:
```yaml
head-config:
    title: Awesome Graph
    description: A long description
    author: Someone
```

### display-config
These fields control the size and spacing of nodes in the graph.
Example:
```yaml
display-config:
    # Horizontal spacing
    horizontal-step-px: 400
    # Vertical spacing
    vertical-step-px: 300
    # Width of each node
    node-box-width-px: 300
    # There is no configuration for the node height from the graph file.
```
The numerical values shown in the configuration above are the default values.

### resources
These are the resources used in the graph.
Resources can be images, pdf-files, html pages, etc.
When clicked on that node, the corresponding link will be opened.
Each node can "linkto" a resource.
Multiple nodes can share the same resource.
Example:
```yaml
# The keys will be used in the "linkto" field of each node
resources:
    # An internal link
    main: resources/main.html
    # An external link
    disinfectants_wiki: https://en.wikipedia.org/wiki/Disinfectant
    # Local or web images work too
    algae_image: resources/pexels-daria-klet-8479585.jpg
    # Local or web pdf references also work.
    # (Note: #view=fit is not part of the filename. It is added to resize the pdf view)
    microorganisms_pdf: resources/microorganisms.pdf#view=fit
```
An explanation of using these references will be provided below.

### nodes
This is a list of data dictionaries for each node in the graph.
Example:
```yaml
nodes:
      # Expects lowercase, without space, must be unique
    - name: tap_water
      # This the the text that will be shown on the node
      title: Tap Water
      # Text that will be shown below title
      subtitle: Node about tap water
      # Resource information for this node
      linkto:
          # Resource defined in section resources
          resource: main
          # The id of the target section in the resource page (optional)
          target: tap-water
      depends-on:
          # List of dependencies (their name, not title)
          - pure_water
          - impurities

    - name: pure_water
      # If title is not provided, it will be guessed based on name
      linkto:
          resource: main
          target: pure-water
```

### algo-config
These fields control the node placement, direction, etc of the graph generation
algorithm. Example:
```yaml
algo-config:
    # Supported: bottom2top (default), top2bottom
    level-strategy: top2bottom
    # Supported: child2parent (default), parent2child
    arrow-direction: child2parent
    # Supported: ascend (default), descend
    node-sorting: ascend
```

[More details on algo-config](docs/algo-config/README.md)

## User Interface

The graph generated is a basic HTML web-page. But we have added a few features to make it
easy to use.

1. A resource connected to a node (if available) can be accessed by left-clicking the 
   node's title. Middle-click or Ctrl-click will open the same resource in a new tab.
2. Clicking on the connection box (those small circles with a +) will move the view to
   the target node. The target node will be highlighted in this case.
3. When an image resource is viewed in the same page of the graph, keys "[" and "]"
   can be used to control the size/zoom of the image.

## External Examples

We are currently building graphs for topics covered in class 11 and 12 (plus-one and plus-two).
You will find the example graphs in this repository:
[class-11-12-india](https://github.com/charstorm/class-11-12-india) (See readme.md).

## Status

❗The project is still in beta phase.

## License

All files in this repo (except those in `src/linkitall_vendor/`) follow the MIT License.
See LICENSE for details. The files in the vendor directory have their own licenses.
