package main

import (
	"strings"
	"testing"
)

// TestObjectConstruction tests object construction syntax
func TestObjectConstruction(t *testing.T) {
	tests := []struct {
		name    string
		filter  string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:   "field shorthand",
			filter: "{name, age}",
			input:  `{"name":"John","age":30,"city":"NYC"}`,
			want:   `{"name":"John","age":30}`,
		},
		{
			name:   "field rename",
			filter: "{n: .name, a: .age}",
			input:  `{"name":"John","age":30}`,
			want:   `{"n":"John","a":30}`,
		},
		{
			name:   "computed field",
			filter: `{name, doubled: (.age * 2)}`,
			input:  `{"name":"John","age":30}`,
			want:   `{"name":"John","doubled":60}`,
		},
		{
			name:   "nested object",
			filter: `{user: {name, age}, city}`,
			input:  `{"name":"John","age":30,"city":"NYC"}`,
			want:   `{"user":{"name":"John","age":30},"city":"NYC"}`,
		},
		{
			name:   "object from array element",
			filter: ".[0] | {name, age}",
			input:  `[{"name":"Alice","age":25,"role":"dev"},{"name":"Bob","age":30}]`,
			want:   `{"name":"Alice","age":25}`,
		},
		{
			name:   "empty object",
			filter: "{}",
			input:  `{"a":1,"b":2}`,
			want:   `{}`,
		},
		{
			name:   "single field",
			filter: "{name}",
			input:  `{"name":"John","age":30}`,
			want:   `{"name":"John"}`,
		},
		{
			name:   "multiple renames",
			filter: "{firstName: .name, years: .age, location: .city}",
			input:  `{"name":"John","age":30,"city":"NYC"}`,
			want:   `{"firstName":"John","years":30,"location":"NYC"}`,
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

// TestArrayConstruction tests array construction syntax
func TestArrayConstruction(t *testing.T) {
	tests := []struct {
		name    string
		filter  string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:   "simple array",
			filter: "[.a, .b, .c]",
			input:  `{"a":1,"b":2,"c":3}`,
			want:   `[1,2,3]`,
		},
		{
			name:   "array with computation",
			filter: "[.x, (.x * 2), (.x * 3)]",
			input:  `{"x":5}`,
			want:   `[5,10,15]`,
		},
		{
			name:   "empty array",
			filter: "[]",
			input:  `{"a":1}`,
			want:   `[]`,
		},
		{
			name:   "single element array",
			filter: "[.name]",
			input:  `{"name":"John"}`,
			want:   `["John"]`,
		},
		{
			name:   "range function",
			filter: "[range(5)]",
			input:  `null`,
			want:   `[0,1,2,3,4]`,
		},
		{
			name:   "range with start and end",
			filter: "[range(2;5)]",
			input:  `null`,
			want:   `[2,3,4]`,
		},
		{
			name:   "array from iteration",
			filter: "[.[] | . * 2]",
			input:  `[1,2,3]`,
			want:   `[2,4,6]`,
		},
		{
			name:   "nested arrays",
			filter: "[[.a, .b], [.c, .d]]",
			input:  `{"a":1,"b":2,"c":3,"d":4}`,
			want:   `[[1,2],[3,4]]`,
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

// TestMapWithConstruction tests map() with object/array construction
func TestMapWithConstruction(t *testing.T) {
	tests := []struct {
		name    string
		filter  string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:   "map to object with selected fields",
			filter: "map({name, age})",
			input:  `[{"name":"Alice","age":25,"role":"dev"},{"name":"Bob","age":30,"role":"admin"}]`,
			want:   `[{"name":"Alice","age":25},{"name":"Bob","age":30}]`,
		},
		{
			name:   "map to object with rename",
			filter: "map({n: .name, a: .age})",
			input:  `[{"name":"Alice","age":25},{"name":"Bob","age":30}]`,
			want:   `[{"n":"Alice","a":25},{"n":"Bob","a":30}]`,
		},
		{
			name:   "map to object with computation",
			filter: "map({name, ageInMonths: (.age * 12)})",
			input:  `[{"name":"Alice","age":25},{"name":"Bob","age":30}]`,
			want:   `[{"name":"Alice","ageInMonths":300},{"name":"Bob","ageInMonths":360}]`,
		},
		{
			name:   "map to extract single field",
			filter: "map(.name)",
			input:  `[{"name":"Alice","age":25},{"name":"Bob","age":30}]`,
			want:   `["Alice","Bob"]`,
		},
		{
			name:   "map with nested object",
			filter: "map({person: {name, age}, role})",
			input:  `[{"name":"Alice","age":25,"role":"dev"}]`,
			want:   `[{"person":{"name":"Alice","age":25},"role":"dev"}]`,
		},
		{
			name:   "map to array",
			filter: "map([.name, .age])",
			input:  `[{"name":"Alice","age":25},{"name":"Bob","age":30}]`,
			want:   `[["Alice",25],["Bob",30]]`,
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

// TestConditionalConstruction tests construction with select()
func TestConditionalConstruction(t *testing.T) {
	tests := []struct {
		name    string
		filter  string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:   "array of filtered objects",
			filter: "[.[] | select(.age > 25)]",
			input:  `[{"name":"Alice","age":20},{"name":"Bob","age":30},{"name":"Charlie","age":35}]`,
			want:   `[{"name":"Bob","age":30},{"name":"Charlie","age":35}]`,
		},
		{
			name:   "array of filtered and transformed",
			filter: "[.[] | select(.active) | {name, role}]",
			input:  `[{"name":"Alice","role":"dev","active":true},{"name":"Bob","role":"admin","active":false}]`,
			want:   `[{"name":"Alice","role":"dev"}]`,
		},
		{
			name:   "object from first match",
			filter: ".[] | select(.age > 25) | {name, age}",
			input:  `[{"name":"Alice","age":20},{"name":"Bob","age":30}]`,
			want:   `{"name":"Bob","age":30}`,
		},
		{
			name:   "array of names from filter",
			filter: "[.[] | select(.age >= 30) | .name]",
			input:  `[{"name":"Alice","age":25},{"name":"Bob","age":30},{"name":"Charlie","age":35}]`,
			want:   `["Bob","Charlie"]`,
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

// TestComplexConstruction tests complex nested construction
func TestComplexConstruction(t *testing.T) {
	tests := []struct {
		name    string
		filter  string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:   "deeply nested object",
			filter: `{user: {profile: {name, age}, contact: {email}}, status}`,
			input:  `{"name":"John","age":30,"email":"john@example.com","status":"active"}`,
			want:   `{"user":{"profile":{"name":"John","age":30},"contact":{"email":"john@example.com"}},"status":"active"}`,
		},
		{
			name:   "array of objects with arrays",
			filter: "map({name, scores: [.math, .english, .science]})",
			input:  `[{"name":"Alice","math":90,"english":85,"science":88}]`,
			want:   `[{"name":"Alice","scores":[90,85,88]}]`,
		},
		{
			name:   "object with computed nested array",
			filter: `{name, doubledScores: [.scores[] | . * 2]}`,
			input:  `{"name":"Alice","scores":[1,2,3]}`,
			want:   `{"name":"Alice","doubledScores":[2,4,6]}`,
		},
		{
			name:   "mixed construction",
			filter: `{info: {name, age}, tags: [.role, .status], count: 1}`,
			input:  `{"name":"John","age":30,"role":"admin","status":"active"}`,
			want:   `{"info":{"name":"John","age":30},"tags":["admin","active"],"count":1}`,
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
