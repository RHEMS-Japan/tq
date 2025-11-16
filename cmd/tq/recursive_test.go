package main

import (
	"strings"
	"testing"
)

// TestBasicRecursiveDescent tests basic recursive descent (..) syntax
func TestBasicRecursiveDescent(t *testing.T) {
	tests := []struct {
		name    string
		filter  string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:   "recursive descent - all values",
			filter: "..",
			input:  `{"a":1,"b":{"c":2}}`,
			want:   `{"a":1,"b":{"c":2}}` + "\n" + `1` + "\n" + `{"c":2}` + "\n" + `2`,
		},
		{
			name:   "recursive descent - simple object",
			filter: "..",
			input:  `{"x":10}`,
			want:   `{"x":10}` + "\n" + `10`,
		},
		{
			name:   "recursive descent - nested object",
			filter: "..",
			input:  `{"a":{"b":{"c":1}}}`,
			want:   `{"a":{"b":{"c":1}}}` + "\n" + `{"b":{"c":1}}` + "\n" + `{"c":1}` + "\n" + `1`,
		},
		{
			name:   "recursive descent - array",
			filter: "..",
			input:  `[1,2,[3,4]]`,
			want:   `[1,2,[3,4]]` + "\n" + `1` + "\n" + `2` + "\n" + `[3,4]` + "\n" + `3` + "\n" + `4`,
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

			got := strings.TrimSpace(result)
			want := tt.want

			// For multiple results, normalize each line
			if strings.Contains(want, "\n") {
				gotLines := strings.Split(got, "\n")
				wantLines := strings.Split(want, "\n")

				if len(gotLines) != len(wantLines) {
					t.Errorf("applyJQ() returned %d results, want %d\nGot:\n%s\nWant:\n%s",
						len(gotLines), len(wantLines), got, want)
					return
				}

				for i := range gotLines {
					gotNorm := normalizeJSON(strings.TrimSpace(gotLines[i]))
					wantNorm := normalizeJSON(strings.TrimSpace(wantLines[i]))
					if gotNorm != wantNorm {
						t.Errorf("applyJQ() result[%d] = %v, want %v", i, gotNorm, wantNorm)
					}
				}
			} else {
				gotNorm := normalizeJSON(got)
				wantNorm := normalizeJSON(want)
				if gotNorm != wantNorm {
					t.Errorf("applyJQ() = %v, want %v", gotNorm, wantNorm)
				}
			}
		})
	}
}

// TestRecursiveWithTypeFilter tests recursive descent with type filtering
func TestRecursiveWithTypeFilter(t *testing.T) {
	tests := []struct {
		name    string
		filter  string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:   "find all numbers",
			filter: `.. | select(type == "number")`,
			input:  `{"a":1,"b":{"c":2,"d":"text"},"e":3}`,
			want:   `1` + "\n" + `2` + "\n" + `3`,
		},
		{
			name:   "find all strings",
			filter: `.. | select(type == "string")`,
			input:  `{"a":"hello","b":{"c":"world"},"d":123}`,
			want:   `"hello"` + "\n" + `"world"`,
		},
		{
			name:   "find all objects",
			filter: `.. | select(type == "object")`,
			input:  `{"a":{"b":1},"c":{"d":2}}`,
			want:   `{"a":{"b":1},"c":{"d":2}}` + "\n" + `{"b":1}` + "\n" + `{"d":2}`,
		},
		{
			name:   "find all arrays",
			filter: `.. | select(type == "array")`,
			input:  `{"a":[1,2],"b":{"c":[3,4]}}`,
			want:   `[1,2]` + "\n" + `[3,4]`,
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

			got := strings.TrimSpace(result)
			want := tt.want

			gotLines := strings.Split(got, "\n")
			wantLines := strings.Split(want, "\n")

			if len(gotLines) != len(wantLines) {
				t.Errorf("applyJQ() returned %d results, want %d", len(gotLines), len(wantLines))
				return
			}

			for i := range gotLines {
				gotNorm := normalizeJSON(strings.TrimSpace(gotLines[i]))
				wantNorm := normalizeJSON(strings.TrimSpace(wantLines[i]))
				if gotNorm != wantNorm {
					t.Errorf("applyJQ() result[%d] = %v, want %v", i, gotNorm, wantNorm)
				}
			}
		})
	}
}

