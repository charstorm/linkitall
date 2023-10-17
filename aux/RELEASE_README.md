# Linkitall

This package contains the Linkitall release files (amd64) for Windows and Linux operating systems.
There is no installation required.

Steps for running the tool are given below.

## Windows

1. Extract the archive to a suitable location
2. Execute the batch script `start_server.bat` found inside the extracted directory.
   Just double clicking the script file should be enough in most cases.
3. Provide the path to the directory with graph file when asked (example is provided in the GitHub repository)
4. Allow permission for networking if asked
5. Go to `http://127.0.0.1:8101` in a browser to see the graph generated

The tool will be running in server mode. It follows a simple `update-generate-refresh` cycle.

1. Edit/Update the graph file
2. Press enter in the script window to generate the new HTML for the graph
3. Refresh the web page
4. When development is complete, press `q` to quit the tool!

## Linux

Use the following command to run the tool:

```bash
linkitall --overwrite --release --serve --listen ":8101" -i path/to/graph/dir
```

See the main `README.md` of the linitall repo to understand the meaning of these flags.
For the usage see the `update-generate-refresh` part explained in the Windows section above.

## License

[The license for the project](https://github.com/charstorm/linkitall/blob/main/LICENSE)

External packages:

* You can find the list of external Go packages [here](https://github.com/charstorm/linkitall/blob/main/src/go.mod).
  Their licenses can be found in their respective repositories.
* The project uses the Javascript package `leader-line` to draw all the connections.
  License for it can be found [here](https://anseki.github.io/leader-line/).

