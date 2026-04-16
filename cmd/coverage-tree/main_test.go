package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// projectTestdata повертає шлях до директорії testdata у корені проєкту.
func projectTestdata(fileName string) string {
	pcs := [1]uintptr{}
	runtime.Callers(2, pcs[:])
	frame, _ := runtime.CallersFrames(pcs[:]).Next()
	currentFile := frame.File

	return filepath.Join(filepath.Dir(currentFile), "..", "..", "testdata", fileName)
}

// TestStringSliceFlagString перевіряє рядкове представлення прапорця.
func TestStringSliceFlagString(test *testing.T) {
	var flag stringSliceFlag

	if flag.String() != "" {
		test.Errorf("Очікувався порожній рядок, отримано %q", flag.String())
	}

	_ = flag.Set("a")
	_ = flag.Set("b")

	result := flag.String()
	if result != "a, b" {
		test.Errorf("Очікувався рядок %q, отримано %q", "a, b", result)
	}
}

// TestStringSliceFlagSet перевіряє додавання значень до прапорця.
func TestStringSliceFlagSet(test *testing.T) {
	var flag stringSliceFlag

	if err := flag.Set("first"); err != nil {
		test.Fatalf("Set повернув помилку: %v", err)
	}

	if err := flag.Set("second"); err != nil {
		test.Fatalf("Set повернув помилку: %v", err)
	}

	if len(flag) != 2 {
		test.Errorf("Очікувалось 2 елементи, отримано %d", len(flag))
	}

	if flag[0] != "first" || flag[1] != "second" {
		test.Errorf("Неочікувані значення: %v", []string(flag))
	}
}

// TestRunCoverageTreeSuccess перевіряє успішну генерацію звіту.
func TestRunCoverageTreeSuccess(test *testing.T) {
	params := runParams{
		inputPath:       projectTestdata("simple.out"),
		outputPath:      filepath.Join(test.TempDir(), "out.html"),
		modulePrefix:    "github.com/myorg/myproject/",
		noAutodetect:    true,
		excludeSuffixes: []string{"_test.go"},
		excludeDirs:     []string{"mocks"},
		theme:           "dark",
		language:        "uk",
	}

	if err := runCoverageTree(params); err != nil {
		test.Fatalf("runCoverageTree повернув помилку: %v", err)
	}
}

// TestRunCoverageTreeDefaultExcludes перевіряє заповнення exclude за замовчуванням.
func TestRunCoverageTreeDefaultExcludes(test *testing.T) {
	params := runParams{
		inputPath:    projectTestdata("simple.out"),
		outputPath:   filepath.Join(test.TempDir(), "out.html"),
		modulePrefix: "github.com/myorg/myproject/",
		noAutodetect: true,
		theme:        "dark",
		language:     "uk",
	}

	if err := runCoverageTree(params); err != nil {
		test.Fatalf("runCoverageTree повернув помилку: %v", err)
	}
}

// TestRunCoverageTreeAutodetectPrefix перевіряє автодетект префікса модуля.
func TestRunCoverageTreeAutodetectPrefix(test *testing.T) {
	params := runParams{
		inputPath:  projectTestdata("simple.out"),
		outputPath: filepath.Join(test.TempDir(), "out.html"),
		theme:      "dark",
		language:   "uk",
	}

	if err := runCoverageTree(params); err != nil {
		test.Fatalf("runCoverageTree повернув помилку: %v", err)
	}
}

// TestRunCoverageTreeFileNotFound перевіряє помилку при відсутньому файлі.
func TestRunCoverageTreeFileNotFound(test *testing.T) {
	params := runParams{
		inputPath:    "/nonexistent/coverage.out",
		outputPath:   filepath.Join(test.TempDir(), "out.html"),
		noAutodetect: true,
		theme:        "dark",
		language:     "uk",
	}

	err := runCoverageTree(params)
	if err == nil {
		test.Error("runCoverageTree має повертати помилку для неіснуючого файлу")
	}
}

// TestRunCoverageTreeEmptyStats перевіряє помилку при відсутності даних покриття.
func TestRunCoverageTreeEmptyStats(test *testing.T) {
	params := runParams{
		inputPath:    projectTestdata("empty.out"),
		outputPath:   filepath.Join(test.TempDir(), "out.html"),
		noAutodetect: true,
		theme:        "dark",
		language:     "uk",
	}

	err := runCoverageTree(params)
	if err == nil {
		test.Error("runCoverageTree має повертати помилку при порожніх даних покриття")
	}
}

