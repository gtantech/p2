package routes

import (
	"slices"
	"testing"
)

func TestRemoveDuplicates(t *testing.T) {
	tests := []struct {
		testName string
		values   []string
		expected []string
	}{
		{
			testName: "abcac",
			values:   []string{"a", "b", "a", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			testName: "aabc",
			values:   []string{"a", "a", "b", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			testName: "abcc",
			values:   []string{"a", "b", "c", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			testName: "aabcc",
			values:   []string{"a", "a", "b", "c", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			testName: "aabbcc",
			values:   []string{"a", "a", "b", "b", "c", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			testName: "empty",
			values:   []string{},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			if got, want := removeDuplicates(tt.values), tt.expected; !slices.Equal(got, want) {
				t.Errorf("got %v, want %v", got, want)
			}
		})
	}
}

func TestValuesNotInMap(t *testing.T) {
	tests := []struct {
		testName  string
		values    []string
		keysInMap []string
		expected  []string
	}{
		{
			testName:  "abcac-ab",
			values:    []string{"a", "b", "a", "c"},
			keysInMap: []string{"a", "b"},
			expected:  []string{"c"},
		},
		{
			testName:  "abc-a",
			values:    []string{"a", "b", "c"},
			keysInMap: []string{"a"},
			expected:  []string{"b", "c"},
		},
		{
			testName:  "abc-empty",
			values:    []string{"a", "b", "c"},
			keysInMap: []string{},
			expected:  []string{"a", "b", "c"},
		},
		{
			testName:  "empty-a",
			values:    []string{},
			keysInMap: []string{"a"},
			expected:  []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			mapToCheck := make(map[string]struct{})
			for _, keys := range tt.keysInMap {
				mapToCheck[keys] = struct{}{}
			}
			if got, want := valuesNotInMap(tt.values, mapToCheck), tt.expected; !slices.Equal(got, want) {
				t.Errorf("got %v, want %v", got, want)
			}
		})
	}
}

func TestRemoveEmptyString(t *testing.T) {
	tests := []struct {
		testName string
		values   []string
		expected []string
	}{
		{
			testName: "ab_c",
			values:   []string{"a", "b", "", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			testName: "abc_",
			values:   []string{"a", "b", "c", ""},
			expected: []string{"a", "b", "c"},
		},
		{
			testName: "_abc",
			values:   []string{"", "a", "b", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			testName: "_",
			values:   []string{""},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			if got, want := removeEmptyString(tt.values), tt.expected; !slices.Equal(got, want) {
				t.Errorf("got %v, want %v", got, want)
			}
		})
	}
}
