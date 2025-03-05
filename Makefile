.PHONY: all run-orchestrator run-agent clean

# Цель по умолчанию — запуск всего проекта
all: run

# Запуск оркестратора и агента в фоновом режиме
run:
	@echo "Starting orchestrator and agent..."
	@go run ./cmd/orchestrator/main.go & go run ./cmd/agent/main.go &

# Запуск только оркестратора (для отладки)
run-orchestrator:
	@echo "Starting orchestrator..."
	@go run ./cmd/orchestrator/main.go

# Запуск только агента (для отладки)
run-agent:
	@echo "Starting agent..."
	@go run ./cmd/agent/main.go

# Очистка (остановка процессов, если нужно)
clean:
	@echo "Stopping all processes..."
	@-pkill -f "go run ./cmd/orchestrator/main.go" || true
	@-pkill -f "go run ./cmd/agent/main.go" || true