// TestRunCoverageTreeRenderError перевіряє помилку при невалідній темі.
func TestRunCoverageTreeRenderError(test *testing.T) {
	params := runParams{
		inputPath:    projectTestdata("simple.out"),
		outputPath:   filepath.Join(test.TempDir(), "out.html"),
		modulePrefix: "github.com/myorg/myproject/",
		noAutodetect: true,
		theme:        "invalid-theme",
		language:     "uk",
	}

	err := runCoverageTree(params)
	if err == nil {
		test.Error("runCoverageTree має повертати помилку при невалідній темі")
	}
}

// TestResolveModulePrefixExplicit перевіряє повернення явно вказаного префікса.
func TestResolveModulePrefixExplicit(test *testing.T) {
	prefix, err := resolveModulePrefix("anything", "my-prefix/", false)
	if err != nil {
		test.Fatalf("resolveModulePrefix повернув помилку: %v", err)
	}

	if prefix != "my-prefix/" {
		test.Errorf("Очікувався префікс %q, отримано %q", "my-prefix/", prefix)
	}
}

// TestResolveModulePrefixNoAutodetect перевіряє вимкнення автодетекту.
func TestResolveModulePrefixNoAutodetect(test *testing.T) {
	prefix, err := resolveModulePrefix("anything", "", true)
	if err != nil {
		test.Fatalf("resolveModulePrefix повернув помилку: %v", err)
	}

	if prefix != "" {
		test.Errorf("Очікувався порожній префікс, отримано %q", prefix)
	}
}

// TestResolveModulePrefixAutodetect перевіряє автодетект префікса.
func TestResolveModulePrefixAutodetect(test *testing.T) {
	prefix, err := resolveModulePrefix(projectTestdata("simple.out"), "", false)
	if err != nil {
		test.Fatalf("resolveModulePrefix повернув помилку: %v", err)
	}

	expected := "github.com/myorg/myproject/"
	if prefix != expected {
		test.Errorf("Очікувався префікс %q, отримано %q", expected, prefix)
	}
}

// makeUnreadable створює тимчасовий файл із забороненим доступом на читання.
// Повертає шлях до файлу і функцію відновлення прав доступу.
func makeUnreadable(test *testing.T) string {
	test.Helper()

	inputPath := filepath.Join(test.TempDir(), "coverage.out")
	if err := os.WriteFile(inputPath, []byte("mode: set\n"), 0o644); err != nil {
		test.Fatal(err)
	}

	if err := os.Chmod(inputPath, 0o000); err != nil {
		test.Skip("неможливо встановити права доступу: " + err.Error())
	}

	test.Cleanup(func() { _ = os.Chmod(inputPath, 0o644) })

	return inputPath
}

// TestRunCoverageTreeAutodetectError перевіряє помилку автодетекту при нечитабельному файлі.
func TestRunCoverageTreeAutodetectError(test *testing.T) {
	inputPath := makeUnreadable(test)
	params := runParams{
		inputPath:  inputPath,
		outputPath: filepath.Join(test.TempDir(), "out.html"),
		theme:      "dark",
		language:   "uk",
	}

	err := runCoverageTree(params)
	if err == nil {
		test.Error("runCoverageTree має повертати помилку при помилці автодетекту")
	}
}

// TestRunCoverageTreeParseError перевіряє помилку парсингу при нечитабельному файлі.
func TestRunCoverageTreeParseError(test *testing.T) {
	inputPath := makeUnreadable(test)
	params := runParams{
		inputPath:    inputPath,
		outputPath:   filepath.Join(test.TempDir(), "out.html"),
		noAutodetect: true,
		theme:        "dark",
		language:     "uk",
	}

	err := runCoverageTree(params)
	if err == nil {
		test.Error("runCoverageTree має повертати помилку при помилці парсингу")
	}
}

// TestNormalizeBoolArgsMergeTrue перевіряє злиття "--flag true" в "--flag=true".
func TestNormalizeBoolArgsMergeTrue(test *testing.T) {
	boolFlags := map[string]bool{"no-autodetect": true, "version": true}

	args := []string{"--no-autodetect", "true", "coverage.out", "report.html"}
	result := normalizeBoolArgs(args, boolFlags)

	expected := []string{"--no-autodetect=true", "coverage.out", "report.html"}
	if strings.Join(result, "|") != strings.Join(expected, "|") {
		test.Errorf("Очікувався %v, отримано %v", expected, result)
	}
}

