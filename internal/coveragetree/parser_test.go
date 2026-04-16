package coveragetree_test

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/romantelychko/coverage-tree/internal/coveragetree"
)

// testdataPath повертає шлях до директорії testdata відносно кореня проєкту.
func testdataPath(fileName string) string {
	pcs := [1]uintptr{}
	runtime.Callers(2, pcs[:])
	frame, _ := runtime.CallersFrames(pcs[:]).Next()
	currentFile := frame.File
	projectRoot := filepath.Join(filepath.Dir(currentFile), "..", "..")

	return filepath.Join(projectRoot, "testdata", fileName)
}

// testModulePrefix — тестовий префікс модуля, що повторюється у кількох тестах.
const testModulePrefix = "github.com/myorg/myproject/"

// TestParseCoverageSimple перевіряє базовий парсинг coverage.out.
func TestParseCoverageSimple(test *testing.T) {
	test.Parallel()

	inputPath := testdataPath("simple.out")
	prefix := testModulePrefix
	suffixes := coveragetree.DefaultExcludeSuffixes()
	directories := coveragetree.DefaultExcludeDirs()

	fileStats, err := coveragetree.ParseCoverage(inputPath, prefix, suffixes, directories)
	if err != nil {
		test.Fatalf("ParseCoverage повернув помилку: %v", err)
	}

	// Перевіряємо кількість файлів (без _mock.go, _test.go, mocks/)
	expectedFileCount := 4
	if len(fileStats) != expectedFileCount {
		test.Errorf("Очікувалось %d файлів, отримано %d", expectedFileCount, len(fileStats))
	}

	// Перевірка конкретного файлу
	handlerStats, exists := fileStats["internal/users/handler.go"]
	if !exists {
		test.Fatal("Файл internal/users/handler.go не знайдено у результатах")
	}

	// handler.go: блок 10.30,12.2 (1 stmt, covered) + блок 14.40,18.2 (3 stmts, not covered)
	if handlerStats.Statements != 4 {
		test.Errorf("handler.go: очікувалось 4 statements, отримано %d", handlerStats.Statements)
	}

	if handlerStats.Covered != 1 {
		test.Errorf("handler.go: очікувалось 1 covered, отримано %d", handlerStats.Covered)
	}
}

// TestParseCoverageDeduplicate перевіряє дедуплікацію блоків при -coverpkg=./...
func TestParseCoverageDeduplicate(test *testing.T) {
	test.Parallel()

	inputPath := testdataPath("coverpkg.out")
	prefix := testModulePrefix
	suffixes := coveragetree.DefaultExcludeSuffixes()
	directories := coveragetree.DefaultExcludeDirs()

	fileStats, err := coveragetree.ParseCoverage(inputPath, prefix, suffixes, directories)
	if err != nil {
		test.Fatalf("ParseCoverage повернув помилку: %v", err)
	}

	// Має бути тільки 1 файл (service.go) — mock та mocks/ виключені
	if len(fileStats) != 1 {
		test.Errorf("Очікувався 1 файл, отримано %d", len(fileStats))
	}

	serviceStats, exists := fileStats["internal/users/service.go"]
	if !exists {
		test.Fatal("Файл internal/users/service.go не знайдено")
	}

	// Два блоки: 10.30,15.2 (4 stmts) та 17.40,20.2 (2 stmts)
	// Для mode:set — max(1,0)=1 для першого, max(0,1)=1 для другого
	// Тому обидва covered
	if serviceStats.Statements != 6 {
		test.Errorf("service.go: очікувалось 6 statements, отримано %d", serviceStats.Statements)
	}

	if serviceStats.Covered != 6 {
		test.Errorf("service.go: очікувалось 6 covered (дедуплікація max), отримано %d", serviceStats.Covered)
	}
}

// TestParseCoverageAtomicMode перевіряє сумування count для mode:atomic.
func TestParseCoverageAtomicMode(test *testing.T) {
	test.Parallel()

	inputPath := testdataPath("atomic.out")
	prefix := testModulePrefix

	fileStats, err := coveragetree.ParseCoverage(inputPath, prefix, nil, nil)
	if err != nil {
		test.Fatalf("ParseCoverage повернув помилку: %v", err)
	}

	serviceStats, exists := fileStats["internal/users/service.go"]
	if !exists {
		test.Fatal("Файл internal/users/service.go не знайдено")
	}

	// Блок 10.30,15.2: count = 3 + 5 = 8 (atomic сумує) — covered
	// Блок 17.40,20.2: count = 0 + 0 = 0 — not covered
	if serviceStats.Statements != 6 {
		test.Errorf("service.go: очікувалось 6 statements, отримано %d", serviceStats.Statements)
	}

	if serviceStats.Covered != 4 {
		test.Errorf("service.go: очікувалось 4 covered, отримано %d", serviceStats.Covered)
	}
}

