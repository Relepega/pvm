VERSION = 0.1.0

# go tool dist list | grep windows
build:
	rm -rf build
	mkdir build
	cp -r ./scripts ./build/scripts
	GOOS=windows GOARCH=amd64 go build -o ./build ./cmd/main/pvm.go
	cd build && tar -a -c -f pvm-$(VERSION).zip Scripts pvm.exe

	