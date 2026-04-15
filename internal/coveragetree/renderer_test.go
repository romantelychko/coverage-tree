package coveragetree_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/romantelychko/coverage-tree/internal/coveragetree"
)

// renderAndReadHTML викликає RenderHTML та повертає вміст згенерованого файлу.
func renderAndReadHTML(
	t *testing.T,
	treeJSON *coveragetree.TreeJSON,
	config coveragetree.Config,
) string {
	t.Helper()

	err := coveragetree.RenderHTML(treeJSON, config)
	if err != nil {
		t.Fatalf("RenderHTML повернув помилку: %v", err)
	}

	content, err := os.ReadFile(config.OutputPath)
	if err != nil {
		t.Fatalf("Не вдалося прочитати файл: %v", err)
	}

	return string(content)
}

// TestRenderHTMLDarkTheme перевіряє генерацію HTML зі стандартною темною темою.
func TestRenderHTMLDarkTheme(t *testing.T) {
	treeJSON := &coveragetree.TreeJSON{
		Name:       "root",
		Statements: 100,
		Covered:    75,
		Coverage:   75.0,
		Children:   make([]*coveragetree.TreeJSON, 0),
		Files:      make([]*coveragetree.FileJSON, 0),
	}

	config := coveragetree.Config{
		OutputPath: filepath.Join(t.TempDir(), "report.html"),
		Theme:      "dark",
		Language:   "uk",
	}

	htmlContent := renderAndReadHTML(t, treeJSON, config)

	if !strings.Contains(htmlContent, "Дерево покриття тестами") {
		t.Error("HTML не містить українського заголовка")
	}

	if !strings.Contains(htmlContent, "--bg: #1e1e1e") {
		t.Error("HTML не містить CSS-змінних темної теми")
	}

	if !strings.Contains(htmlContent, "Розгорнути все") {
		t.Error("HTML не містить українського тексту кнопки")
	}

	if !strings.Contains(htmlContent, `lang="uk"`) {
		t.Error("HTML не містить атрибуту lang=uk")
	}
}

// TestRenderHTMLLightThemeEnglish перевіряє генерацію HTML зі світлою темою англійською.
func TestRenderHTMLLightThemeEnglish(t *testing.T) {
	treeJSON := &coveragetree.TreeJSON{
		Name:       "root",
		Statements: 50,
		Covered:    25,
		Coverage:   50.0,
		Children:   make([]*coveragetree.TreeJSON, 0),
		Files:      make([]*coveragetree.FileJSON, 0),
	}

	config := coveragetree.Config{
		OutputPath: filepath.Join(t.TempDir(), "report-en.html"),
		Theme:      "light",
		Language:   "en",
	}

	htmlContent := renderAndReadHTML(t, treeJSON, config)

	if !strings.Contains(htmlContent, "Test Coverage Tree") {
		t.Error("HTML не містить англійського заголовка")
	}

	if !strings.Contains(htmlContent, "--bg: #ffffff") {
		t.Error("HTML не містить CSS-змінних світлої теми")
	}

	if !strings.Contains(htmlContent, "Expand all") {
		t.Error("HTML не містить англійського тексту кнопки")
	}

	if !strings.Contains(htmlContent, `lang="en"`) {
		t.Error("HTML не містить атрибуту lang=en")
	}
}

// TestRenderHTMLCustomTitle перевіряє генерацію HTML з кастомним заголовком.
func TestRenderHTMLCustomTitle(t *testing.T) {
	treeJSON := &coveragetree.TreeJSON{
		Name:       "root",
		Statements: 10,
		Covered:    10,
		Coverage:   100.0,
		Children:   make([]*coveragetree.TreeJSON, 0),
		Files:      make([]*coveragetree.FileJSON, 0),
	}

	config := coveragetree.Config{
		OutputPath: filepath.Join(t.TempDir(), "report-custom.html"),
		Theme:      "dark",
		Language:   "uk",
		Title:      "API Coverage Report",
	}

	htmlContent := renderAndReadHTML(t, treeJSON, config)

	if !strings.Contains(htmlContent, "API Coverage Report") {
		t.Error("HTML не містить кастомного заголовка")
	}
}