// TestParseCoverageEmpty перевіряє парсинг порожнього coverage.out.
func TestParseCoverageEmpty(test *testing.T) {
	test.Parallel()

	inputPath := testdataPath("empty.out")

	fileStats, err := coveragetree.ParseCoverage(inputPath, "", nil, nil)
	if err != nil {
		test.Fatalf("ParseCoverage повернув помилку: %v", err)
	}

	if len(fileStats) != 0 {
		test.Errorf("Очікувалось 0 файлів, отримано %d", len(fileStats))
	}
}

// TestParseCoverageFileNotFound перевіряє помилку при відсутньому файлі.
func TestParseCoverageFileNotFound(test *testing.T) {
	test.Parallel()

	_, err := coveragetree.ParseCoverage("/nonexistent/coverage.out", "", nil, nil)
	if err == nil {
		test.Error("ParseCoverage має повертати помилку для неіснуючого файлу")
	}
}

// TestDetectModulePrefix перевіряє автодетект префікса модуля.
func TestDetectModulePrefix(test *testing.T) {
	test.Parallel()

	inputPath := testdataPath("simple.out")

	prefix, err := coveragetree.DetectModulePrefix(inputPath)
	if err != nil {
		test.Fatalf("DetectModulePrefix повернув помилку: %v", err)
	}

	expectedPrefix := testModulePrefix
	if prefix != expectedPrefix {
		test.Errorf("Очікувався префікс %q, отримано %q", expectedPrefix, prefix)
	}
}

// TestDetectModulePrefixEmpty перевіряє автодетект для порожнього файлу.
func TestDetectModulePrefixEmpty(test *testing.T) {
	test.Parallel()

	inputPath := testdataPath("empty.out")

	prefix, err := coveragetree.DetectModulePrefix(inputPath)
	if err != nil {
		test.Fatalf("DetectModulePrefix повернув помилку: %v", err)
	}

	if prefix != "" {
		test.Errorf("Очікувався порожній префікс, отримано %q", prefix)
	}
}

// TestParseCoverageMixedLines перевіряє пропуск порожніх та невалідних рядків.
func TestParseCoverageMixedLines(test *testing.T) {
	test.Parallel()

	inputPath := testdataPath("mixed-lines.out")
	prefix := testModulePrefix

	fileStats, err := coveragetree.ParseCoverage(inputPath, prefix, nil, nil)
	if err != nil {
		test.Fatalf("ParseCoverage повернув помилку: %v", err)
	}

	if len(fileStats) != 2 {
		test.Errorf("Очікувалось 2 файли, отримано %d", len(fileStats))
	}
}

// TestDetectModulePrefixFileNotFound перевіряє помилку при відсутньому файлі.
func TestDetectModulePrefixFileNotFound(test *testing.T) {
	test.Parallel()

	_, err := coveragetree.DetectModulePrefix("/nonexistent/coverage.out")
	if err == nil {
		test.Error("DetectModulePrefix має повертати помилку для неіснуючого файлу")
	}
}

// TestDetectModulePrefixNoColon перевіряє пропуск рядків без двокрапки.
func TestDetectModulePrefixNoColon(test *testing.T) {
	test.Parallel()

	inputPath := testdataPath("no-colon-lines.out")

	prefix, err := coveragetree.DetectModulePrefix(inputPath)
	if err != nil {
		test.Fatalf("DetectModulePrefix повернув помилку: %v", err)
	}

	expectedPrefix := testModulePrefix
	if prefix != expectedPrefix {
		test.Errorf("Очікувався префікс %q, отримано %q", expectedPrefix, prefix)
	}
}

// TestDetectModulePrefixNoModulePrefix перевіряє повернення порожнього префікса
// коли відома директорія стоїть на початку шляху (без модульного префікса).
func TestDetectModulePrefixNoModulePrefix(test *testing.T) {
	test.Parallel()

	inputPath := testdataPath("no-prefix.out")

	prefix, err := coveragetree.DetectModulePrefix(inputPath)
	if err != nil {
		test.Fatalf("DetectModulePrefix повернув помилку: %v", err)
	}

	if prefix != "" {
		test.Errorf("Очікувався порожній префікс, отримано %q", prefix)
	}
}
