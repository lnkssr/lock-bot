BINARY = cmd/lock-bot/main.go
NAME = build/lock-bot

init:
	rm -f ${NAME}
	mkdir build
	go mod download
	go build -o ${NAME} ${BINARY}

dev: 
	go mod download
	go run ${BINARY}

delete: 
	rm -rf build

docker: docker_del
	@docker compose up --build -d

docker_del:
	@docker compose down --rmi all --volumes --remove-orphans

logs:
	@docker compose logs -f app

.PHONY: logs docker docker_del dev init delete 