// TestRenderHTMLInvalidTheme перевіряє помилку при невідомій темі.
func TestRenderHTMLInvalidTheme(t *testing.T) {
	treeJSON := &coveragetree.TreeJSON{
		Name:       "root",
		Statements: 10,
		Covered:    5,
		Coverage:   50.0,
		Children:   make([]*coveragetree.TreeJSON, 0),
		Files:      make([]*coveragetree.FileJSON, 0),
	}

	config := coveragetree.Config{
		OutputPath: filepath.Join(t.TempDir(), "report.html"),
		Theme:      "neon",
		Language:   "uk",
	}

	err := coveragetree.RenderHTML(treeJSON, config)
	if err == nil {
		t.Error("RenderHTML має повертати помилку для невідомої теми")
	}
}

// TestRenderHTMLInvalidLanguage перевіряє помилку при невідомій мові.
func TestRenderHTMLInvalidLanguage(t *testing.T) {
	treeJSON := &coveragetree.TreeJSON{
		Name:       "root",
		Statements: 10,
		Covered:    5,
		Coverage:   50.0,
		Children:   make([]*coveragetree.TreeJSON, 0),
		Files:      make([]*coveragetree.FileJSON, 0),
	}

	config := coveragetree.Config{
		OutputPath: filepath.Join(t.TempDir(), "report.html"),
		Theme:      "dark",
		Language:   "de",
	}

	err := coveragetree.RenderHTML(treeJSON, config)
	if err == nil {
		t.Error("RenderHTML має повертати помилку для невідомої мови")
	}
}

// TestRenderHTMLContainsTreeData перевіряє що JSON-дані дерева присутні в HTML.
func TestRenderHTMLContainsTreeData(t *testing.T) {
	treeJSON := &coveragetree.TreeJSON{
		Name:       "root",
		Statements: 100,
		Covered:    73,
		Coverage:   73.0,
		Children: []*coveragetree.TreeJSON{
			{
				Name:       "internal",
				Statements: 80,
				Covered:    60,
				Coverage:   75.0,
				Children:   make([]*coveragetree.TreeJSON, 0),
				Files: []*coveragetree.FileJSON{
					{Name: "service.go", Statements: 80, Covered: 60, Coverage: 75.0},
				},
			},
		},
		Files: []*coveragetree.FileJSON{
			{Name: "main.go", Statements: 20, Covered: 13, Coverage: 65.0},
		},
	}

	config := coveragetree.Config{
		OutputPath: filepath.Join(t.TempDir(), "report-data.html"),
		Theme:      "dark",
		Language:   "uk",
	}

	htmlContent := renderAndReadHTML(t, treeJSON, config)

	if !strings.Contains(htmlContent, `"service.go"`) {
		t.Error("HTML не містить назви файлу з JSON-даних")
	}

	if !strings.Contains(htmlContent, `"internal"`) {
		t.Error("HTML не містить назви директорії з JSON-даних")
	}
}

// TestRenderHTMLWriteError перевіряє помилку при неможливості записати файл.
func TestRenderHTMLWriteError(t *testing.T) {
	treeJSON := &coveragetree.TreeJSON{
		Name:       "root",
		Statements: 10,
		Covered:    5,
		Coverage:   50.0,
		Children:   make([]*coveragetree.TreeJSON, 0),
		Files:      make([]*coveragetree.FileJSON, 0),
	}

	config := coveragetree.Config{
		OutputPath: filepath.Join("/nonexistent/directory", "report.html"),
		Theme:      "dark",
		Language:   "uk",
	}

	err := coveragetree.RenderHTML(treeJSON, config)
	if err == nil {
		t.Error("RenderHTML має повертати помилку при неможливості записати файл")
	}
}
