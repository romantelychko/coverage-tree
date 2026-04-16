VERSION ?= dev

.PHONY: help
help: ## Показує це повідомлення
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-40s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## Будує бінарний файл coverage-tree
	go build -ldflags "-s -w -X main.version=$(VERSION)" -trimpath -o ./bin/coverage-tree ./cmd/coverage-tree

.PHONY: lint
lint: ## Виправляє помилки статичного аналізу коду
	go vet ./...
	golangci-lint run ./... --fix

.PHONY: test
test: ## Виконує тести
	go test ./... -v -race -timeout=5m -cover -coverprofile=./coverage.out -covermode=atomic -coverpkg=./...

.PHONY: install
install: ## Встановлює бінарний файл coverage-tree
	go install -ldflags "-s -w -X main.version=$(VERSION)" -trimpath ./cmd/coverage-tree

.PHONY: update
update: ## Оновлення залежностей Go
	# Оновлення залежностей
	go get -u ./...
	go mod tidy
	go mod download

.PHONY: clean
clean: ## Видаляє бінарний файл coverage-tree
	rm -f ./bin/coverage-tree
