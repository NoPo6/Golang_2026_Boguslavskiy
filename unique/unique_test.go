package unique

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnique_Success(t *testing.T) {
	tests := []struct {
		name     string
		commands Commands
		lines    []string
		expected []string
	}{
		{
			name:     "empty input",
			commands: Commands{},
			lines:    []string{},
			expected: []string{},
		},
		{
			name:     "single line",
			commands: Commands{},
			lines:    []string{"hello"},
			expected: []string{"hello"},
		},
		{
			name:     "default keeps one per group",
			commands: Commands{},
			lines:    []string{"a", "a", "b", "c", "c", "c", "d"},
			expected: []string{"a", "b", "c", "d"},
		},
		{
			name:     "default all unique",
			commands: Commands{},
			lines:    []string{"a", "b", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "count flag",
			commands: Commands{Count: true},
			lines:    []string{"a", "a", "b", "c", "c", "c"},
			expected: []string{"2 a", "1 b", "3 c"},
		},
		{
			name:     "count flag single group",
			commands: Commands{Count: true},
			lines:    []string{"x", "x", "x"},
			expected: []string{"3 x"},
		},
		{
			name:     "duplicates flag",
			commands: Commands{Duplicates: true},
			lines:    []string{"a", "a", "b", "c", "c", "c"},
			expected: []string{"a", "c"},
		},
		{
			name:     "duplicates flag no duplicates",
			commands: Commands{Duplicates: true},
			lines:    []string{"a", "b", "c"},
			expected: []string{},
		},
		{
			name:     "unique flag",
			commands: Commands{Unique: true},
			lines:    []string{"a", "a", "b", "c", "c", "c", "d"},
			expected: []string{"b", "d"},
		},
		{
			name:     "unique flag all duplicated",
			commands: Commands{Unique: true},
			lines:    []string{"a", "a", "b", "b"},
			expected: []string{},
		},
		{
			name:     "ignore case",
			commands: Commands{IgnoreCase: true},
			lines:    []string{"Apple", "apple", "APPLE", "banana"},
			expected: []string{"Apple", "banana"},
		},
		{
			name:     "ignore case with count",
			commands: Commands{Count: true, IgnoreCase: true},
			lines:    []string{"Apple", "apple", "APPLE"},
			expected: []string{"3 Apple"},
		},
		{
			name:     "skip fields",
			commands: Commands{SkipFields: 1},
			lines:    []string{"1 apple", "2 apple", "3 banana"},
			expected: []string{"1 apple", "3 banana"},
		},
		{
			name:     "skip fields with count",
			commands: Commands{Count: true, SkipFields: 1},
			lines:    []string{"1 apple", "2 apple", "3 banana"},
			expected: []string{"2 1 apple", "1 3 banana"},
		},
		{
			name:     "skip two fields",
			commands: Commands{SkipFields: 2},
			lines:    []string{"a b c", "x y c", "a b d"},
			expected: []string{"a b c", "a b d"},
		},
		{
			name:     "skip chars",
			commands: Commands{SkipChars: 2},
			lines:    []string{"XXapple", "YYapple", "ZZbanana"},
			expected: []string{"XXapple", "ZZbanana"},
		},
		{
			name:     "skip chars with count",
			commands: Commands{Count: true, SkipChars: 2},
			lines:    []string{"XXapple", "YYapple", "ZZbanana"},
			expected: []string{"2 XXapple", "1 ZZbanana"},
		},
		{
			name:     "skip fields and chars together",
			commands: Commands{SkipFields: 1, SkipChars: 2},
			lines:    []string{"1 XXapple", "2 YYapple", "3 ZZbanana"},
			expected: []string{"1 XXapple", "3 ZZbanana"},
		},
		{
			name:     "skip chars with russian letters",
			commands: Commands{SkipChars: 3},
			lines:    []string{"котлета", "собаки", "кошка"},
			expected: []string{"котлета", "собаки", "кошка"},
		},
		{
			name:     "skip chars longer than line groups everything",
			commands: Commands{SkipChars: 100},
			lines:    []string{"a", "b"},
			expected: []string{"a"},
		},
		{
			name:     "skip fields more than present groups everything",
			commands: Commands{SkipFields: 10},
			lines:    []string{"a b c", "d e f"},
			expected: []string{"a b c"},
		},
		{
			name:     "empty lines grouped",
			commands: Commands{},
			lines:    []string{"", "", "a", ""},
			expected: []string{"", "a", ""},
		},
		{
			name:     "duplicates with ignore case",
			commands: Commands{Duplicates: true, IgnoreCase: true},
			lines:    []string{"Apple", "apple", "banana"},
			expected: []string{"Apple"},
		},
		{
			name:     "unique with ignore case",
			commands: Commands{Unique: true, IgnoreCase: true},
			lines:    []string{"Apple", "apple", "banana", "Cherry", "cherry"},
			expected: []string{"banana"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Unique(tt.commands, tt.lines)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestUnique_Failure(t *testing.T) {
	tests := []struct {
		name     string
		commands Commands
		lines    []string
		expected []string
	}{
		{
			name:     "duplicates on all unique lines returns empty",
			commands: Commands{Duplicates: true},
			lines:    []string{"a", "b", "c", "d"},
			expected: []string{},
		},
		{
			name:     "unique on all duplicated lines returns empty",
			commands: Commands{Unique: true},
			lines:    []string{"a", "a", "b", "b", "c", "c"},
			expected: []string{},
		},
		{
			name:     "count on empty input returns empty",
			commands: Commands{Count: true},
			lines:    []string{},
			expected: []string{},
		},
		{
			name:     "duplicates on empty input returns empty",
			commands: Commands{Duplicates: true},
			lines:    []string{},
			expected: []string{},
		},
		{
			name:     "unique on empty input returns empty",
			commands: Commands{Unique: true},
			lines:    []string{},
			expected: []string{},
		},
		{
			name:     "skip all chars groups everything into one",
			commands: Commands{SkipChars: 5},
			lines:    []string{"hello", "world"},
			expected: []string{"hello"},
		},
		{
			name:     "duplicates after skip fields no matches",
			commands: Commands{Duplicates: true, SkipFields: 1},
			lines:    []string{"1 a", "2 b", "3 c"},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Unique(tt.commands, tt.lines)
			require.Equal(t, tt.expected, result)
		})
	}
}
