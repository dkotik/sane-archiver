default:
	# mkdir -p tests/data
	# go test
	# cd cmd/sane-archiver && go test
	# cd cmd/sane-archiver && go run . keygen
	cd cmd/sane-archiver && go install .
linux:
	echo "Building for Linux..."
	cd cmd/sane-archiver; \
	GOOS=linux GOARCH=amd64 go build -trimpath -o ../../build/sane-archiver
arm:
	echo "Building for Raspberry Pi..."
	cd cmd/sane-archiver; \
	GOOS=linux GOARCH=arm GOARM=5 go build -trimpath -o ../../build/sane-archiver-arm
macos:
	echo "Building for MacOS..."
	cd cmd/sane-archiver; \
	GOOS=darwin GOARCH=amd64 go build -trimpath -o ../../build/sane-archiver-macos
clean:
	rm -rf build
	rm -rf tests/data
