BINARY_NAME = Connect4
build:
	GOOS=windows go build -o bin/${BINARY_NAME}-windows.exe
	GOOS=linux go build -o bin/${BINARY_NAME}-linux
