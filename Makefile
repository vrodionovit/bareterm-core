# Makefile

# Go команда
GO := go

# Имя исполняемого файла
BINARY_NAME := myapp

# Папка для сборки
BUILD_DIR := build

# Основной пакет
MAIN_PACKAGE := .

# Флаги для сборки
BUILD_FLAGS := -v

# Операционные системы и архитектуры для сборки
OS_ARCH := darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 windows/amd64

# Цель по умолчанию
.DEFAULT_GOAL := all

# Создание папки build
$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

# Сборка для всех платформ
.PHONY: all
all: $(BUILD_DIR) $(addprefix build-,$(OS_ARCH))

# Шаблон для сборки под конкретную ОС и архитектуру
define GOBUILD
.PHONY: build-$(1)/$(2)
build-$(1)/$(2): $(BUILD_DIR)
	GOOS=$(1) GOARCH=$(2) $(GO) build $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)_$(1)_$(2)$(if $(filter windows,$(1)),.exe) $(MAIN_PACKAGE)
endef

# Генерация целей для каждой ОС и архитектуры
$(foreach OSARCH,$(OS_ARCH),$(eval $(call GOBUILD,$(firstword $(subst /, ,$(OSARCH))),$(lastword $(subst /, ,$(OSARCH))))))

# Очистка
.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)

# Запуск тестов
.PHONY: test
test:
	$(GO) test -v ./...

# Проверка кода
.PHONY: lint
lint:
	golangci-lint run