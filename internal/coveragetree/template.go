package coveragetree

import _ "embed"

//go:embed template/coverage_tree.html
var embeddedTemplate string

// GetTemplate повертає вбудований HTML-шаблон.
func GetTemplate() string {
	return embeddedTemplate
}
