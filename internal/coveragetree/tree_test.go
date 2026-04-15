package coveragetree_test

import (
	"testing"

	"github.com/romantelychko/coverage-tree/internal/coveragetree"
)

// findChild шукає дочірній вузол за назвою в TreeJSON.
func findChild(node *coveragetree.TreeJSON, name string) *coveragetree.TreeJSON {
	for _, child := range node.Children {
		if child.Name == name {
			return child
		}
	}

	return nil
}

// TestBuildTreeStructure перевіряє побудову дерева з файлових шляхів.
func TestBuildTreeStructure(test *testing.T) {
	fileStats := map[string]*coveragetree.FileStats{
		"internal/users/service.go":   {Statements: 10, Covered: 8},
		"internal/users/handler.go":   {Statements: 5, Covered: 3},
		"internal/auth/middleware.go": {Statements: 6, Covered: 6},
		"cmd/server/main.go":          {Statements: 3, Covered: 3},
	}

	tree := coveragetree.BuildTree(fileStats)
	coveragetree.Aggregate(tree)
	treeJSON := coveragetree.ToJSON(tree, "root")

	// Перевіряємо наявність кореневих директорій
	if findChild(treeJSON, "internal") == nil {
		test.Error("Директорія 'internal' не знайдена в корені дерева")
	}

	if findChild(treeJSON, "cmd") == nil {
		test.Error("Директорія 'cmd' не знайдена в корені дерева")
	}

	// Перевіряємо вкладеність
	internalNode := findChild(treeJSON, "internal")
	if internalNode == nil {
		test.Fatal("Директорія 'internal' не знайдена")
	}

	if findChild(internalNode, "users") == nil {
		test.Error("Директорія 'users' не знайдена в 'internal'")
	}

	if findChild(internalNode, "auth") == nil {
		test.Error("Директорія 'auth' не знайдена в 'internal'")
	}

	// Перевіряємо файли у вузлі
	usersNode := findChild(internalNode, "users")
	if usersNode == nil {
		test.Fatal("Директорія 'users' не знайдена")
	}

	if len(usersNode.Files) != 2 {
		test.Errorf("Очікувалось 2 файли в 'users', отримано %d", len(usersNode.Files))
	}
}

// TestAggregateTree перевіряє рекурсивну агрегацію статистики.
func TestAggregateTree(test *testing.T) {
	fileStats := map[string]*coveragetree.FileStats{
		"internal/users/service.go":   {Statements: 10, Covered: 8},
		"internal/users/handler.go":   {Statements: 5, Covered: 3},
		"internal/auth/middleware.go": {Statements: 6, Covered: 6},
		"cmd/server/main.go":          {Statements: 3, Covered: 3},
	}

	tree := coveragetree.BuildTree(fileStats)
	coveragetree.Aggregate(tree)
	treeJSON := coveragetree.ToJSON(tree, "root")

	// Загальна статистика кореня
	expectedStatements := 24
	expectedCovered := 20

	if treeJSON.Statements != expectedStatements {
		test.Errorf("Корінь: очікувалось %d statements, отримано %d", expectedStatements, treeJSON.Statements)
	}

	if treeJSON.Covered != expectedCovered {
		test.Errorf("Корінь: очікувалось %d covered, отримано %d", expectedCovered, treeJSON.Covered)
	}

	// Статистика internal
	internalNode := findChild(treeJSON, "internal")
	if internalNode == nil {
		test.Fatal("Директорія 'internal' не знайдена")
	}

	expectedInternalStatements := 21
	expectedInternalCovered := 17

	if internalNode.Statements != expectedInternalStatements {
		test.Errorf(
			"internal: очікувалось %d statements, отримано %d",
			expectedInternalStatements, internalNode.Statements,
		)
	}

	if internalNode.Covered != expectedInternalCovered {
		test.Errorf(
			"internal: очікувалось %d covered, отримано %d",
			expectedInternalCovered, internalNode.Covered,
		)
	}
}

// TestToJSON перевіряє конвертацію дерева в JSON-структуру.
func TestToJSON(test *testing.T) {
	fileStats := map[string]*coveragetree.FileStats{
		"internal/users/service.go": {Statements: 10, Covered: 7},
		"cmd/main.go":               {Statements: 5, Covered: 5},
	}

	tree := coveragetree.BuildTree(fileStats)
	coveragetree.Aggregate(tree)
	treeJSON := coveragetree.ToJSON(tree, "root")

	if treeJSON.Name != "root" {
		test.Errorf("Очікувалось ім'я 'root', отримано %q", treeJSON.Name)
	}

	if treeJSON.Statements != 15 {
		test.Errorf("Очікувалось 15 statements, отримано %d", treeJSON.Statements)
	}

	if treeJSON.Covered != 12 {
		test.Errorf("Очікувалось 12 covered, отримано %d", treeJSON.Covered)
	}

	// 12/15 = 80.0%
	if treeJSON.Coverage != 80.0 {
		test.Errorf("Очікувалось 80.0%% coverage, отримано %.1f%%", treeJSON.Coverage)
	}

	// Перевіряємо сортування дочірніх вузлів
	if len(treeJSON.Children) != 2 {
		test.Fatalf("Очікувалось 2 дочірніх вузли, отримано %d", len(treeJSON.Children))
	}

	if treeJSON.Children[0].Name != "cmd" {
		test.Errorf("Перший дочірній вузол має бути 'cmd', отримано %q", treeJSON.Children[0].Name)
	}

	if treeJSON.Children[1].Name != "internal" {
		test.Errorf("Другий дочірній вузол має бути 'internal', отримано %q", treeJSON.Children[1].Name)
	}
}

// TestToJSONWithZeroStatements перевіряє коректну обробку файлу з 0 statements.
func TestToJSONWithZeroStatements(test *testing.T) {
	fileStats := map[string]*coveragetree.FileStats{
		"empty.go": {Statements: 0, Covered: 0},
	}

	tree := coveragetree.BuildTree(fileStats)
	coveragetree.Aggregate(tree)
	treeJSON := coveragetree.ToJSON(tree, "root")

	if treeJSON.Coverage != 0.0 {
		test.Errorf("Очікувалось 0.0%% coverage для порожнього файлу, отримано %.1f%%", treeJSON.Coverage)
	}
}

// TestCalculateCoverage перевіряє обчислення відсотка покриття.
func TestCalculateCoverage(test *testing.T) {
	testCases := []struct {
		name       string
		statements int
		covered    int
		expected   float64
	}{
		{"повне покриття", 10, 10, 100.0},
		{"нульове покриття", 10, 0, 0.0},
		{"часткове покриття", 3, 1, 33.3},
		{"нуль statements", 0, 0, 0.0},
		{"дві третини", 3, 2, 66.7},
	}

	for _, testCase := range testCases {
		test.Run(testCase.name, func(test *testing.T) {
			result := coveragetree.CalculateCoverage(testCase.statements, testCase.covered)
			if result != testCase.expected {
				test.Errorf(
					"CalculateCoverage(%d, %d) = %.1f, очікувалось %.1f",
					testCase.statements, testCase.covered, result, testCase.expected,
				)
			}
		})
	}
}
