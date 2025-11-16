package main

import (
	"strings"
	"testing"
)

// TestBasicVariableBinding tests basic variable binding syntax
func TestBasicVariableBinding(t *testing.T) {
	tests := []struct {
		name    string
		filter  string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:   "simple variable binding",
			filter: `.age as $a | {name, age: $a}`,
			input:  `{"name":"Alice","age":30}`,
			want:   `{"name":"Alice","age":30}`,
		},
		{
			name:   "variable in computation",
			filter: `.age as $a | {name, age: $a, next_year: ($a + 1)}`,
			input:  `{"name":"Bob","age":25}`,
			want:   `{"name":"Bob","age":25,"next_year":26}`,
		},
		{
			name:   "variable with string",
			filter: `.name as $n | {original: $n, upper: ($n | ascii_upcase)}`,
			input:  `{"name":"john"}`,
			want:   `{"original":"john","upper":"JOHN"}`,
		},
		{
			name:   "variable with number",
			filter: `.price as $p | {price: $p, double: ($p * 2), half: ($p / 2)}`,
			input:  `{"price":100}`,
			want:   `{"price":100,"double":200,"half":50}`,
		},
		{
			name:   "variable in select",
			filter: `.age as $a | select($a > 25) | {name, age: $a}`,
			input:  `{"name":"Charlie","age":30}`,
			want:   `{"name":"Charlie","age":30}`,
		},
		{
			name:   "variable with false select condition",
			filter: `.age as $a | select($a > 50) | {name, age: $a}`,
			input:  `{"name":"Dave","age":30}`,
			want:   ``,
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

// TestMultipleVariables tests multiple variable bindings
func TestMultipleVariables(t *testing.T) {
	tests := []struct {
		name    string
		filter  string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:   "two variables",
			filter: `.price as $p | .quantity as $q | {price: $p, quantity: $q, total: ($p * $q)}`,
			input:  `{"price":10,"quantity":5}`,
			want:   `{"price":10,"quantity":5,"total":50}`,
		},
		{
			name:   "three variables",
			filter: `.a as $x | .b as $y | .c as $z | {sum: ($x + $y + $z)}`,
			input:  `{"a":1,"b":2,"c":3}`,
			want:   `{"sum":6}`,
		},
		{
			name:   "variables with different types",
			filter: `.name as $n | .age as $a | .active as $act | {info: "\($n) is \($a)", active: $act}`,
			input:  `{"name":"Alice","age":30,"active":true}`,
			want:   `{"info":"Alice is 30","active":true}`,
		},
		{
			name:   "chained variable assignments",
			filter: `.salary as $s | ($s * 0.1) as $bonus | {salary: $s, bonus: $bonus, total: ($s + $bonus)}`,
			input:  `{"salary":100000}`,
			want:   `{"salary":100000,"bonus":10000,"total":110000}`,
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

// TestVariablesInArrays tests variables with array operations
func TestVariablesInArrays(t *testing.T) {
	tests := []struct {
		name    string
		filter  string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:   "variable in array map",
			filter: `[.[] | . as $x | $x * 2]`,
			input:  `[1,2,3,4,5]`,
			want:   `[2,4,6,8,10]`,
		},
		{
			name:   "variable with array iteration",
			filter: `.[] | . as $item | {value: $item, doubled: ($item * 2)}`,
			input:  `[10,20,30]`,
			want:   `{"value":10,"doubled":20}` + "\n" + `{"value":20,"doubled":40}` + "\n" + `{"value":30,"doubled":60}`,
		},
		{
			name:   "variable with array filter",
			filter: `[.[] | . as $n | select($n > 15) | $n]`,
			input:  `[10,20,30,5,25]`,
			want:   `[20,30,25]`,
		},
		{
			name:   "array element variable",
			filter: `.[0] as $first | .[1] as $second | {first: $first, second: $second, sum: ($first + $second)}`,
			input:  `[100,200,300]`,
			want:   `{"first":100,"second":200,"sum":300}`,
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

// TestNestedVariables tests variables across nested iterations
func TestNestedVariables(t *testing.T) {
	tests := []struct {
		name    string
		filter  string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:   "variable across nested array",
			filter: `.[] | .name as $n | .scores[] | {user: $n, score: .}`,
			input:  `[{"name":"Alice","scores":[85,90]},{"name":"Bob","scores":[70,75]}]`,
			want:   `{"user":"Alice","score":85}` + "\n" + `{"user":"Alice","score":90}` + "\n" + `{"user":"Bob","score":70}` + "\n" + `{"user":"Bob","score":75}`,
		},
		{
			name:   "nested object with variable",
			filter: `.user.name as $n | .user.age as $a | {name: $n, age: $a, info: "\($n) is \($a) years old"}`,
			input:  `{"user":{"name":"Charlie","age":35}}`,
			want:   `{"name":"Charlie","age":35,"info":"Charlie is 35 years old"}`,
		},
		{
			name:   "variable in nested select",
			filter: `.[] | .category as $cat | .items[] | select(.price > 100) | {category: $cat, item: .name, price: .price}`,
			input:  `[{"category":"Electronics","items":[{"name":"Laptop","price":1200},{"name":"Mouse","price":25}]},{"category":"Books","items":[{"name":"Novel","price":150}]}]`,
			want:   `{"category":"Electronics","item":"Laptop","price":1200}` + "\n" + `{"category":"Books","item":"Novel","price":150}`,
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

// TestVariablesWithAggregation tests variables with aggregation functions
func TestVariablesWithAggregation(t *testing.T) {
	tests := []struct {
		name    string
		filter  string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:   "variable with add aggregation",
			filter: `[.[]] | add as $total | . | map(. / $total)`,
			input:  `[5,10,15,20]`,
			want:   `[0.1,0.2,0.3,0.4]`,
		},
		{
			name:   "variable with length",
			filter: `. | length as $count | {count: $count, items: .}`,
			input:  `[1,2,3]`,
			want:   `{"count":3,"items":[1,2,3]}`,
		},
		{
			name:   "variable with max/min",
			filter: `[.[]] | max as $mx | min as $mn | {max: $mx, min: $mn, range: ($mx - $mn)}`,
			input:  `[10,50,30,20,40]`,
			want:   `{"max":50,"min":10,"range":40}`,
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

// TestVariablesWithConditionals tests variables with conditional logic
func TestVariablesWithConditionals(t *testing.T) {
	tests := []struct {
		name    string
		filter  string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:   "variable in if condition",
			filter: `.age as $a | if $a >= 18 then "adult" else "minor" end`,
			input:  `{"age":25}`,
			want:   `"adult"`,
		},
		{
			name:   "variable in if condition - false",
			filter: `.age as $a | if $a >= 18 then "adult" else "minor" end`,
			input:  `{"age":15}`,
			want:   `"minor"`,
		},
		{
			name:   "multiple variables in conditional",
			filter: `.price as $p | .discount as $d | if $d > 0 then ($p * (1 - $d)) else $p end`,
			input:  `{"price":100,"discount":0.2}`,
			want:   `80`,
		},
		{
			name:   "variable with elif",
			filter: `.score as $s | if $s >= 90 then "A" elif $s >= 80 then "B" elif $s >= 70 then "C" else "F" end`,
			input:  `{"score":85}`,
			want:   `"B"`,
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
