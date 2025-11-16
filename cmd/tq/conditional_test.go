package main

import (
	"strings"
	"testing"
)

// TestBasicConditional tests basic if-then-else syntax
func TestBasicConditional(t *testing.T) {
	tests := []struct {
		name    string
		filter  string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:   "if-then-else true",
			filter: `if .age > 25 then "adult" else "young" end`,
			input:  `{"age":30}`,
			want:   `"adult"`,
		},
		{
			name:   "if-then-else false",
			filter: `if .age > 25 then "adult" else "young" end`,
			input:  `{"age":20}`,
			want:   `"young"`,
		},
		{
			name:   "numeric comparison",
			filter: `if .salary > 90000 then "high" else "normal" end`,
			input:  `{"salary":95000}`,
			want:   `"high"`,
		},
		{
			name:   "string comparison",
			filter: `if .role == "admin" then "admin user" else "regular user" end`,
			input:  `{"role":"admin"}`,
			want:   `"admin user"`,
		},
		{
			name:   "boolean condition",
			filter: `if .active then "yes" else "no" end`,
			input:  `{"active":true}`,
			want:   `"yes"`,
		},
		{
			name:   "boolean condition false",
			filter: `if .active then "yes" else "no" end`,
			input:  `{"active":false}`,
			want:   `"no"`,
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

// TestElifConditional tests elif (else if) syntax
func TestElifConditional(t *testing.T) {
	tests := []struct {
		name    string
		filter  string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:   "elif first condition true",
			filter: `if .score >= 90 then "A" elif .score >= 80 then "B" else "C" end`,
			input:  `{"score":95}`,
			want:   `"A"`,
		},
		{
			name:   "elif second condition true",
			filter: `if .score >= 90 then "A" elif .score >= 80 then "B" else "C" end`,
			input:  `{"score":85}`,
			want:   `"B"`,
		},
		{
			name:   "elif else clause",
			filter: `if .score >= 90 then "A" elif .score >= 80 then "B" else "C" end`,
			input:  `{"score":75}`,
			want:   `"C"`,
		},
		{
			name:   "multiple elif",
			filter: `if .n >= 100 then "A" elif .n >= 75 then "B" elif .n >= 50 then "C" else "D" end`,
			input:  `{"n":60}`,
			want:   `"C"`,
		},
		{
			name:   "salary grading",
			filter: `if .salary >= 100000 then "senior" elif .salary >= 85000 then "mid" else "junior" end`,
			input:  `{"salary":90000}`,
			want:   `"mid"`,
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

// TestNestedConditional tests nested if-then-else
func TestNestedConditional(t *testing.T) {
	tests := []struct {
		name    string
		filter  string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:   "nested both true",
			filter: `if .active then (if .premium then "premium-active" else "standard-active" end) else "inactive" end`,
			input:  `{"active":true,"premium":true}`,
			want:   `"premium-active"`,
		},
		{
			name:   "nested outer true inner false",
			filter: `if .active then (if .premium then "premium-active" else "standard-active" end) else "inactive" end`,
			input:  `{"active":true,"premium":false}`,
			want:   `"standard-active"`,
		},
		{
			name:   "nested outer false",
			filter: `if .active then (if .premium then "premium-active" else "standard-active" end) else "inactive" end`,
			input:  `{"active":false,"premium":true}`,
			want:   `"inactive"`,
		},
		{
			name:   "triple nested",
			filter: `if .a then (if .b then (if .c then "abc" else "ab" end) else "a" end) else "none" end`,
			input:  `{"a":true,"b":true,"c":true}`,
			want:   `"abc"`,
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

// TestConditionalWithLogic tests conditionals with logical operators
func TestConditionalWithLogic(t *testing.T) {
	tests := []struct {
		name    string
		filter  string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:   "and operator both true",
			filter: `if .age > 18 and .verified then "allowed" else "denied" end`,
			input:  `{"age":25,"verified":true}`,
			want:   `"allowed"`,
		},
		{
			name:   "and operator one false",
			filter: `if .age > 18 and .verified then "allowed" else "denied" end`,
			input:  `{"age":25,"verified":false}`,
			want:   `"denied"`,
		},
		{
			name:   "or operator one true",
			filter: `if .admin or .moderator then "staff" else "user" end`,
			input:  `{"admin":false,"moderator":true}`,
			want:   `"staff"`,
		},
		{
			name:   "or operator both false",
			filter: `if .admin or .moderator then "staff" else "user" end`,
			input:  `{"admin":false,"moderator":false}`,
			want:   `"user"`,
		},
		{
			name:   "not operator",
			filter: `if .active | not then "disabled" else "enabled" end`,
			input:  `{"active":false}`,
			want:   `"disabled"`,
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

// TestConditionalInConstruction tests conditionals in object/array construction
func TestConditionalInConstruction(t *testing.T) {
	tests := []struct {
		name    string
		filter  string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:   "conditional in object",
			filter: `{name, status: (if .active then "Active" else "Inactive" end)}`,
			input:  `{"name":"John","active":true}`,
			want:   `{"name":"John","status":"Active"}`,
		},
		{
			name:   "conditional value in object",
			filter: `{name, level: (if .score >= 90 then "A" elif .score >= 80 then "B" else "C" end)}`,
			input:  `{"name":"Alice","score":85}`,
			want:   `{"name":"Alice","level":"B"}`,
		},
		{
			name:   "multiple conditionals in object",
			filter: `{name, status: (if .active then "active" else "inactive" end), tier: (if .premium then "premium" else "free" end)}`,
			input:  `{"name":"Bob","active":true,"premium":false}`,
			want:   `{"name":"Bob","status":"active","tier":"free"}`,
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

// TestConditionalInStringInterpolation tests conditionals in string interpolation
func TestConditionalInStringInterpolation(t *testing.T) {
	tests := []struct {
		name    string
		filter  string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:   "conditional in interpolation",
			filter: `"\(.name) is \(if .active then "active" else "inactive" end)"`,
			input:  `{"name":"Alice","active":true}`,
			want:   `"Alice is active"`,
		},
		{
			name:   "multiple conditionals in string",
			filter: `"User: \(.name) - Status: \(if .active then "✓" else "✗" end) - Tier: \(if .premium then "Premium" else "Free" end)"`,
			input:  `{"name":"Bob","active":false,"premium":true}`,
			want:   `"User: Bob - Status: ✗ - Tier: Premium"`,
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

// TestConditionalWithEmpty tests conditionals with empty
func TestConditionalWithEmpty(t *testing.T) {
	tests := []struct {
		name    string
		filter  string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:   "return value when true",
			filter: `.[] | if .age > 25 then .name else empty end`,
			input:  `[{"name":"Alice","age":30},{"name":"Bob","age":20}]`,
			want:   `"Alice"`,
		},
		{
			name:   "empty when false",
			filter: `.[] | if .active then .name else empty end`,
			input:  `[{"name":"Alice","active":true},{"name":"Bob","active":false},{"name":"Charlie","active":true}]`,
			want:   `"Alice"` + "\n" + `"Charlie"`,
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
