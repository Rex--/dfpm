

all: linux_amd64 darwin_arm64 windows_amd64

.PHONY: linux_amd64
linux_amd64:
	GOOS=linux GOARCH=amd64 go build -o dist/dfpm-linux-amd64 .

.PHONY: darwin_arm64
darwin_arm64:
	GOOS=darwin GOARCH=arm64 go build -o dist/dfpm-darwin-arm64 .

.PHONY: windows_amd64
windows_amd64:
	GOOS=windows GOARCH=amd64 go build -o dist/dfpm-windows-amd64.exe .

.PHONY: dev
dev:
	go build .

.PHONY: clean
clean:
	rm -f ./dfpm ./dfpm.exe
	rm -rf ./dist/