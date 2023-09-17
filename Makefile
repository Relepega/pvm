VERSION = 0.1.0-b1

.PHONY: build

# go tool dist list | grep windows
build:
	rm -rf build
	mkdir build
	GOOS=windows GOARCH=amd64 go build -o ./build ./cmd/main/pvm.go
	cd build && tar -a -c -f pvm-v$(VERSION).zip pvm.exe
