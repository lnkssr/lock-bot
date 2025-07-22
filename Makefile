BINARY = cmd/lock-bot/main.go
NAME = build/lock-bot

install: 
	@rm -rf build 
	@mkdir build
	@go mod download
	@go build -o ${NAME} ${BINARY}

check:
	@./scripts/check.sh

dev: check
	@go mod download
	@go run ${BINARY}

clear: 
	@rm -rf build

init: down
	@docker compose up --build -d

down:
	@docker compose down --rmi all --volumes --remove-orphans

logs:
	@docker compose logs -f app

.PHONY: logs init down dev init delete 
