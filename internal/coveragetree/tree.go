package coveragetree

import (
	"math"
	"sort"
	"strings"
)

// TreeNode представляє вузол дерева (директорію) з дочірніми вузлами та файлами.
type TreeNode struct {
	children   map[string]*TreeNode
	files      []*TreeFile
	statements int
	covered    int
}

// TreeFile представляє файл у дереві покриття.
type TreeFile struct {
	name       string
	statements int
	covered    int
}

// TreeJSON — JSON-представлення вузла дерева для вставки в HTML-шаблон.
type TreeJSON struct {
	Name       string      `json:"name"`
	Statements int         `json:"statements"`
	Covered    int         `json:"covered"`
	Coverage   float64     `json:"coverage"`
	Children   []*TreeJSON `json:"children"`
	Files      []*FileJSON `json:"files"`
}

// FileJSON — JSON-представлення файлу для вставки в HTML-шаблон.
type FileJSON struct {
	Name       string  `json:"name"`
	Statements int     `json:"statements"`
	Covered    int     `json:"covered"`
	Coverage   float64 `json:"coverage"`
}

// BuildTree створює дерево директорій з файлових шляхів та їхньої статистики покриття.
func BuildTree(fileStats map[string]*FileStats) *TreeNode {
	root := &TreeNode{
		children: make(map[string]*TreeNode),
	}

	for filePath, stats := range fileStats {
		parts := strings.Split(filePath, "/")
		currentNode := root

		// Створюємо директорії (всі частини шляху крім останньої)
		for _, directory := range parts[:len(parts)-1] {
			if _, exists := currentNode.children[directory]; !exists {
				currentNode.children[directory] = &TreeNode{
					children: make(map[string]*TreeNode),
				}
			}

			currentNode = currentNode.children[directory]
		}

		// Додаємо файл до листкового вузла
		fileName := parts[len(parts)-1]
		currentNode.files = append(currentNode.files, &TreeFile{
			name:       fileName,
			statements: stats.Statements,
			covered:    stats.Covered,
		})
	}

	return root
}

// Aggregate рекурсивно агрегує статистику покриття від листків до кореня.
func Aggregate(node *TreeNode) {
	totalStatements := 0
	totalCovered := 0

	// Агрегація файлів поточного вузла
	for _, file := range node.files {
		totalStatements += file.statements
		totalCovered += file.covered
	}

	// Рекурсивна агрегація дочірніх вузлів
	for _, child := range node.children {
		Aggregate(child)
		totalStatements += child.statements
		totalCovered += child.covered
	}

	node.statements = totalStatements
	node.covered = totalCovered
}

// ToJSON конвертує дерево у JSON-структуру для HTML-шаблону.
func ToJSON(node *TreeNode, name string) *TreeJSON {
	coverage := CalculateCoverage(node.statements, node.covered)

	result := &TreeJSON{
		Name:       name,
		Statements: node.statements,
		Covered:    node.covered,
		Coverage:   coverage,
		Children:   make([]*TreeJSON, 0),
		Files:      make([]*FileJSON, 0),
	}

	// Сортуємо дочірні вузли за назвою
	childNames := make([]string, 0, len(node.children))
	for childName := range node.children {
		childNames = append(childNames, childName)
	}

	sort.Strings(childNames)

	for _, childName := range childNames {
		result.Children = append(result.Children, ToJSON(node.children[childName], childName))
	}

	// Сортуємо файли за назвою
	sortedFiles := make([]*TreeFile, len(node.files))
	copy(sortedFiles, node.files)
	sort.Slice(sortedFiles, func(first, second int) bool {
		return sortedFiles[first].name < sortedFiles[second].name
	})

	for _, file := range sortedFiles {
		fileCoverage := CalculateCoverage(file.statements, file.covered)
		result.Files = append(result.Files, &FileJSON{
			Name:       file.name,
			Statements: file.statements,
			Covered:    file.covered,
			Coverage:   fileCoverage,
		})
	}

	return result
}

// coverageScale використовується для округлення покриття до 1 десяткового знаку.
const coverageScale = 1000

// coverageRoundFactor — дільник для отримання 1 десяткового знаку.
const coverageRoundFactor = 10

// CalculateCoverage обчислює відсоток покриття з округленням до 1 десяткового знаку.
func CalculateCoverage(statements int, covered int) float64 {
	if statements == 0 {
		return 0.0
	}

	return math.Round(float64(covered)/float64(statements)*coverageScale) / coverageRoundFactor
}
