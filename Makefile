default:
	go test
	# cd cmd/sane-archiver && \
	# go test && \
	# go run .
linux:
	export GOOS=linux; \
	export GOARCH=amd64; \
	cd cmd/sane-archiver; \
	go build -o ../../build/sane-archiver
arm:
	export GOOS=linux; \
	export GOARCH=arm; \
	export GOARM=5; \
	cd cmd/sane-archiver; \
	go build -o ../../build/sane-archiver
install:
	sudo cp build/sane-archiver /usr/bin/sane-archiver
clean:
	rm -r -f build
