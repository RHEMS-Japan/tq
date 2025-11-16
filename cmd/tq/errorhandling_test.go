package main

import (
	"strings"
	"testing"
)

// TestOptionalAccess tests the optional access operator (?)
func TestOptionalAccess(t *testing.T) {
	tests := []struct {
		name    string
		filter  string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:   "optional access - field exists",
			filter: ".name?",
			input:  `{"name":"John","age":30}`,
			want:   `"John"`,
		},
		{
			name:   "optional access - field missing",
			filter: ".unknown?",
			input:  `{"name":"John"}`,
			want:   `null`,
		},
		{
			name:   "optional access - nested field exists",
			filter: ".user.name?",
			input:  `{"user":{"name":"Alice"}}`,
			want:   `"Alice"`,
		},
		{
			name:   "optional access - nested field missing",
			filter: ".user.email?",
			input:  `{"user":{"name":"Alice"}}`,
			want:   `null`,
		},
		{
			name:   "optional access - array index valid",
			filter: ".[0]?",
			input:  `[1,2,3]`,
			want:   `1`,
		},
		{
			name:   "optional access - array index out of bounds",
			filter: ".[99]?",
			input:  `[1,2,3]`,
			want:   `null`,
		},
		{
			name:   "optional access - chained",
			filter: ".user?.name?",
			input:  `{"user":null}`,
			want:   `null`,
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

// TestTryCatch tests try-catch error handling
func TestTryCatch(t *testing.T) {
	tests := []struct {
		name    string
		filter  string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:   "try-catch - no error",
			filter: `try .name catch "error"`,
			input:  `{"name":"John"}`,
			want:   `"John"`,
		},
		{
			name:   "try-catch - with error",
			filter: `try .invalid catch "default"`,
			input:  `{"name":"John"}`,
			want:   `null`,
		},
		{
			name:   "try-catch - division by zero",
			filter: `try (1/0) catch "division error"`,
			input:  `null`,
			want:   `"division error"`,
		},
		{
			name:   "try-catch - type error",
			filter: `try (.name | tonumber) catch 0`,
			input:  `{"name":"not a number"}`,
			want:   `0`,
		},
		{
			name:   "try-catch - array access",
			filter: `try .[10] catch "out of bounds"`,
			input:  `[1,2,3]`,
			want:   `null`,
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

// TestAlternativeOperatorAdvanced tests advanced alternative operator usage
func TestAlternativeOperatorAdvanced(t *testing.T) {
	tests := []struct {
		name    string
		filter  string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:   "alternative - first is null",
			filter: `.a // "default"`,
			input:  `{"a":null}`,
			want:   `"default"`,
		},
		{
			name:   "alternative - first is false",
			filter: `.a // "default"`,
			input:  `{"a":false}`,
			want:   `"default"`,
		},
		{
			name:   "alternative - first has value",
			filter: `.a // "default"`,
			input:  `{"a":"value"}`,
			want:   `"value"`,
		},
		{
			name:   "alternative - chain",
			filter: `.a // .b // .c // "none"`,
			input:  `{"a":null,"b":false,"c":"found"}`,
			want:   `"found"`,
		},
		{
			name:   "alternative - all null",
			filter: `.a // .b // "fallback"`,
			input:  `{"a":null,"b":null}`,
			want:   `"fallback"`,
		},
		{
			name:   "alternative - missing field",
			filter: `.missing // "not found"`,
			input:  `{"name":"John"}`,
			want:   `"not found"`,
		},
		{
			name:   "alternative - with zero",
			filter: `.count // 0`,
			input:  `{"name":"test"}`,
			want:   `0`,
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

// TestErrorHandlingInConstruction tests error handling in object/array construction
func TestErrorHandlingInConstruction(t *testing.T) {
	tests := []struct {
		name    string
		filter  string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:   "optional in object",
			filter: `{name, email: .email?}`,
			input:  `{"name":"John"}`,
			want:   `{"name":"John","email":null}`,
		},
		{
			name:   "alternative in object",
			filter: `{name, role: (.role // "user")}`,
			input:  `{"name":"Alice"}`,
			want:   `{"name":"Alice","role":"user"}`,
		},
		{
			name:   "try-catch in object",
			filter: `{name, age: (try .age catch 0)}`,
			input:  `{"name":"Bob"}`,
			want:   `{"name":"Bob","age":null}`,
		},
		{
			name:   "multiple alternatives",
			filter: `{name, contact: (.email // .phone // "none")}`,
			input:  `{"name":"Charlie","phone":"123"}`,
			want:   `{"name":"Charlie","contact":"123"}`,
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

// TestErrorHandlingCombinations tests combining error handling techniques
func TestErrorHandlingCombinations(t *testing.T) {
	tests := []struct {
		name    string
		filter  string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:   "optional with alternative",
			filter: `.user.name? // "anonymous"`,
			input:  `{"user":{}}`,
			want:   `"anonymous"`,
		},
		{
			name:   "try-catch with alternative",
			filter: `(try .value catch null) // "default"`,
			input:  `{"other":"field"}`,
			want:   `"default"`,
		},
		{
			name:   "nested optional access",
			filter: `.a?.b?.c? // "not found"`,
			input:  `{"a":{"b":{}}}`,
			want:   `"not found"`,
		},
		{
			name:   "optional in array iteration",
			filter: `[.items[]?.name?] | map(. // "unnamed")`,
			input:  `{"items":[{"name":"A"},{"id":2},{"name":"C"}]}`,
			want:   `["A","unnamed","C"]`,
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

// TestErrorHandlingWithConditionals tests error handling with conditionals
func TestErrorHandlingWithConditionals(t *testing.T) {
	tests := []struct {
		name    string
		filter  string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:   "conditional with optional",
			filter: `if .active? then "yes" else "no" end`,
			input:  `{"name":"test"}`,
			want:   `"no"`,
		},
		{
			name:   "conditional with alternative",
			filter: `if (.status // "unknown") == "active" then "running" else "stopped" end`,
			input:  `{"name":"service"}`,
			want:   `"stopped"`,
		},
		{
			name:   "try-catch in conditional",
			filter: `if (try .age catch 0) > 18 then "adult" else "minor" end`,
			input:  `{"name":"John"}`,
			want:   `"minor"`,
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
