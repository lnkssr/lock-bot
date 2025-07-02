BINARY = cmd/lock-bot/main.go

init:
	rm -f lock-bot
	go build -o lock-bot ${BINARY}

dev: 
	go run ${BINARY}

delete: 
	rm -f lock-bot