// TestNormalizeBoolArgsMergeFalse перевіряє злиття "--flag false" в "--flag=false".
func TestNormalizeBoolArgsMergeFalse(test *testing.T) {
	boolFlags := map[string]bool{"no-autodetect": true, "version": true}

	args := []string{"--no-autodetect", "false"}
	result := normalizeBoolArgs(args, boolFlags)

	expected := []string{"--no-autodetect=false"}
	if strings.Join(result, "|") != strings.Join(expected, "|") {
		test.Errorf("Очікувався %v, отримано %v", expected, result)
	}
}

// TestNormalizeBoolArgsSingleDash перевіряє злиття для прапорця з одним мінусом.
func TestNormalizeBoolArgsSingleDash(test *testing.T) {
	boolFlags := map[string]bool{"no-autodetect": true}

	args := []string{"-no-autodetect", "true"}
	result := normalizeBoolArgs(args, boolFlags)

	expected := []string{"-no-autodetect=true"}
	if strings.Join(result, "|") != strings.Join(expected, "|") {
		test.Errorf("Очікувався %v, отримано %v", expected, result)
	}
}

// TestNormalizeBoolArgsNoNextArg перевіряє що прапорець без наступного аргументу не змінюється.
func TestNormalizeBoolArgsNoNextArg(test *testing.T) {
	boolFlags := map[string]bool{"no-autodetect": true}

	args := []string{"--no-autodetect"}
	result := normalizeBoolArgs(args, boolFlags)

	if strings.Join(result, "|") != "--no-autodetect" {
		test.Errorf("Очікувався [--no-autodetect], отримано %v", result)
	}
}

// TestNormalizeBoolArgsNonBoolUnchanged перевіряє що не-булевий прапорець не змінюється.
func TestNormalizeBoolArgsNonBoolUnchanged(test *testing.T) {
	boolFlags := map[string]bool{"no-autodetect": true}

	args := []string{"--theme", "true"}
	result := normalizeBoolArgs(args, boolFlags)

	expected := []string{"--theme", "true"}
	if strings.Join(result, "|") != strings.Join(expected, "|") {
		test.Errorf("Очікувався %v, отримано %v", expected, result)
	}
}

// TestNormalizeBoolArgsAlreadyMerged перевіряє що "=value" форма не змінюється.
func TestNormalizeBoolArgsAlreadyMerged(test *testing.T) {
	boolFlags := map[string]bool{"no-autodetect": true}

	args := []string{"--no-autodetect=true", "coverage.out", "report.html"}
	result := normalizeBoolArgs(args, boolFlags)

	expected := []string{"--no-autodetect=true", "coverage.out", "report.html"}
	if strings.Join(result, "|") != strings.Join(expected, "|") {
		test.Errorf("Очікувався %v, отримано %v", expected, result)
	}
}

// TestNormalizeBoolArgsEmpty перевіряє обробку порожнього списку аргументів.
func TestNormalizeBoolArgsEmpty(test *testing.T) {
	boolFlags := map[string]bool{"no-autodetect": true}

	result := normalizeBoolArgs([]string{}, boolFlags)

	if len(result) != 0 {
		test.Errorf("Очікувався порожній результат, отримано %v", result)
	}
}

// TestPrintDefaultsDoubleDash перевіряє що printDefaults виводить прапорці з "--".
func TestPrintDefaultsDoubleDash(test *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		test.Fatalf("os.Pipe: %v", err)
	}

	origStderr := os.Stderr
	os.Stderr = w

	printDefaults()

	w.Close()

	os.Stderr = origStderr

	var buf bytes.Buffer

	_, _ = io.Copy(&buf, r)

	output := buf.String()
	if output == "" {
		test.Fatal("printDefaults нічого не вивела")
	}

	for _, line := range strings.Split(output, "\n") {
		// Рядки з прапорцями починаються з "  -"; перевіряємо що не "  -X" (один мінус)
		if strings.HasPrefix(line, "  -") && !strings.HasPrefix(line, "  --") {
			test.Errorf("printDefaults вивела прапорець з одним мінусом: %q", line)
		}
	}
}
