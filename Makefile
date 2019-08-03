.PHONY: dev
dev: main.go
	go build
	auth-server.exe