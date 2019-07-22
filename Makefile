default:
	# go test -run TestS3
	# go test
	cd cmd/sane-archiver && go test
linux:
	echo "Building for Linux..."
	cd cmd/sane-archiver; \
	GOOS=linux GOARCH=amd64 go build -o ../../build/sane-archiver
arm:
	echo "Building for Raspberry Pi..."
	cd cmd/sane-archiver; \
	GOOS=linux GOARCH=arm GOARM=5 go build -o ../../build/sane-archiver-arm
macos:
	echo "Building for MacOS..."
	cd cmd/sane-archiver; \
	GOOS=darwin GOARCH=amd64 go build -o ../../build/sane-archiver-macos
install: linux
	cp build/sane-archiver ~/.local/bin/sane-archiver
clean:
	rm -r -f build
