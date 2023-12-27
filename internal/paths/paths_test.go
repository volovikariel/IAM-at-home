package paths

import (
	"reflect"
	"testing"
	"testing/fstest"
)

func TestGetInputPaths(t *testing.T) {
	infs := fstest.MapFS{
		"file.d2":   {},
		"file.x.d2": {},
		"file.d2d2": {},
		"file.go":   {},
		"file":      {},
		// 1 level deep
		"foo/":          {},
		"foo/file.d2":   {},
		"foo/file.x.d2": {},
		"foo/file.d2d2": {},
		"foo/file.go":   {},
		"foo/file":      {},
		// 2 levels deep
		"foo/d2/":          {},
		"foo/d2/file.d2":   {},
		"foo/d2/file.x.d2": {},
		"foo/d2/file.d2d2": {},
		"foo/d2/file.go":   {},
		"foo/d2/file":      {},
	}
	testCases := []struct {
		name     string
		recurse  bool
		expected []string
	}{
		{
			name:    "with recursion",
			recurse: true,
			expected: []string{
				"file.d2",
				"file.x.d2",
				"foo/d2/file.d2",
				"foo/d2/file.x.d2",
				"foo/file.d2",
				"foo/file.x.d2",
			},
		},
		{
			name:    "without recursion",
			recurse: false,
			expected: []string{
				"file.d2",
				"file.x.d2",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := GetInputPaths(infs, ".d2", tc.recurse)
			if !reflect.DeepEqual(tc.expected, got) {
				t.Errorf("\nWant:\n%q\nGot:\n%q\n", tc.expected, got)
			}
		})
	}
}
