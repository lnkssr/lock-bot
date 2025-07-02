BINARY = cmd/lock-bot/main.go

init:
	rm -f lock-bot
	go mod download
	go build -o lock-bot ${BINARY}

dev: 
	go mod download
	go run ${BINARY}

delete: 
	rm -f lock-bot