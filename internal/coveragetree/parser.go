package coveragetree

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// coverageLinePattern — регулярний вираз для рядка coverage.out.
var coverageLinePattern = regexp.MustCompile(`^(.+):(\d+\.\d+),(\d+\.\d+)\s+(\d+)\s+(\d+)$`)

// knownRootDirectories — відомі кореневі директорії Go-проєктів для автодетекту префікса.
var knownRootDirectories = map[string]bool{
	"cmd": true, "internal": true, "pkg": true, "shared": true,
	"api": true, "app": true, "lib": true, "src": true,
	"server": true, "service": true, "services": true,
}

// coverageBlock зберігає інформацію про один блок покриття з дедуплікацією.
type coverageBlock struct {
	filePath   string
	statements int
	count      int
}

// FileStats містить статистику покриття для одного файлу.
type FileStats struct {
	// Statements — загальна кількість виразів у файлі
	Statements int
	// Covered — кількість покритих виразів
	Covered int
}

// ParseCoverage парсить файл coverage.out та повертає статистику покриття по файлах.
// З -coverpkg=./... один блок може з'являтись багато разів (по разу на кожен тестовий пакет).
// Для mode:atomic count сумується, для mode:set/count береться максимальне значення.
func ParseCoverage(
	inputPath string,
	modulePrefix string,
	excludeSuffixes []string,
	excludeDirs []string,
) (map[string]*FileStats, error) {
	file, err := os.Open(inputPath) //nolint:gosec // G304: CLI tool reads user-provided files
	if err != nil {
		return nil, fmt.Errorf("Не вдалося відкрити файл %s: %w", inputPath, err)
	}

	defer func() { _ = file.Close() }()

	blocks := make(map[string]*coverageBlock)
	mode := "set"
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Парсинг рядка mode:
		if after, ok := strings.CutPrefix(line, "mode:"); ok {
			mode = strings.TrimSpace(after)

			continue
		}

		if line == "" {
			continue
		}

		match := coverageLinePattern.FindStringSubmatch(line)
		if match == nil {
			continue
		}

		filePath := match[1]
		blockStart := match[2]
		blockEnd := match[3]

		statements, err := strconv.Atoi(match[4])
		if err != nil {
			continue
		}

		count, err := strconv.Atoi(match[5])
		if err != nil {
			continue
		}

		// Обрізка префікса модуля
		if modulePrefix != "" && strings.HasPrefix(filePath, modulePrefix) {
			filePath = strings.TrimPrefix(filePath, modulePrefix)
			filePath = strings.TrimLeft(filePath, "/")
		}

		// Виключення файлів за суфіксами та директоріями
		if ShouldExclude(filePath, excludeSuffixes, excludeDirs) {
			continue
		}

		// Унікальний ключ блоку: файл + координати
		blockKey := filePath + ":" + blockStart + "," + blockEnd

		if existing, exists := blocks[blockKey]; !exists {
			blocks[blockKey] = &coverageBlock{
				filePath:   filePath,
				statements: statements,
				count:      count,
			}
		} else {
			mergeBlock(existing, count, mode)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("Помилка читання файлу %s: %w", inputPath, err)
	}

	return aggregateBlocks(blocks), nil
}

// mergeBlock дедуплікує блок покриття: atomic — сума, set/count — максимум.
func mergeBlock(existing *coverageBlock, count int, mode string) {
	if mode == "atomic" {
		existing.count += count
	} else if count > existing.count {
		existing.count = count
	}
}

// aggregateBlocks агрегує блоки покриття по файлах.
func aggregateBlocks(blocks map[string]*coverageBlock) map[string]*FileStats {
	fileStats := make(map[string]*FileStats)

	for _, block := range blocks {
		stats, exists := fileStats[block.filePath]
		if !exists {
			stats = &FileStats{}
			fileStats[block.filePath] = stats
		}

		stats.Statements += block.statements
		if block.count > 0 {
			stats.Covered += block.statements
		}
	}

	return fileStats
}

// DetectModulePrefix автоматично визначає префікс модуля з файлу coverage.out.
// Шукає відомі кореневі директорії (cmd, internal, pkg тощо) та будує prefix зі шляху до них.
func DetectModulePrefix(inputPath string) (string, error) {
	file, err := os.Open(inputPath) //nolint:gosec // G304: CLI tool reads user-provided files
	if err != nil {
		return "", fmt.Errorf("Не вдалося відкрити файл %s: %w", inputPath, err)
	}

	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "mode:") {
			continue
		}

		// Знаходимо останній ':' перед координатами блоку
		colonIndex := strings.LastIndex(line, ":")
		if colonIndex == -1 {
			continue
		}

		filePath := line[:colonIndex]
		parts := strings.Split(filePath, "/")

		for index, part := range parts {
			if knownRootDirectories[part] {
				prefix := strings.Join(parts[:index], "/")
				if prefix != "" {
					return prefix + "/", nil
				}

				return "", nil
			}
		}
	}

	return "", nil
}
