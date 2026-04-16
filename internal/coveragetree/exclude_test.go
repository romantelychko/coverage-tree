package coveragetree_test

import (
	"testing"

	"github.com/romantelychko/coverage-tree/internal/coveragetree"
)

// TestShouldExcludeBySuffix перевіряє виключення файлів за суфіксом.
func TestShouldExcludeBySuffix(test *testing.T) {
	test.Parallel()

	suffixes := []string{"_mock.go", "_test.go"}
	directories := []string{"mocks"}

	testCases := []struct {
		name     string
		filePath string
		expected bool
	}{
		{
			name:     "звичайний файл не виключається",
			filePath: "internal/users/service.go",
			expected: false,
		},
		{
			name:     "mock файл виключається",
			filePath: "internal/users/service_mock.go",
			expected: true,
		},
		{
			name:     "test файл виключається",
			filePath: "internal/users/service_test.go",
			expected: true,
		},
		{
			name:     "файл у директорії mocks виключається",
			filePath: "internal/users/mocks/repository.go",
			expected: true,
		},
		{
			name:     "файл з mocks у назві але не директорія не виключається",
			filePath: "internal/users/mocks_helper.go",
			expected: false,
		},
		{
			name:     "файл без шляху не виключається",
			filePath: "main.go",
			expected: false,
		},
		{
			name:     "mock файл без шляху виключається",
			filePath: "service_mock.go",
			expected: true,
		},
	}

	for _, testCase := range testCases {
		test.Run(testCase.name, func(test *testing.T) {
			test.Parallel()

			result := coveragetree.ShouldExclude(testCase.filePath, suffixes, directories)
			if result != testCase.expected {
				test.Errorf("ShouldExclude(%q) = %v, очікувалось %v", testCase.filePath, result, testCase.expected)
			}
		})
	}
}

// TestShouldExcludeWithEmptyFilters перевіряє що без фільтрів нічого не виключається.
func TestShouldExcludeWithEmptyFilters(test *testing.T) {
	test.Parallel()

	result := coveragetree.ShouldExclude("internal/users/service_mock.go", nil, nil)
	if result {
		test.Error("ShouldExclude з порожніми фільтрами має повертати false")
	}
}
