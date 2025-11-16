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
tq --version  # Should output: tq version 0.2.0
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

Since `tq` uses `jq` internally, it supports most jq features including:

### Basic Filters
- `.` - Identity
- `.key` - Field access
- `.[index]` - Array indexing
- `.[]` - Array/object iteration
- `.[start:end]` - Array slicing

### Operators

#### Arithmetic Operators
- `+` - Addition (numbers), concatenation (strings/arrays), merge (objects)
- `-` - Subtraction (numbers), array difference
- `*` - Multiplication (numbers), string repetition
- `/` - Division
- `%` - Modulo (remainder)

```bash
# Number operations
tq '.price + 10' data.toon          # Add 10 to price
tq '.price * 1.1' data.toon         # 10% increase

# String concatenation
tq '.first + " " + .last' data.toon # "John Doe"

# Array concatenation
tq '[1,2] + [3,4]'                  # [1,2,3,4]

# Object merge
tq '{a:1} + {b:2}'                  # {a:1, b:2}
```

#### Comparison Operators
- `==` - Equal
- `!=` - Not equal
- `<` - Less than
- `<=` - Less than or equal
- `>` - Greater than
- `>=` - Greater than or equal

```bash
# Numeric comparison
tq '.age > 25' data.toon            # true/false
tq '.price >= 100' data.toon

# String comparison
tq '.status == "active"' data.toon

# Use with select()
tq '.users[] | select(.age > 25)' data.toon
```

#### Logical Operators
- `and` - Logical AND
- `or` - Logical OR
- `not` - Logical NOT

```bash
# Combine conditions
tq '.users[] | select(.age > 25 and .active)' data.toon
tq '.users[] | select(.role == "admin" or .role == "owner")' data.toon
tq '.active | not' data.toon        # Negate boolean
```

#### Special Operators
- `|` - Pipe (chain operations)
- `//` - Alternative operator (returns right side if left is null/false)

```bash
# Alternative operator for defaults
tq '.optional // "default"' data.toon
tq '.a // .b // "none"' data.toon   # Chain alternatives
```

### Built-in Functions

#### Array Functions
- `length` - Get array/object/string length
- `reverse` - Reverse array
- `sort` - Sort array (ascending)
- `sort_by(expr)` - Sort by expression
- `unique` - Remove duplicates
- `group_by(expr)` - Group by expression
- `add` - Sum numbers or concatenate strings/arrays
- `min`, `max` - Minimum/maximum value
- `first`, `last` - First/last element
- `flatten` - Flatten nested arrays
- `map(expr)` - Transform each element

#### Object Functions
- `keys` - Get object keys (sorted)
- `values` - Get object values
- `has(key)` - Check if key exists
- `in(object)` - Check if key is in object
- `to_entries` - Convert object to key-value pairs
- `from_entries` - Convert key-value pairs to object
- `with_entries(expr)` - Transform entries

#### Type Functions
- `type` - Get type name
- `tonumber` - Convert to number
- `tostring` - Convert to string

#### String Functions

**Basic String Operations**
- `startswith(str)` - Check prefix
- `endswith(str)` - Check suffix
- `contains(str)` - Check substring
- `split(sep)` - Split string into array
- `join(sep)` - Join array elements into string

**Case Conversion**
- `ascii_upcase` - Convert to uppercase
- `ascii_downcase` - Convert to lowercase

**Trimming**
- `ltrimstr(str)` - Remove prefix string
- `rtrimstr(str)` - Remove suffix string

**Regular Expressions**
- `test(regex)` - Test if string matches regex (returns boolean)
- `match(regex)` - Match string against regex (returns match object)
- `sub(regex; replacement)` - Replace first match
- `gsub(regex; replacement)` - Replace all matches

**Format Functions**
- `@base64` - Base64 encode
- `@uri` - URL encode
- `@csv` - CSV format
- `@json` - JSON encode
- `@html` - HTML encode
- `@base64d` - Base64 decode

**String Interpolation**
```bash
# Embed expressions in strings
tq '"\(.name) is \(.age) years old"' data.toon
# Result: "John is 30 years old"
```

### Conditionals & Logic

#### if-then-else Expressions
Conditional logic for complex decision making:

```bash
# Basic conditional
tq 'if .age > 18 then "adult" else "minor" end' data.toon

# With elif (else if)
tq 'if .score >= 90 then "A" elif .score >= 80 then "B" else "C" end' data.toon

# Nested conditionals
tq 'if .active then (if .premium then "premium" else "standard" end) else "inactive" end' data.toon

# In object construction
tq '{name, status: (if .active then "Active" else "Inactive" end)}' data.toon

# In string interpolation
tq '"\(.name): \(if .verified then "✓" else "✗" end)"' data.toon

# With logical operators
tq 'if .age > 18 and .verified then "allowed" else "denied" end' data.toon
```

#### Error Handling
Robust error handling for missing fields and invalid operations:

