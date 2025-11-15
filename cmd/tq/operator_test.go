package main

import (
	"strings"
	"testing"
)

// TestArithmeticOperators tests arithmetic operators
func TestArithmeticOperators(t *testing.T) {
	tests := []struct {
		name    string
		filter  string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:   "addition with numbers",
			filter: ".a + .b",
			input:  `{"a":10,"b":5}`,
			want:   "15",
		},
		{
			name:   "subtraction",
			filter: ".a - .b",
			input:  `{"a":10,"b":3}`,
			want:   "7",
		},
		{
			name:   "multiplication",
			filter: ".a * .b",
			input:  `{"a":4,"b":5}`,
			want:   "20",
		},
		{
			name:   "division",
			filter: ".a / .b",
			input:  `{"a":20,"b":4}`,
			want:   "5",
		},
		{
			name:   "modulo",
			filter: ".a % .b",
			input:  `{"a":10,"b":3}`,
			want:   "1",
		},
		{
			name:   "string concatenation",
			filter: `.first + " " + .last`,
			input:  `{"first":"John","last":"Doe"}`,
			want:   `"John Doe"`,
		},
		{
			name:   "array concatenation",
			filter: ".a + .b",
			input:  `{"a":[1,2],"b":[3,4]}`,
			want:   "[1,2,3,4]",
		},
		{
			name:   "object merge",
			filter: ".a + .b",
			input:  `{"a":{"x":1},"b":{"y":2}}`,
			want:   `{"x":1,"y":2}`,
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

// TestComparisonOperators tests comparison operators
func TestComparisonOperators(t *testing.T) {
	tests := []struct {
		name    string
		filter  string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:   "equal numbers",
			filter: ".a == .b",
			input:  `{"a":5,"b":5}`,
			want:   "true",
		},
		{
			name:   "not equal numbers",
			filter: ".a != .b",
			input:  `{"a":5,"b":3}`,
			want:   "true",
		},
		{
			name:   "greater than",
			filter: ".a > .b",
			input:  `{"a":10,"b":5}`,
			want:   "true",
		},
		{
			name:   "greater than or equal",
			filter: ".a >= .b",
			input:  `{"a":5,"b":5}`,
			want:   "true",
		},
		{
			name:   "less than",
			filter: ".a < .b",
			input:  `{"a":3,"b":10}`,
			want:   "true",
		},
		{
			name:   "less than or equal",
			filter: ".a <= .b",
			input:  `{"a":5,"b":5}`,
			want:   "true",
		},
		{
			name:   "string equality",
			filter: `.name == "John"`,
			input:  `{"name":"John"}`,
			want:   "true",
		},
		{
			name:   "boolean equality",
			filter: ".active == true",
			input:  `{"active":true}`,
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

// TestLogicalOperators tests logical operators
func TestLogicalOperators(t *testing.T) {
	tests := []struct {
		name    string
		filter  string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:   "and - both true",
			filter: ".a and .b",
			input:  `{"a":true,"b":true}`,
			want:   "true",
		},
		{
			name:   "and - one false",
			filter: ".a and .b",
			input:  `{"a":true,"b":false}`,
			want:   "false",
		},
		{
			name:   "or - both false",
			filter: ".a or .b",
			input:  `{"a":false,"b":false}`,
			want:   "false",
		},
		{
			name:   "or - one true",
			filter: ".a or .b",
			input:  `{"a":false,"b":true}`,
			want:   "true",
		},
		{
			name:   "not - negate true",
			filter: ".a | not",
			input:  `{"a":true}`,
			want:   "false",
		},
		{
			name:   "not - negate false",
			filter: ".a | not",
			input:  `{"a":false}`,
			want:   "true",
		},
		{
			name:   "complex condition with and",
			filter: ".age > 25 and .active == true",
			input:  `{"age":30,"active":true}`,
			want:   "true",
		},
		{
			name:   "complex condition with or",
			filter: `.role == "admin" or .role == "owner"`,
			input:  `{"role":"admin"}`,
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

// TestAlternativeOperator tests the alternative operator (//)
func TestAlternativeOperator(t *testing.T) {
	tests := []struct {
		name    string
		filter  string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:   "alternative - first is null",
			filter: ".a // .b",
			input:  `{"a":null,"b":"default"}`,
			want:   `"default"`,
		},
		{
			name:   "alternative - first is false",
			filter: ".a // .b",
			input:  `{"a":false,"b":"default"}`,
			want:   `"default"`,
		},
		{
			name:   "alternative - first is valid",
			filter: ".a // .b",
			input:  `{"a":"value","b":"default"}`,
			want:   `"value"`,
		},
		{
			name:   "alternative - chain",
			filter: `.a // .b // "none"`,
			input:  `{"a":null,"b":null}`,
			want:   `"none"`,
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

// TestOperatorsWithSelect tests operators in select() expressions
func TestOperatorsWithSelect(t *testing.T) {
	tests := []struct {
		name    string
		filter  string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:   "select with greater than",
			filter: ".[] | select(.age > 25)",
			input:  `[{"name":"Alice","age":30},{"name":"Bob","age":20}]`,
			want:   `{"name":"Alice","age":30}`,
		},
		{
			name:   "select with and",
			filter: `.[] | select(.age > 25 and .active == true)`,
			input:  `[{"age":30,"active":true},{"age":30,"active":false}]`,
			want:   `{"age":30,"active":true}`,
		},
		{
			name:   "select with or",
			filter: `.[] | select(.role == "admin" or .role == "owner")`,
			input:  `[{"name":"A","role":"admin"},{"name":"B","role":"user"},{"name":"C","role":"owner"}]`,
			want:   `{"name":"A","role":"admin"}` + "\n" + `{"name":"C","role":"owner"}`,
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

			// For multiple results, we need to normalize each line
			gotLines := strings.Split(strings.TrimSpace(result), "\n")
			wantLines := strings.Split(strings.TrimSpace(tt.want), "\n")

			if len(gotLines) != len(wantLines) {
				t.Errorf("applyJQ() returned %d results, want %d\nGot: %v\nWant: %v",
					len(gotLines), len(wantLines), result, tt.want)
				return
			}

			for i := range gotLines {
				got := normalizeJSON(strings.TrimSpace(gotLines[i]))
				want := normalizeJSON(strings.TrimSpace(wantLines[i]))
				if got != want {
					t.Errorf("applyJQ() result[%d] = %v, want %v", i, got, want)
				}
			}
		})
	}
}
