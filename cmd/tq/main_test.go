package main

import (
	"os"
	"testing"
)

func TestFindScript(t *testing.T) {
	tests := []struct {
		name       string
		scriptName string
		wantFound  bool
	}{
		{
			name:       "find toon-to-json.js",
			scriptName: "toon-to-json.js",
			wantFound:  true,
		},
		{
			name:       "find json-to-toon.js",
			scriptName: "json-to-toon.js",
			wantFound:  true,
		},
		{
			name:       "non-existent script",
			scriptName: "non-existent.js",
			wantFound:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := findScript(tt.scriptName)
			found := result != ""
			if found != tt.wantFound {
				t.Errorf("findScript(%q) found=%v, want %v", tt.scriptName, found, tt.wantFound)
			}
		})
	}
}

func TestToonToJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "simple object",
			input:   "name: John\nage: 30",
			wantErr: false,
		},
		{
			name:    "empty input",
			input:   "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := toonToJSON(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("toonToJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestApplyJQ(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		filter  string
		color   bool
		wantErr bool
	}{
		{
			name:    "identity filter",
			input:   `{"name":"John"}`,
			filter:  ".",
			color:   false,
			wantErr: false,
		},
		{
			name:    "field extraction",
			input:   `{"name":"John","age":30}`,
			filter:  ".name",
			color:   false,
			wantErr: false,
		},
		{
			name:    "with color",
			input:   `{"name":"John"}`,
			filter:  ".",
			color:   true,
			wantErr: false,
		},
		{
			name:    "invalid filter",
			input:   `{"name":"John"}`,
			filter:  "..invalid",
			color:   false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := applyJQ(tt.input, tt.filter, tt.color)
			if (err != nil) != tt.wantErr {
				t.Errorf("applyJQ() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestJSONToTOON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "simple object",
			input:   `{"name":"John","age":30}`,
			wantErr: false,
		},
		{
			name:    "array",
			input:   `[1,2,3]`,
			wantErr: false,
		},
		{
			name:    "invalid JSON",
			input:   `{invalid}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := jsonToTOON(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("jsonToTOON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIntegration(t *testing.T) {
	// Test with actual sample file
	input, err := os.ReadFile("../../testdata/sample.toon")
	if err != nil {
		t.Skip("testdata not available")
		return
	}

	// Convert TOON to JSON
	jsonData, err := toonToJSON(string(input))
	if err != nil {
		t.Fatalf("toonToJSON failed: %v", err)
	}

	// Apply filter without color
	result, err := applyJQ(jsonData, ".name", false)
	if err != nil {
		t.Fatalf("applyJQ failed: %v", err)
	}

	// Convert back to TOON
	_, err = jsonToTOON(result)
	if err != nil {
		t.Fatalf("jsonToTOON failed: %v", err)
	}

	// Test with color
	_, err = applyJQ(jsonData, ".", true)
	if err != nil {
		t.Fatalf("applyJQ with color failed: %v", err)
	}
}
