package coveragetree

import "slices"

import "strings"

// ShouldExclude перевіряє чи файл або директорія мають бути виключені зі звіту.
// Файл виключається якщо його ім'я закінчується на один із суфіксів,
// або якщо будь-яка батьківська директорія є у списку виключених.
func ShouldExclude(filePath string, excludeSuffixes []string, excludeDirs []string) bool {
	// Перевірка суфіксів імені файлу
	for _, suffix := range excludeSuffixes {
		if strings.HasSuffix(filePath, suffix) {
			return true
		}
	}

	// Перевірка директорій у шляху (без останнього компонента — це файл)
	parts := strings.Split(filePath, "/")
	directoryParts := parts[:len(parts)-1]

	for _, directory := range directoryParts {
		if slices.Contains(excludeDirs, directory) {
			return true
		}
	}

	return false
}
