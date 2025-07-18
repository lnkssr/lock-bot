BINARY = cmd/lock-bot/main.go
NAME = build/lock-bot

init: 
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

docker: docker_del
	@docker compose up --build -d

docker_del:
	@docker compose down --rmi all --volumes --remove-orphans

logs:
	@docker compose logs -f app

.PHONY: logs docker docker_del dev init delete 
