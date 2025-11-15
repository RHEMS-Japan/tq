package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const version = "0.1.0"

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "-v" || os.Args[1] == "--version") {
		fmt.Printf("tq version %s\n", version)
		os.Exit(0)
	}

	if len(os.Args) > 1 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
		printHelp()
		os.Exit(0)
	}

	// Parse command-line flags
	filter := "."
	outputFormat := "toon" // toon, json, compact, raw
	colorOutput := false
	var inputFile string
	filterSet := false

	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "--json" {
			outputFormat = "json"
		} else if arg == "--compact" || arg == "-c" {
			outputFormat = "compact"
		} else if arg == "--raw" || arg == "-r" {
			outputFormat = "raw"
		} else if arg == "--color" || arg == "-C" {
			colorOutput = true
		} else if arg == "--no-color" || arg == "-M" {
			colorOutput = false
		} else if !strings.HasPrefix(arg, "-") && !filterSet {
			filter = arg
			filterSet = true
		} else if !strings.HasPrefix(arg, "-") && filterSet && inputFile == "" {
			inputFile = arg
		}
	}

	// Auto-detect color support if not explicitly set
	if !colorOutput && os.Getenv("NO_COLOR") == "" {
		// Check if output is a terminal
		if fileInfo, _ := os.Stdout.Stat(); (fileInfo.Mode() & os.ModeCharDevice) != 0 {
			colorOutput = true
		}
	}

	// Read input
	var input []byte
	var err error

	if inputFile != "" {
		input, err = os.ReadFile(inputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file '%s': %v\n", inputFile, err)
			os.Exit(1)
		}
	} else {
		// Check if stdin is piped
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) != 0 {
			// No piped input and no file specified
			fmt.Fprintf(os.Stderr, "Error: No input provided\n")
			fmt.Fprintf(os.Stderr, "Usage: tq [filter] [file] or cat file | tq [filter]\n")
			fmt.Fprintf(os.Stderr, "Try 'tq --help' for more information\n")
			os.Exit(1)
		}

		input, err = io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
			os.Exit(1)
		}
	}

	// Convert TOON to JSON using Node.js script
	jsonData, err := toonToJSON(string(input))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing TOON:\n")
		fmt.Fprintf(os.Stderr, "  %v\n", err)
		fmt.Fprintf(os.Stderr, "\nPlease check your TOON syntax:\n")
		fmt.Fprintf(os.Stderr, "  - Verify proper indentation (2 spaces)\n")
		fmt.Fprintf(os.Stderr, "  - Check array declarations: arrayName[N]{fields}:\n")
		fmt.Fprintf(os.Stderr, "  - Ensure keys are followed by colons\n")
		os.Exit(1)
	}

	// Apply jq filter
	useColor := colorOutput && (outputFormat == "json" || outputFormat == "compact")
	result, err := applyJQ(jsonData, filter, useColor)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error applying filter '%s':\n", filter)
		fmt.Fprintf(os.Stderr, "  %v\n", err)
		fmt.Fprintf(os.Stderr, "\nFilter syntax reference:\n")
		fmt.Fprintf(os.Stderr, "  .field          - Extract field\n")
		fmt.Fprintf(os.Stderr, "  .[0]            - Get array element\n")
		fmt.Fprintf(os.Stderr, "  .[]             - Iterate array\n")
		fmt.Fprintf(os.Stderr, "  select(expr)    - Filter by condition\n")
		fmt.Fprintf(os.Stderr, "\nSee 'tq --help' for more examples\n")
		os.Exit(1)
	}

	// Output result based on format
	switch outputFormat {
	case "json":
		// Pretty-print JSON
		prettyJSON, err := formatJSON(result, true)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error formatting JSON: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(prettyJSON)

	case "compact":
		// Compact JSON (single line)
		fmt.Print(result)

	case "raw":
		// Raw values without quotes (useful for strings)
		lines := strings.Split(strings.TrimSpace(result), "\n")
		for _, line := range lines {
			if line == "" {
				continue
			}
			// Try to parse as JSON and extract raw value
			var v interface{}
			if err := json.Unmarshal([]byte(line), &v); err == nil {
				if str, ok := v.(string); ok {
					fmt.Println(str)
				} else {
					fmt.Println(line)
				}
			} else {
				fmt.Println(line)
			}
		}

	default: // "toon"
		// Handle multiple JSON values (jq can output multiple values)
		lines := strings.Split(strings.TrimSpace(result), "\n")
		for i, line := range lines {
			if line == "" {
				continue
			}
			// Convert each JSON line back to TOON
			toonOutput, err := jsonToTOON(line)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error converting to TOON: %v\n", err)
				os.Exit(1)
			}
			fmt.Print(toonOutput)
			// Add separator between multiple results
			if i < len(lines)-1 && lines[i+1] != "" {
				fmt.Println("---")
			}
		}
	}
}

