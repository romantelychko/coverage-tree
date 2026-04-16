// Генерація HTML-звіту покриття Go-коду у вигляді дерева директорій.
// Парсить coverage.out та генерує інтерактивний HTML з collapsible tree.
//
// Використання: coverage-tree [flags] <coverage.out> <output.html>
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/romantelychko/coverage-tree/internal/coveragetree"
)

// version встановлюється через -ldflags при білді.
var version = "dev"

// requiredArgCount — кількість обов'язкових позиційних аргументів.
const requiredArgCount = 2

// stringSliceFlag реалізує flag.Value для прапорців, що можуть бути вказані кілька разів.
type stringSliceFlag []string

// String повертає рядкове представлення значень прапорця.
func (flag *stringSliceFlag) String() string {
	return strings.Join(*flag, ", ")
}

// Set додає нове значення до списку.
func (flag *stringSliceFlag) Set(value string) error {
	*flag = append(*flag, value)

	return nil
}

// normalizeBoolArgs перетворює ["--flag", "true"] на ["--flag=true"] для відомих булевих прапорців,
// щоб запобігти потраплянню "true"/"false" у позиційні аргументи.
func normalizeBoolArgs(args []string, boolFlags map[string]bool) []string {
	result := make([]string, 0, len(args))

	for i := 0; i < len(args); i++ {
		arg := args[i]

		name := ""

		switch {
		case strings.HasPrefix(arg, "--"):
			name = strings.TrimPrefix(arg, "--")
		case strings.HasPrefix(arg, "-"):
			name = strings.TrimPrefix(arg, "-")
		}

		// Якщо це булевий прапорець без вбудованого значення і наступний токен — "true"/"false"
		if name != "" && !strings.Contains(name, "=") && boolFlags[name] {
			if i+1 < len(args) && (args[i+1] == "true" || args[i+1] == "false") {
				result = append(result, arg+"="+args[i+1])
				i++

				continue
			}
		}

		result = append(result, arg)
	}

	return result
}

// printDefaults виводить список прапорців із префіксом "--" замість стандартного "-".
func printDefaults() {
	flag.VisitAll(func(f *flag.Flag) {
		typeName, usage := flag.UnquoteUsage(f)

		s := fmt.Sprintf("  --%s", f.Name)
		if typeName != "" {
			s += " " + typeName
		}

		s += "\n    \t" + usage

		if f.DefValue != "" && f.DefValue != "false" {
			s += fmt.Sprintf(" (default %q)", f.DefValue)
		}

		fmt.Fprintln(os.Stderr, s)
	})
}

// runParams містить усі параметри для запуску генерації звіту.
type runParams struct {
	inputPath       string
	outputPath      string
	modulePrefix    string
	noAutodetect    bool
	excludeSuffixes []string
	excludeDirs     []string
	theme           string
	language        string
	title           string
}

// runCoverageTree виконує повний цикл: парсинг → побудова дерева → генерація HTML.
func runCoverageTree(p runParams) error {
	if _, err := os.Stat(p.inputPath); os.IsNotExist(err) {
		return fmt.Errorf("файл не знайдено: %s", p.inputPath)
	}

	if len(p.excludeSuffixes) == 0 {
		p.excludeSuffixes = coveragetree.DefaultExcludeSuffixes()
	}

	if len(p.excludeDirs) == 0 {
		p.excludeDirs = coveragetree.DefaultExcludeDirs()
	}

	resolvedPrefix, err := resolveModulePrefix(p.inputPath, p.modulePrefix, p.noAutodetect)
	if err != nil {
		return fmt.Errorf("помилка автодетекту префікса: %w", err)
	}

	fileStats, err := coveragetree.ParseCoverage(
		p.inputPath, resolvedPrefix, p.excludeSuffixes, p.excludeDirs,
	)
	if err != nil {
		return fmt.Errorf("помилка парсингу: %w", err)
	}

	if len(fileStats) == 0 {
		return fmt.Errorf("дані покриття не знайдено")
	}

	fmt.Printf("Файлів: %d\n", len(fileStats))

	tree := coveragetree.BuildTree(fileStats)
	coveragetree.Aggregate(tree)
	treeJSON := coveragetree.ToJSON(tree, "root")

	config := coveragetree.Config{
		InputPath:       p.inputPath,
		OutputPath:      p.outputPath,
		ModulePrefix:    resolvedPrefix,
		ExcludeSuffixes: p.excludeSuffixes,
		ExcludeDirs:     p.excludeDirs,
		Theme:           p.theme,
		Language:        p.language,
		Title:           p.title,
	}

	if err := coveragetree.RenderHTML(treeJSON, config); err != nil {
		return fmt.Errorf("помилка генерації HTML: %w", err)
	}

	fmt.Printf("Покриття: %.1f%% (%d / %d)\n", treeJSON.Coverage, treeJSON.Covered, treeJSON.Statements)
	fmt.Printf("Звіт: %s\n", p.outputPath)

	return nil
}

