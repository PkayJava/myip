build:
	go mod tidy

	GOOS=linux GOARCH=amd64 go build -o bin/linux/amd64/myip cmd/main.go
	GOOS=linux GOARCH=arm64 go build -o bin/linux/arm64/myip cmd/main.go

	GOOS=darwin GOARCH=amd64 go build -o bin/darwin/amd64/myip cmd/main.go
	GOOS=darwin GOARCH=arm64 go build -o bin/darwin/arm64/myip cmd/main.go

	GOOS=windows GOARCH=amd64 go build -o bin/windows/amd64/myip.exe cmd/main.go
	GOOS=windows GOARCH=arm64 go build -o bin/windows/arm64/myip.exe cmd/main.go
