default:
	# go test -run TestS3
	go test
linux:
	echo "Building for Linux..."
	export GOOS=linux; \
	export GOARCH=amd64; \
	cd cmd/sane-archiver; \
	go build -o ../../build/sane-archiver
arm:
	echo "Building for Raspberry Pi..."
	export GOOS=linux; \
	export GOARCH=arm; \
	export GOARM=5; \
	cd cmd/sane-archiver; \
	go build -o ../../build/sane-archiver-arm
macos:
	echo "Building for MacOS..."
	export GOOS=darwin; \
	export GOARCH=amd64; \
	cd cmd/sane-archiver; \
	go build -o ../../build/sane-archiver-macos
install: linux
	cp build/sane-archiver ~/.local/bin/sane-archiver
clean:
	rm -r -f build