// resolveModulePrefix визначає префікс модуля: з прапорця або автодетектом.
func resolveModulePrefix(inputPath, modulePrefix string, noAutodetect bool) (string, error) {
	if modulePrefix != "" || noAutodetect {
		return modulePrefix, nil
	}

	detected, err := coveragetree.DetectModulePrefix(inputPath)
	if err != nil {
		return "", err
	}

	if detected != "" {
		fmt.Printf("Префікс модуля: %s\n", detected)
	}

	return detected, nil
}

func main() {
	modulePrefix := flag.String("module-prefix", "", "Префікс модуля для обрізки шляхів (за замовчуванням: автодетект)")
	theme := flag.String("theme", "dark", "Колірна тема HTML-звіту: dark або light")
	language := flag.String("lang", "uk", "Мова інтерфейсу HTML-звіту: uk або en")
	title := flag.String("title", "", "Заголовок HTML-звіту (за замовчуванням: визначається мовою)")
	noAutodetect := flag.Bool("no-autodetect", false, "Вимкнути автодетект префікса модуля")
	showVersion := flag.Bool("version", false, "Вивести версію та вийти")

	var excludeSuffixes stringSliceFlag

	var excludeDirs stringSliceFlag

	const excludeSuffixUsage = "Суфікс файлів для виключення" +
		" (можна вказати кілька разів, за замовчуванням: _mock.go, _test.go)"

	const excludeDirUsage = "Директорія для виключення" +
		" (можна вказати кілька разів, за замовчуванням: mocks)"

	flag.Var(&excludeSuffixes, "exclude-suffix", excludeSuffixUsage)
	flag.Var(&excludeDirs, "exclude-dir", excludeDirUsage)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "CoverageTree v%s - генерація HTML-звіту покриття Go-коду\n\n", version)
		fmt.Fprintf(os.Stderr, "Використання: coverage-tree [flags] <coverage.out> <output.html>\n\n")
		fmt.Fprintf(os.Stderr, "Аргументи:\n")
		fmt.Fprintf(os.Stderr, "  coverage.out     шлях до файлу coverage profile\n")
		fmt.Fprintf(os.Stderr, "  output.html      шлях до вихідного HTML-файлу\n\n")
		fmt.Fprintf(os.Stderr, "Прапорці:\n")
		printDefaults()
	}

	boolFlags := map[string]bool{
		"no-autodetect": true,
		"version":       true,
	}

	os.Args = append(os.Args[:1], normalizeBoolArgs(os.Args[1:], boolFlags)...)

	flag.Parse()

	if *showVersion {
		fmt.Printf("coverage-tree v%s\n", version)
		os.Exit(0)
	}

	positionalArgs := flag.Args()
	if len(positionalArgs) != requiredArgCount {
		fmt.Fprintf(os.Stderr, "Помилка: потрібно вказати два аргументи: <coverage.out> <output.html>\n\n")
		flag.Usage()
		os.Exit(1)
	}

	err := runCoverageTree(runParams{
		inputPath:       positionalArgs[0],
		outputPath:      positionalArgs[1],
		modulePrefix:    *modulePrefix,
		noAutodetect:    *noAutodetect,
		excludeSuffixes: excludeSuffixes,
		excludeDirs:     excludeDirs,
		theme:           *theme,
		language:        *language,
		title:           *title,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Помилка: %v\n", err)
		os.Exit(1)
	}
}
