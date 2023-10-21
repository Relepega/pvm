VERSION = 0.1.0.b2

.PHONY: build

# go tool dist list | grep windows
build:
	rm -rf build
	mkdir build

	GOOS=windows GOARCH=386 go build -o ./build ./cmd/main/pvm.go
	cd build && tar -a -c -f pvm-v$(VERSION)-x86.zip pvm.exe

	GOOS=windows GOARCH=amd64 go build -o ./build ./cmd/main/pvm.go
	cd build && tar -a -c -f pvm-v$(VERSION)-x64.zip pvm.exe
