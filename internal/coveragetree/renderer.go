package coveragetree

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"text/template"
)

// filePermissions — права доступу до вихідного HTML-файлу.
const filePermissions = 0o644

// ErrUnknownTheme повертається коли передано невідому тему.
var ErrUnknownTheme = errors.New("невідома тема")

// ErrUnknownLanguage повертається коли передано невідому мову.
var ErrUnknownLanguage = errors.New("невідома мова")

// themeVariables містить CSS-змінні для кожної теми.
var themeVariables = map[string]string{
	"dark": `--bg: #1e1e1e;
      --bh: #2a2d2e;
      --bf: #2a2d2e;
      --t: #cccccc;
      --td: #858585;
      --bd: #404040;
      --g: #4ec94e;
      --y: #e6a817;
      --r: #f44747;
      --bb: #3c3c3c;
      --a: #d4d4d4;
      --btn: #3c3c3c;
      --bar-g: #4ec94e;
      --bar-y: #e6a817;
      --bar-r: #f44747`,
	"light": `--bg: #ffffff;
      --bh: #f0f0f0;
      --bf: #f5f5f5;
      --t: #333333;
      --td: #777777;
      --bd: #e0e0e0;
      --g: #2d8f2d;
      --y: #b8860b;
      --r: #d32f2f;
      --bb: #e8e8e8;
      --a: #1a1a1a;
      --btn: #f0f0f0;
      --bar-g: #2d8f2d;
      --bar-y: #b8860b;
      --bar-r: #d32f2f`,
}

// localization містить рядки інтерфейсу для кожної мови.
type localization struct {
	Title          string
	SummaryLabel   string
	ExpandAll      string
	CollapseAll    string
	HeaderName     string
	HeaderCoverage string
	HeaderLines    string
}

// localizations — словник локалізацій.
var localizations = map[string]localization{
	"uk": {
		Title:          "Дерево покриття тестами",
		SummaryLabel:   "Загальне покриття",
		ExpandAll:      "Розгорнути все",
		CollapseAll:    "Згорнути все",
		HeaderName:     "Назва",
		HeaderCoverage: "Покриття",
		HeaderLines:    "Рядки",
	},
	"en": {
		Title:          "Test Coverage Tree",
		SummaryLabel:   "Total coverage",
		ExpandAll:      "Expand all",
		CollapseAll:    "Collapse all",
		HeaderName:     "Name",
		HeaderCoverage: "Coverage",
		HeaderLines:    "Lines",
	},
}

// templateData містить дані для підстановки в HTML-шаблон.
type templateData struct {
	Lang           string
	Title          string
	ThemeVars      string
	SummaryLabel   string
	ExpandAll      string
	CollapseAll    string
	HeaderName     string
	HeaderCoverage string
	HeaderLines    string
	TreeData       string
}

// RenderHTML генерує HTML-звіт покриття та записує його у файл.
func RenderHTML(treeJSON *TreeJSON, config Config) error {
	// Серіалізація дерева в JSON
	jsonData, err := json.Marshal(treeJSON)
	if err != nil {
		return fmt.Errorf("Помилка серіалізації дерева в JSON: %w", err)
	}

	// Визначення теми
	themeVars, exists := themeVariables[config.Theme]
	if !exists {
		return fmt.Errorf("%w: %s (доступні: dark, light)", ErrUnknownTheme, config.Theme)
	}

	// Визначення локалізації
	locale, exists := localizations[config.Language]
	if !exists {
		return fmt.Errorf("%w: %s (доступні: uk, en)", ErrUnknownLanguage, config.Language)
	}

	// Заголовок: кастомний або з локалізації
	title := locale.Title
	if config.Title != "" {
		title = config.Title
	}

	// Підготовка даних для шаблону
	data := templateData{
		Lang:           config.Language,
		Title:          title,
		ThemeVars:      themeVars,
		SummaryLabel:   locale.SummaryLabel,
		ExpandAll:      locale.ExpandAll,
		CollapseAll:    locale.CollapseAll,
		HeaderName:     locale.HeaderName,
		HeaderCoverage: locale.HeaderCoverage,
		HeaderLines:    locale.HeaderLines,
		TreeData:       string(jsonData),
	}

	// Парсинг та виконання шаблону
	htmlTemplate, err := template.New("coverage").Parse(GetTemplate())
	if err != nil {
		return fmt.Errorf("Помилка парсингу HTML-шаблону: %w", err)
	}

	var builder strings.Builder
	if err := htmlTemplate.Execute(&builder, data); err != nil {
		return fmt.Errorf("Помилка генерації HTML: %w", err)
	}

	// Запис у файл
	if err := os.WriteFile(config.OutputPath, []byte(builder.String()), filePermissions); err != nil {
		return fmt.Errorf("Помилка запису файлу %s: %w", config.OutputPath, err)
	}

	return nil
}
