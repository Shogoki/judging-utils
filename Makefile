BINARY_NAME=judging-tools

build: 
	go build -o bin/${BINARY_NAME} main.go

build_all: 
	GOARCH=amd64 GOOS=darwin go build -o bin/${BINARY_NAME}-macos-amd64 main.go
	GOARCH=amd64 GOOS=linux go build -o bin/${BINARY_NAME}-linux-amd64 main.go
	GOARCH=amd64 GOOS=windows go build -o bin/${BINARY_NAME}-windows-amd64.exe main.go
	GOARCH=arm64 GOOS=darwin go build -o bin/${BINARY_NAME}-macos-arm64 main.go
	GOARCH=arm64 GOOS=linux go build -o bin/${BINARY_NAME}-linux-arm64 main.go
test:
	echo "No Tests yet"
run: build
	bin/${BINARY_NAME}
clean:
	go clean
	rm bin/${BINARY_NAME}  2> /dev/null
	rm bin/${BINARY_NAME}-*  2> /dev/null
