# tq - TOON Query Processor

`tq` is a command-line tool for querying and transforming [TOON](https://github.com/toon-format/toon) (Token-Oriented Object Notation) files, similar to how `jq` works for JSON.

## What is TOON?

TOON is a compact, human-readable format for serializing JSON data, optimized for Large Language Model prompts. It uses 30-60% fewer tokens than formatted JSON while remaining fully lossless and human-readable.

## Features

- **Element Extraction**: Extract fields and array elements using familiar `.key` and `.[index]` syntax
- **Filtering**: Filter data with `select()` conditions
- **Output Formatting**: Output as TOON (default) or JSON
- **Data Transformation**: Transform data using `map()`, aggregations, and more
- **jq-Compatible**: Uses jq's powerful query syntax under the hood

## Installation

### Prerequisites

- **Node.js** 18 or later ([download](https://nodejs.org/))
- **jq** 1.6 or later ([installation guide](https://stedolan.github.io/jq/download/))
  - macOS: `brew install jq`
  - Ubuntu/Debian: `sudo apt-get install jq`
- **Go** 1.21 or later (for building from source) ([download](https://golang.org/dl/))

### Quick Install (Recommended)

```bash
git clone https://github.com/RHEMS-japan/tq.git
cd tq
./install.sh
```

This will:
1. Check that all prerequisites are installed
2. Install Node.js dependencies
3. Build the `tq` binary
4. Install to `~/.local/bin/tq`
5. Copy required scripts to `~/.tq/scripts`

**Add to PATH** (if not already):

```bash
# For zsh (macOS default)
echo 'export PATH="$PATH:$HOME/.local/bin"' >> ~/.zshrc
source ~/.zshrc

# For bash
echo 'export PATH="$PATH:$HOME/.local/bin"' >> ~/.bashrc
source ~/.bashrc
```

**Verify installation**:
```bash
tq --version  # Should output: tq version 0.1.0
```

### Manual Installation

```bash
# Clone and build
git clone https://github.com/RHEMS-japan/tq.git
cd tq
npm install
go build -o tq ./cmd/tq

# Install manually
mkdir -p ~/.local/bin ~/.tq/scripts
cp tq ~/.local/bin/
cp scripts/*.js ~/.tq/scripts/
cp -r node_modules ~/.tq/

# Add to PATH if needed
echo 'export PATH="$PATH:$HOME/.local/bin"' >> ~/.bashrc  # or ~/.zshrc
```

## Quick Start

```bash
# Pretty print a TOON file
tq '.' data.toon

# Extract a field
tq '.name' data.toon

# Get an array element
tq '.users[0]' data.toon

# Filter an array
tq '.users[] | select(.age > 25)' data.toon

# Transform data
tq '.users | map({name, email})' data.toon

# Output as JSON
tq --json '.' data.toon
```

## Command-Line Usage

```
tq [options] [filter] [file]

Options:
  --json         Output as JSON instead of TOON
  -h, --help     Show help message
  -v, --version  Show version

Filter:
  jq-compatible filter expression (default: ".")

Input:
  File path or stdin
```

## Examples

### 1. Element Extraction

```bash
# Extract a specific field
$ tq '.company' data.toon
Acme Corp

# Navigate nested objects
$ tq '.address.city' data.toon
San Francisco

# Get first array element
$ tq '.employees[0]' data.toon
id: 1
name: Alice Smith
role: Engineer
salary: 95000
active: true

# Extract field from all array elements
$ tq '.employees[].name' data.toon
Alice Smith
---
Bob Johnson
---
Charlie Brown
```

### 2. Filtering

```bash
# Filter by condition
$ tq '.employees[] | select(.salary > 90000)' data.toon
id: 1
name: Alice Smith
role: Engineer
salary: 95000
active: true
---
id: 3
name: Charlie Brown
role: Manager
salary: 110000
active: true

# Multiple conditions
$ tq '.employees[] | select(.active == true and .role == "Engineer")' data.toon
id: 1
name: Alice Smith
role: Engineer
salary: 95000
active: true
```

### 3. Output Formatting

```bash
# Default: TOON format
$ tq '.' data.toon
name: John Doe
age: 30
email: john@example.com

# JSON format
$ tq --json '.' data.toon
{"name":"John Doe","age":30,"email":"john@example.com"}
```

### 4. Data Transformation

```bash
# Extract specific fields
$ tq '.users | map({name, email})' data.toon
[3]{name,email}:
  Alice,alice@example.com
  Bob,bob@example.com
  Charlie,charlie@example.com

# Calculate new values
$ tq '.employees | map({name, newSalary: (.salary * 1.1)})' data.toon
[5]{name,newSalary}:
  Alice Smith,104500
  Bob Johnson,93500
  Charlie Brown,121000
  Diana Prince,107800
  Eve Wilson,95700

# Sort and transform
$ tq '.employees | sort_by(.salary) | reverse | .[0]' data.toon
id: 3
name: Charlie Brown
role: Manager
salary: 110000
active: true
```

See [EXAMPLES.md](EXAMPLES.md) for more comprehensive examples.

## How It Works

`tq` works by converting TOON to JSON, applying jq filters, then converting back to TOON:

```
TOON → JSON → jq filter → JSON → TOON
```

This approach leverages:
- The official [@toon-format/toon](https://www.npmjs.com/package/@toon-format/toon) TypeScript library for parsing
- [jq](https://stedolan.github.io/jq/) for powerful querying
- Go for fast, portable execution

## Supported jq Features

Since `tq` uses `jq` internally, it supports most jq features:

- **Basic filters**: `.`, `.key`, `.[index]`, `.[]`
- **Operators**: `|`, `,`, `+`, `-`, `*`, `/`, `%`
- **Conditionals**: `if-then-else`, `select()`
- **Functions**: `map()`, `select()`, `sort_by()`, `group_by()`, `unique()`, etc.
- **String interpolation**: `\(expr)`
- **Array slicing**: `.[start:end]`
- **Object construction**: `{key: value}`

For complete jq syntax, see the [jq manual](https://stedolan.github.io/jq/manual/).

## Development

```bash
# Install dependencies
npm install
go mod download

# Build
go build -o tq ./cmd/tq

# Run tests
go test ./...

# Run with sample data
./tq '.' testdata/sample.toon
```

## Project Structure

```
tq/
├── cmd/tq/              # Main application
│   └── main.go
├── scripts/             # Node.js helper scripts
│   ├── toon-to-json.js  # TOON → JSON converter
│   └── json-to-toon.js  # JSON → TOON converter
├── testdata/            # Sample TOON files
│   ├── sample.toon
│   ├── users.toon
│   ├── company.toon
│   └── products.toon
├── EXAMPLES.md          # Comprehensive examples
└── README.md            # This file
```

## Uninstall

To completely remove `tq`:

```bash
rm -f ~/.local/bin/tq
rm -rf ~/.tq
```

If you added the PATH export to your shell config, remove this line:
```bash
export PATH="$PATH:$HOME/.local/bin"
```

## Limitations

- Requires both Node.js (for TOON parsing) and jq (for querying) to be installed
- May be slower than pure implementations due to multiple conversions
- Some very advanced jq features may not work perfectly with TOON's structure

## Future Plans

- [ ] Native Go TOON parser (remove Node.js dependency)
- [ ] Performance optimizations
- [ ] Additional output formats
- [ ] Streaming support for large files
- [ ] Syntax highlighting in output

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Related Projects

- [TOON Format](https://github.com/toon-format/toon) - Official TOON format specification and TypeScript implementation
- [jq](https://stedolan.github.io/jq/) - Command-line JSON processor
- [TOON Specification](https://github.com/toon-format/spec) - Detailed TOON format specification

## Acknowledgments

- Built on top of the excellent [TOON format](https://github.com/toon-format/toon) by the toon-format team
- Powered by [jq](https://stedolan.github.io/jq/) for query processing
- Inspired by the simplicity and power of jq

---

Made with ❤️ by [RHEMS-japan](https://github.com/RHEMS-japan)