// TestRecursiveWithFieldSearch tests finding specific fields recursively
func TestRecursiveWithFieldSearch(t *testing.T) {
	tests := []struct {
		name    string
		filter  string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:   "find all 'name' fields",
			filter: `.. | .name? // empty`,
			input:  `{"name":"John","user":{"name":"Alice","age":30},"item":{"price":100}}`,
			want:   `"John"` + "\n" + `"Alice"`,
		},
		{
			name:   "find specific field value",
			filter: `.. | select(.name? == "Alice")`,
			input:  `{"users":[{"name":"Alice","age":30},{"name":"Bob","age":25}]}`,
			want:   `{"name":"Alice","age":30}`,
		},
		{
			name:   "find nested field",
			filter: `.. | .email? // empty`,
			input:  `{"user":{"profile":{"email":"test@example.com"}},"admin":{"email":"admin@example.com"}}`,
			want:   `"test@example.com"` + "\n" + `"admin@example.com"`,
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

			got := strings.TrimSpace(result)
			want := tt.want

			gotLines := strings.Split(got, "\n")
			wantLines := strings.Split(want, "\n")

			if len(gotLines) != len(wantLines) {
				t.Errorf("applyJQ() returned %d results, want %d", len(gotLines), len(wantLines))
				return
			}

			for i := range gotLines {
				gotNorm := normalizeJSON(strings.TrimSpace(gotLines[i]))
				wantNorm := normalizeJSON(strings.TrimSpace(wantLines[i]))
				if gotNorm != wantNorm {
					t.Errorf("applyJQ() result[%d] = %v, want %v", i, gotNorm, wantNorm)
				}
			}
		})
	}
}

// TestRecursiveWithConditions tests recursive descent with conditional filters
func TestRecursiveWithConditions(t *testing.T) {
	tests := []struct {
		name    string
		filter  string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:   "find numbers greater than 10",
			filter: `.. | select(type == "number" and . > 10)`,
			input:  `{"a":5,"b":{"c":15,"d":20},"e":8}`,
			want:   `15` + "\n" + `20`,
		},
		{
			name:   "find strings starting with prefix",
			filter: `.. | select(type == "string" and startswith("test"))`,
			input:  `{"a":"test1","b":{"c":"other","d":"test2"}}`,
			want:   `"test1"` + "\n" + `"test2"`,
		},
		{
			name:   "find objects with specific field",
			filter: `.. | select(type == "object" and has("id"))`,
			input:  `{"users":[{"id":1,"name":"Alice"},{"name":"Bob"},{"id":2,"name":"Charlie"}]}`,
			want:   `{"id":1,"name":"Alice"}` + "\n" + `{"id":2,"name":"Charlie"}`,
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

			got := strings.TrimSpace(result)
			want := tt.want

			gotLines := strings.Split(got, "\n")
			wantLines := strings.Split(want, "\n")

			if len(gotLines) != len(wantLines) {
				t.Errorf("applyJQ() returned %d results, want %d", len(gotLines), len(wantLines))
				return
			}

			for i := range gotLines {
				gotNorm := normalizeJSON(strings.TrimSpace(gotLines[i]))
				wantNorm := normalizeJSON(strings.TrimSpace(wantLines[i]))
				if gotNorm != wantNorm {
					t.Errorf("applyJQ() result[%d] = %v, want %v", i, gotNorm, wantNorm)
				}
			}
		})
	}
}

// TestRecursiveWithArrays tests recursive descent with arrays
func TestRecursiveWithArrays(t *testing.T) {
	tests := []struct {
		name    string
		filter  string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:   "find all array elements",
			filter: `.. | select(type == "number")`,
			input:  `{"data":[1,2,3],"nested":{"values":[4,5]}}`,
			want:   `1` + "\n" + `2` + "\n" + `3` + "\n" + `4` + "\n" + `5`,
		},
		{
			name:   "find objects in nested arrays",
			filter: `.. | select(type == "object" and .name?)`,
			input:  `{"items":[{"name":"A","value":1},{"value":2},{"name":"B","value":3}]}`,
			want:   `{"name":"A","value":1}` + "\n" + `{"name":"B","value":3}`,
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

			got := strings.TrimSpace(result)
			want := tt.want

			gotLines := strings.Split(got, "\n")
			wantLines := strings.Split(want, "\n")

			if len(gotLines) != len(wantLines) {
				t.Errorf("applyJQ() returned %d results, want %d", len(gotLines), len(wantLines))
				return
			}

			for i := range gotLines {
				gotNorm := normalizeJSON(strings.TrimSpace(gotLines[i]))
				wantNorm := normalizeJSON(strings.TrimSpace(wantLines[i]))
				if gotNorm != wantNorm {
					t.Errorf("applyJQ() result[%d] = %v, want %v", i, gotNorm, wantNorm)
				}
			}
		})
	}
}
