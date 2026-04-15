// Package coveragetree генерує HTML-звіт покриття Go-коду у вигляді дерева директорій.
package coveragetree

// Config містить конфігурацію для генерації звіту покриття.
type Config struct {
	// InputPath — шлях до файлу coverage.out
	InputPath string
	// OutputPath — шлях до вихідного HTML-файлу
	OutputPath string
	// ModulePrefix — префікс модуля для обрізки шляхів (порожній = автодетект)
	ModulePrefix string
	// NoAutodetect — вимкнути автодетект префікса модуля
	NoAutodetect bool
	// ExcludeSuffixes — суфікси файлів для виключення
	ExcludeSuffixes []string
	// ExcludeDirs — директорії для виключення
	ExcludeDirs []string
	// Theme — колірна тема HTML-звіту (dark або light)
	Theme string
	// Language — мова інтерфейсу HTML-звіту (uk або en)
	Language string
	// Title — заголовок HTML-звіту (порожній = визначається мовою)
	Title string
}

// DefaultExcludeSuffixes повертає суфікси файлів, що виключаються за замовчуванням.
func DefaultExcludeSuffixes() []string {
	return []string{"_mock.go", "_test.go"}
}

// DefaultExcludeDirs повертає директорії, що виключаються за замовчуванням.
func DefaultExcludeDirs() []string {
	return []string{"mocks"}
}

// DefaultConfig повертає конфігурацію за замовчуванням.
func DefaultConfig() Config {
	return Config{
		ExcludeSuffixes: DefaultExcludeSuffixes(),
		ExcludeDirs:     DefaultExcludeDirs(),
		Theme:           "dark",
		Language:        "uk",
	}
}
