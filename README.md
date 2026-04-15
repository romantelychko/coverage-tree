# CoverageTree

Генерація HTML-звіту покриття Go-коду у вигляді інтерактивного дерева директорій.

Парсить `coverage.out` (стандартний вивід `go test -coverprofile`) та генерує HTML з collapsible tree, кольоровими індикаторами покриття та підтримкою тем і локалізації.

[![CI](https://github.com/romantelychko/coverage-tree/actions/workflows/ci.yml/badge.svg)](https://github.com/romantelychko/coverage-tree/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/romantelychko/coverage-tree)](https://goreportcard.com/report/github.com/romantelychko/coverage-tree)
[![Go Reference](https://pkg.go.dev/badge/github.com/romantelychko/coverage-tree.svg)](https://pkg.go.dev/github.com/romantelychko/coverage-tree)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Встановлення

```bash
go install github.com/romantelychko/coverage-tree/cmd/coverage-tree@latest
```

## Використання

```bash
# Базовий виклик
coverage-tree coverage.out coverage-tree.html

# З явним префіксом модуля
coverage-tree --module-prefix github.com/myorg/myproject/ coverage.out report.html

# Англійська мова, світла тема
coverage-tree --lang en --theme light coverage.out report.html

# Додаткові виключення
coverage-tree --exclude-dir vendor --exclude-dir generated coverage.out report.html

# Кастомний заголовок
coverage-tree --title "API Coverage Report" coverage.out report.html
```

## Прапорці

| Прапорець          | Опис                                                      | За замовчуванням       |
|--------------------|-----------------------------------------------------------|------------------------|
| `--module-prefix`  | Префікс модуля для обрізки шляхів                         | автодетект             |
| `--exclude-suffix` | Суфікс файлів для виключення (можна вказати кілька разів) | `_mock.go`, `_test.go` |
| `--exclude-dir`    | Директорія для виключення (можна вказати кілька разів)    | `mocks`                |
| `--theme`          | Колірна тема: `dark` або `light`                          | `dark`                 |
| `--lang`           | Мова інтерфейсу: `uk` або `en`                            | `uk`                   |
| `--title`          | Заголовок HTML-звіту                                      | визначається мовою     |
| `--no-autodetect`  | Вимкнути автодетект префікса модуля                       | `false`                |
| `--version`        | Вивести версію та вийти                                   |                        |

## Інтеграція з Makefile

```makefile
coverage:
	go test -coverprofile=coverage.out -coverpkg=./... ./...
	coverage-tree coverage.out coverage-tree.html
```

## Ліцензія

[MIT](LICENSE)