```bash
# Optional access operator (?) - returns null instead of error
tq '.user.email?' data.toon                    # Returns null if missing
tq '.[99]?' data.toon                          # Returns null if out of bounds

# try-catch - catch errors and provide fallback
tq 'try .invalid catch "default"' data.toon    # Returns "default" on error
tq 'try (1/0) catch "error"' data.toon         # Catches division by zero

# Alternative operator (//) - use default if null/false
tq '.field // "default"' data.toon             # Use default if null/false
tq '.a // .b // "none"' data.toon              # Chain alternatives

# Combining error handling
tq '.user?.email? // "no email"' data.toon     # Optional + alternative
tq '{name, email: (.email // "N/A")}' data.toon  # In object construction
```

#### Other Conditional Tools
- `select(expr)` - Filter by condition (keep only matching items)
- `empty` - Return nothing (useful with conditionals for filtering)

### Variables

Bind values to variables for reuse in complex queries using the `as $var` syntax:

```bash
# Basic variable binding
tq '.age as $a | {name, age: $a, next_year: ($a + 1)}' data.toon

# Multiple variables
tq '.price as $p | .quantity as $q | {total: ($p * $q)}' data.toon

# Variables in select
tq '.users[] | .age as $a | select($a > 25) | {name, age: $a}' data.toon

# Variables across nested iterations
tq '.users[] | .name as $n | .scores[] | {user: $n, score: .}' data.toon

# Variables with aggregations
tq '[.prices[]] | add as $total | . | map(. / $total)' data.toon
# Calculate percentage of each price relative to total

# Chained variable assignments
tq '.salary as $s | ($s * 0.1) as $bonus | {salary: $s, bonus: $bonus, total: ($s + $bonus)}' data.toon
```

**Common patterns:**
- Store intermediate results: `.field as $var | ... use $var multiple times ...`
- Preserve context in iterations: `.name as $n | .items[] | {parent: $n, item: .}`
- Simplify complex expressions: `(.a + .b) as $sum | .c as $other | {sum: $sum, ratio: ($sum / $other)}`

### Construction

#### Object Construction
Create new objects by selecting or transforming fields:

```bash
# Field shorthand - select specific fields
tq '{name, age}' data.toon
# Result: {"name": "John", "age": 30}

# Rename fields
tq '{n: .name, a: .age}' data.toon
# Result: {"n": "John", "a": 30}

# Computed fields
tq '{name, doubled: (.age * 2)}' data.toon
# Result: {"name": "John", "doubled": 60}

# Nested objects
tq '{user: {name, age}, email}' data.toon
# Result: {"user": {"name": "John", "age": 30}, "email": "..."}

# With map() - transform arrays
tq '.users | map({name, email})' data.toon
# Extract only name and email from each user
```

#### Array Construction
Create new arrays from expressions:

```bash
# Simple array
tq '[.a, .b, .c]' data.toon
# Result: [1, 2, 3]

# With computation
tq '[.x, (.x * 2), (.x * 3)]' data.toon
# Result: [5, 10, 15]

# From iteration
tq '[.users[].name]' data.toon
# Result: ["Alice", "Bob", "Charlie"]

# With filtering
tq '[.users[] | select(.age > 25)]' data.toon
# Result: array of users over 25

# Range function
tq '[range(5)]'
# Result: [0, 1, 2, 3, 4]
```

#### String Interpolation
Build strings with embedded expressions:

```bash
tq '"\(.name) is \(.age) years old"' data.toon
# Result: "John is 30 years old"
```

For complete jq syntax, see the [jq manual](https://stedolan.github.io/jq/manual/).

### Quick Reference

```bash
# Operators
tq '.price * 1.1'                      # Arithmetic (10% increase)
tq '.first + " " + .last'              # String concatenation
tq '.users[] | select(.age > 25)'      # Comparison in filter
tq '.a // .b // "default"'             # Alternative operator

# Array operations
tq '.items | length'                    # Count items
tq '.items | sort_by(.price)'          # Sort by price
tq '.items | map(.name)'               # Extract names
tq '.items | unique'                   # Remove duplicates
tq '.prices | add'                     # Sum prices

# Object operations
tq '. | keys'                          # List keys
tq '. | has("field")'                  # Check field exists
tq '. | to_entries'                    # Convert to array

# Type checking
tq '.value | type'                     # Get type
tq '.age | tonumber'                   # Convert to number

# String operations
tq '.email | split("@")'               # Split email
tq '.tags | join(", ")'                # Join with comma
tq '.name | ascii_upcase'              # Convert to uppercase
tq '.url | ltrimstr("https://")'       # Remove prefix
tq '.text | gsub("foo"; "bar")'        # Replace all occurrences
tq '"\(.name) <\(.email)>"'            # String interpolation

# Conditionals
tq 'if .age > 18 then "adult" else "minor" end'  # Basic conditional
tq '{name, tier: (if .premium then "Premium" else "Free" end)}'  # In objects

# Error handling
tq '.user.email? // "no email"'        # Optional access with default
tq 'try .field catch "default"'        # Catch errors

# Variables
tq '.age as $a | {name, age: $a, next: ($a + 1)}'  # Bind and reuse
tq '.price as $p | .qty as $q | {total: ($p * $q)}'  # Multiple variables

# Construction
tq '{name, age}'                       # Select fields
tq '.users | map({name, email})'       # Transform array
tq '[.users[].name]'                   # Build array from field
tq '[range(10)]'                       # Generate range
```

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
