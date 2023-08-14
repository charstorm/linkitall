# linkitall
❗The project is still in beta phase. Expect breaking changes to graph definitions in future.

Linkitall is a tool to create a dependency graph of ideas.
The generated graph gives a strict level-based graphical structure to knowledge.
Linkitall takes a graph definition file in YAML format and generates an HTML file
containing nodes and connections between them.
This makes Linkitall a powerful tool to build teaching materials as well as
reference materials.

This tool is part of the **graphitout** project. 
Checkout our [YouTube channel](https://www.youtube.com/channel/UCSYUPhmh-x85NSslUvz-PUQ) to see
this tool in action.

Our motto is **untangle, refactor, rename**

## Building
For building and installing binary releases, refer to the INSTALL files inside the release archives.
The steps for building the tool from the source files for Linux environment is given here.

Clone the source files locally:
```bash
git clone https://github.com/charstorm/linkitall.git
```
The source files are located in `linkitall/src` directory. The project is written in Golang and
is tested on version `go1.20.7`.
CD to the `src` directory, build the tool, and check the result:
```bash
cd linkitall/src
go build
./linkitall --help
```

To make the tool available globally, it is required to add the path to PATH environment
variable.
```bash
export PATH="$PATH:/path/to/linkitall/src/"
```

## License
All files in this repo, except those in `src/linkitall_assets/vendor/`, follow the MIT License.
See LICENSE for details. The files in the vendor directory have their own licenses.
