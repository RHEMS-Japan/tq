package main

import (
	"strings"
	"testing"
)

// TestStringInterpolation tests string interpolation syntax
func TestStringInterpolation(t *testing.T) {
	tests := []struct {
		name    string
		filter  string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:   "basic interpolation",
			filter: `"\(.name) is \(.age) years old"`,
			input:  `{"name":"John","age":30}`,
			want:   `"John is 30 years old"`,
		},
		{
			name:   "multiple fields",
			filter: `"\(.name) <\(.email)>"`,
			input:  `{"name":"John","email":"john@example.com"}`,
			want:   `"John <john@example.com>"`,
		},
		{
			name:   "with computation",
			filter: `"Total: \(.price * .quantity)"`,
			input:  `{"price":10,"quantity":5}`,
			want:   `"Total: 50"`,
		},
		{
			name:   "nested field access",
			filter: `"\(.user.name) from \(.user.city)"`,
			input:  `{"user":{"name":"Alice","city":"NYC"}}`,
			want:   `"Alice from NYC"`,
		},
		{
			name:   "with string operations",
			filter: `"\(.name | ascii_upcase)"`,
			input:  `{"name":"john"}`,
			want:   `"JOHN"`,
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

// TestStringSplitJoin tests split and join functions
func TestStringSplitJoin(t *testing.T) {
	tests := []struct {
		name    string
		filter  string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:   "split by delimiter",
			filter: `split("@")`,
			input:  `"user@example.com"`,
			want:   `["user","example.com"]`,
		},
		{
			name:   "split by space",
			filter: `split(" ")`,
			input:  `"hello world test"`,
			want:   `["hello","world","test"]`,
		},
		{
			name:   "join with comma",
			filter: `join(", ")`,
			input:  `["a","b","c"]`,
			want:   `"a, b, c"`,
		},
		{
			name:   "join with space",
			filter: `join(" ")`,
			input:  `["hello","world"]`,
			want:   `"hello world"`,
		},
		{
			name:   "join empty array",
			filter: `join(",")`,
			input:  `[]`,
			want:   `""`,
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

// TestStringCase tests case conversion functions
func TestStringCase(t *testing.T) {
	tests := []struct {
		name    string
		filter  string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:   "ascii_upcase",
			filter: "ascii_upcase",
			input:  `"hello world"`,
			want:   `"HELLO WORLD"`,
		},
		{
			name:   "ascii_downcase",
			filter: "ascii_downcase",
			input:  `"HELLO WORLD"`,
			want:   `"hello world"`,
		},
		{
			name:   "ascii_upcase mixed case",
			filter: "ascii_upcase",
			input:  `"Hello World"`,
			want:   `"HELLO WORLD"`,
		},
		{
			name:   "ascii_downcase mixed case",
			filter: "ascii_downcase",
			input:  `"Hello World"`,
			want:   `"hello world"`,
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

// TestStringTrim tests ltrimstr and rtrimstr functions
func TestStringTrim(t *testing.T) {
	tests := []struct {
		name    string
		filter  string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:   "ltrimstr - remove prefix",
			filter: `ltrimstr("https://")`,
			input:  `"https://example.com"`,
			want:   `"example.com"`,
		},
		{
			name:   "ltrimstr - no match",
			filter: `ltrimstr("http://")`,
			input:  `"https://example.com"`,
			want:   `"https://example.com"`,
		},
		{
			name:   "rtrimstr - remove suffix",
			filter: `rtrimstr(".html")`,
			input:  `"index.html"`,
			want:   `"index"`,
		},
		{
			name:   "rtrimstr - no match",
			filter: `rtrimstr(".txt")`,
			input:  `"index.html"`,
			want:   `"index.html"`,
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

// TestStringRegex tests regex functions (test, match)
func TestStringRegex(t *testing.T) {
	tests := []struct {
		name    string
		filter  string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:   "test - match found",
			filter: `test("@example.com$")`,
			input:  `"user@example.com"`,
			want:   `true`,
		},
		{
			name:   "test - no match",
			filter: `test("@example.com$")`,
			input:  `"user@other.com"`,
			want:   `false`,
		},
		{
			name:   "test - pattern at start",
			filter: `test("^hello")`,
			input:  `"hello world"`,
			want:   `true`,
		},
		{
			name:   "test - case sensitive",
			filter: `test("^Hello")`,
			input:  `"hello world"`,
			want:   `false`,
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

// TestStringReplace tests sub and gsub functions
func TestStringReplace(t *testing.T) {
	tests := []struct {
		name    string
		filter  string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:   "sub - replace first",
			filter: `sub("foo"; "bar")`,
			input:  `"foo foo foo"`,
			want:   `"bar foo foo"`,
		},
		{
			name:   "gsub - replace all",
			filter: `gsub("foo"; "bar")`,
			input:  `"foo foo foo"`,
			want:   `"bar bar bar"`,
		},
		{
			name:   "sub - no match",
			filter: `sub("baz"; "qux")`,
			input:  `"foo bar"`,
			want:   `"foo bar"`,
		},
		{
			name:   "gsub - replace spaces",
			filter: `gsub(" "; "_")`,
			input:  `"hello world test"`,
			want:   `"hello_world_test"`,
		},
		{
			name:   "gsub - remove character",
			filter: `gsub("-"; "")`,
			input:  `"hello-world-test"`,
			want:   `"helloworldtest"`,
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

// TestStringFormat tests format functions (@base64, @uri, @json)
func TestStringFormat(t *testing.T) {
	tests := []struct {
		name    string
		filter  string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:   "@uri encode",
			filter: "@uri",
			input:  `"hello world"`,
			want:   `"hello%20world"`,
		},
		{
			name:   "@uri with special chars",
			filter: "@uri",
			input:  `"foo@bar.com"`,
			want:   `"foo%40bar.com"`,
		},
		{
			name:   "@base64 encode",
			filter: "@base64",
			input:  `"hello"`,
			want:   `"aGVsbG8="`,
		},
		{
			name:   "@json encode",
			filter: "@json",
			input:  `"hello"`,
			want:   `"\"hello\""`,
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

// TestStringCombinations tests combining multiple string functions
func TestStringCombinations(t *testing.T) {
	tests := []struct {
		name    string
		filter  string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:   "split and join",
			filter: `split("@") | join(" AT ")`,
			input:  `"user@example.com"`,
			want:   `"user AT example.com"`,
		},
		{
			name:   "trim and upcase",
			filter: `ltrimstr("Mr. ") | ascii_upcase`,
			input:  `"Mr. John Doe"`,
			want:   `"JOHN DOE"`,
		},
		{
			name:   "replace and interpolate",
			filter: `gsub(" "; "_") | "username: \(.)"`,
			input:  `"John Doe"`,
			want:   `"username: John_Doe"`,
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
