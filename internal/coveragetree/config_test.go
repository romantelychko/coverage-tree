package coveragetree_test

import (
	"testing"

	"github.com/romantelychko/coverage-tree/internal/coveragetree"
)

// TestDefaultConfig перевіряє конфігурацію за замовчуванням.
func TestDefaultConfig(test *testing.T) {
	test.Parallel()

	config := coveragetree.DefaultConfig()

	if config.Theme != "dark" {
		test.Errorf("Очікувалась тема 'dark', отримано %q", config.Theme)
	}

	if config.Language != "uk" {
		test.Errorf("Очікувалась мова 'uk', отримано %q", config.Language)
	}

	if len(config.ExcludeSuffixes) == 0 {
		test.Error("ExcludeSuffixes не повинні бути порожніми")
	}

	if len(config.ExcludeDirs) == 0 {
		test.Error("ExcludeDirs не повинні бути порожніми")
	}
}
