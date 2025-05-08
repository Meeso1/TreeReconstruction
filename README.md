# Tree Reconstruction CLI

A minimal CLI application template written in Go.

## Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/treereconstruction.git
cd treereconstruction

# Build the application
go build -o bin/treereconstruction

# Run the application
./bin/treereconstruction
```

## Usage

```bash
# Show help
./bin/treereconstruction --help

# Enable verbose mode
./bin/treereconstruction -v

# Use a custom config file
./bin/treereconstruction -c config.yaml

# Check version
./bin/treereconstruction version
```

## Development

```bash
# Get dependencies
go mod download

# Run tests
go test ./...
```

## Building with Version Information

To build with version information:

```bash
go build -ldflags="-X 'treereconstruction/cmd.Version=v1.0.0' -X 'treereconstruction/cmd.Commit=$(git rev-parse HEAD)' -X 'treereconstruction/cmd.BuildDate=$(date)'" -o bin/treereconstruction
``` 