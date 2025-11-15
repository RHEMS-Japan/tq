package main

import (
	"encoding/json"
	"strings"
	"testing"
)

// TestBuiltinArrayFunctions tests jq built-in array functions
func TestBuiltinArrayFunctions(t *testing.T) {
	tests := []struct {
		name    string
		filter  string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:   "length of array",
			filter: "length",
			input:  `[1,2,3,4,5]`,
			want:   "5",
		},
		{
			name:   "length of object",
			filter: "length",
			input:  `{"a":1,"b":2,"c":3}`,
			want:   "3",
		},
		{
			name:   "length of string",
			filter: "length",
			input:  `"hello"`,
			want:   "5",
		},
		{
			name:   "reverse array",
			filter: "reverse",
			input:  `[1,2,3]`,
			want:   "[3,2,1]",
		},
		{
			name:   "sort array",
			filter: "sort",
			input:  `[3,1,2]`,
			want:   "[1,2,3]",
		},
		{
			name:   "unique array",
			filter: "unique",
			input:  `[1,2,2,3,1]`,
			want:   "[1,2,3]",
		},
		{
			name:   "add numbers",
			filter: "add",
			input:  `[1,2,3,4,5]`,
			want:   "15",
		},
		{
			name:   "add strings",
			filter: "add",
			input:  `["hello"," ","world"]`,
			want:   `"hello world"`,
		},
		{
			name:   "min",
			filter: "min",
			input:  `[5,2,8,1,9]`,
			want:   "1",
		},
		{
			name:   "max",
			filter: "max",
			input:  `[5,2,8,1,9]`,
			want:   "9",
		},
		{
			name:   "first element",
			filter: "first",
			input:  `[1,2,3]`,
			want:   "1",
		},
		{
			name:   "last element",
			filter: "last",
			input:  `[1,2,3]`,
			want:   "3",
		},
		{
			name:   "flatten",
			filter: "flatten",
			input:  `[[1,2],[3,4]]`,
			want:   "[1,2,3,4]",
		},
		{
			name:   "flatten deep",
			filter: "flatten",
			input:  `[[1,[2,3]],[[4]]]`,
			want:   "[1,2,3,4]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := applyJQ(tt.input, tt.filter, false)
			if (err != nil) != tt.wantErr {
				t.Errorf("applyJQ() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}

			// Normalize JSON for comparison
			got := normalizeJSON(strings.TrimSpace(result))
			want := normalizeJSON(tt.want)

			if got != want {
				t.Errorf("applyJQ() = %v, want %v", got, want)
			}
		})
	}
}

// TestBuiltinObjectFunctions tests jq built-in object functions
func TestBuiltinObjectFunctions(t *testing.T) {
	tests := []struct {
		name    string
		filter  string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:   "keys of object",
			filter: "keys",
			input:  `{"c":3,"a":1,"b":2}`,
			want:   `["a","b","c"]`,
		},
		{
			name:   "keys of array",
			filter: "keys",
			input:  `[1,2,3]`,
			want:   `[0,1,2]`,
		},
		{
			name:   "has key true",
			filter: `has("name")`,
			input:  `{"name":"John","age":30}`,
			want:   "true",
		},
		{
			name:   "has key false",
			filter: `has("unknown")`,
			input:  `{"name":"John"}`,
			want:   "false",
		},
		{
			name:   "to_entries",
			filter: "to_entries",
			input:  `{"a":1,"b":2}`,
			want:   `[{"key":"a","value":1},{"key":"b","value":2}]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := applyJQ(tt.input, tt.filter, false)
			if (err != nil) != tt.wantErr {
				t.Errorf("applyJQ() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}

			got := normalizeJSON(strings.TrimSpace(result))
			want := normalizeJSON(tt.want)

			if got != want {
				t.Errorf("applyJQ() = %v, want %v", got, want)
			}
		})
	}
}

// TestBuiltinTypeFunctions tests jq built-in type/info functions
func TestBuiltinTypeFunctions(t *testing.T) {
	tests := []struct {
		name    string
		filter  string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:   "type object",
			filter: "type",
			input:  `{"a":1}`,
			want:   `"object"`,
		},
		{
			name:   "type array",
			filter: "type",
			input:  `[1,2,3]`,
			want:   `"array"`,
		},
		{
			name:   "type string",
			filter: "type",
			input:  `"hello"`,
			want:   `"string"`,
		},
		{
			name:   "type number",
			filter: "type",
			input:  `42`,
			want:   `"number"`,
		},
		{
			name:   "type boolean",
			filter: "type",
			input:  `true`,
			want:   `"boolean"`,
		},
		{
			name:   "type null",
			filter: "type",
			input:  `null`,
			want:   `"null"`,
		},
		{
			name:   "tonumber",
			filter: "tonumber",
			input:  `"42"`,
			want:   "42",
		},
		{
			name:   "tostring",
			filter: "tostring",
			input:  `42`,
			want:   `"42"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := applyJQ(tt.input, tt.filter, false)
			if (err != nil) != tt.wantErr {
				t.Errorf("applyJQ() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}

			got := normalizeJSON(strings.TrimSpace(result))
			want := normalizeJSON(tt.want)

			if got != want {
				t.Errorf("applyJQ() = %v, want %v", got, want)
			}
		})
	}
}

// TestBuiltinStringFunctions tests jq built-in string functions
func TestBuiltinStringFunctions(t *testing.T) {
	tests := []struct {
		name    string
		filter  string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:   "startswith true",
			filter: `startswith("hello")`,
			input:  `"hello world"`,
			want:   "true",
		},
		{
			name:   "startswith false",
			filter: `startswith("world")`,
			input:  `"hello world"`,
			want:   "false",
		},
		{
			name:   "endswith true",
			filter: `endswith("world")`,
			input:  `"hello world"`,
			want:   "true",
		},
		{
			name:   "contains true",
			filter: `contains("llo wor")`,
			input:  `"hello world"`,
			want:   "true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := applyJQ(tt.input, tt.filter, false)
			if (err != nil) != tt.wantErr {
				t.Errorf("applyJQ() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}

			got := normalizeJSON(strings.TrimSpace(result))
			want := normalizeJSON(tt.want)

			if got != want {
				t.Errorf("applyJQ() = %v, want %v", got, want)
			}
		})
	}
}

// normalizeJSON normalizes JSON string for comparison
func normalizeJSON(s string) string {
	s = strings.TrimSpace(s)

	// Try to parse and re-encode as JSON
	var data interface{}
	if err := json.Unmarshal([]byte(s), &data); err == nil {
		if b, err := json.Marshal(data); err == nil {
			return string(b)
		}
	}

	return s
}