func formatJSON(compactJSON string, pretty bool) (string, error) {
	if !pretty {
		return compactJSON, nil
	}

	lines := strings.Split(strings.TrimSpace(compactJSON), "\n")
	var result strings.Builder

	for _, line := range lines {
		if line == "" {
			continue
		}

		var data interface{}
		if err := json.Unmarshal([]byte(line), &data); err != nil {
			return "", err
		}

		formatted, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			return "", err
		}

		result.Write(formatted)
		result.WriteString("\n")
	}

	return result.String(), nil
}

func toonToJSON(toonInput string) (string, error) {
	scriptPath := findScript("toon-to-json.js")
	if scriptPath == "" {
		return "", fmt.Errorf("could not find toon-to-json.js script")
	}

	// Create a temporary file for input
	tmpFile, err := os.CreateTemp("", "toon-*.toon")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write([]byte(toonInput)); err != nil {
		return "", err
	}
	tmpFile.Close()

	// Run the Node.js script
	cmd := exec.Command("node", scriptPath, tmpFile.Name())
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v: %s", err, stderr.String())
	}

	return out.String(), nil
}

func applyJQ(jsonInput, filter string, color bool) (string, error) {
	// Check if jq is installed
	if _, err := exec.LookPath("jq"); err != nil {
		return "", fmt.Errorf("jq is not installed. Please install jq to use tq")
	}

	// Build jq arguments
	args := []string{"-c"}
	if color {
		args = append(args, "-C") // Color output
	} else {
		args = append(args, "-M") // Monochrome output
	}
	args = append(args, filter)

	cmd := exec.Command("jq", args...)
	cmd.Stdin = strings.NewReader(jsonInput)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v: %s", err, stderr.String())
	}

	return out.String(), nil
}

func jsonToTOON(jsonInput string) (string, error) {
	scriptPath := findScript("json-to-toon.js")
	if scriptPath == "" {
		return "", fmt.Errorf("could not find json-to-toon.js script")
	}

	// Validate JSON first
	var js interface{}
	if err := json.Unmarshal([]byte(jsonInput), &js); err != nil {
		return "", fmt.Errorf("invalid JSON: %v", err)
	}

	// Run the Node.js script
	cmd := exec.Command("node", scriptPath)
	cmd.Stdin = strings.NewReader(jsonInput)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v: %s", err, stderr.String())
	}

	return out.String(), nil
}

func findScript(scriptName string) string {
	// Try multiple locations in order
	locations := []string{
		// 1. Relative to current directory (development)
		filepath.Join("scripts", scriptName),
		// 2. Relative to executable (../scripts)
		"",
		// 3. In $HOME/.tq/scripts (installed location)
		filepath.Join(os.Getenv("HOME"), ".tq", "scripts", scriptName),
		// 4. In /usr/local/share/tq/scripts
		filepath.Join("/usr/local/share/tq/scripts", scriptName),
	}

	// Add executable-relative path
	if execPath, err := os.Executable(); err == nil {
		execDir := filepath.Dir(execPath)
		locations[1] = filepath.Join(execDir, "..", "scripts", scriptName)
	}

	for _, location := range locations {
		if location == "" {
			continue
		}
		if _, err := os.Stat(location); err == nil {
			return location
		}
	}

	return ""
}

func printHelp() {
	fmt.Println(`tq - TOON query processor (like jq for TOON format)

Usage:
  tq [options] [filter] [file]
  tq [options] [filter] < file
  cat file | tq [options] [filter]

Options:
  Output formats:
    --json         Output as pretty-printed JSON
    -c, --compact  Output as compact JSON (single line)
    -r, --raw      Output raw values (strings without quotes)
    (default)      Output as TOON format

  Display:
    -C, --color    Force colored output
    -M, --no-color Force monochrome output

  Help:
    -h, --help     Show this help message
    -v, --version  Show version

Filter Syntax (jq-compatible):
  .              Identity (return input as-is)
  .key           Extract field value
  .[0]           Extract array element by index
  .[]            Iterate array elements
  .key1.key2     Navigate nested objects
  .[].field      Extract field from each array element
  select(expr)   Filter elements by condition
  map(expr)      Transform each element
  {key: expr}    Construct new object

Examples:
  # 1. Element extraction
  tq '.' data.toon                    # Pretty print
  tq '.users' data.toon               # Extract 'users' field
  tq '.users[0]' data.toon            # Get first user
  tq '.users[0].name' data.toon       # Get first user's name

  # 2. Filtering
  tq '.users[] | select(.age > 25)' data.toon           # Users older than 25
  tq '.users[] | select(.name == "Alice")' data.toon    # User named Alice

  # 3. Output formatting
  tq '.' data.toon                    # Output as TOON (default)
  tq --json '.' data.toon             # Output as JSON

  # 4. Data transformation
  tq '.users[] | {name, email}' data.toon               # Extract specific fields
  tq '.users | map({name, older: (.age + 1)})' data.toon # Transform data

For more information, visit: https://github.com/RHEMS-japan/tq`)
}
