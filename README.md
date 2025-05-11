# Tree Reconstruction CLI

A CLI application for tree reconstruction written in Go.

## Installation

```bash
# Clone the repository
git clone https://github.com/Meeso1/treereconstruction.git
cd treereconstruction

# Build the application
make build

# Run the application
./bin/tree-reconstruction
```

## Usage

```bash
# Show help
./bin/tree-reconstruction help

# Check version
./bin/tree-reconstruction version

# Reconstruct a tree from distance matrix in input_file.txt
./bin/tree-reconstruction reconstruct -i input_file.txt
```

## Development

```bash
# Get dependencies
go mod download

# Run tests
make test
```

## Building with Version Information

To build with version information:

```bash
go build -ldflags="-X 'treereconstruction/cmd.Version=v1.0.0' -X 'treereconstruction/cmd.Commit=$(git rev-parse HEAD)' -X 'treereconstruction/cmd.BuildDate=$(date)'" -o bin/treereconstruction
``` 