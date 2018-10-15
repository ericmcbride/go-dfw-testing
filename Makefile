REPO := github.com/ericmcbride/go-dfw-testing

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o service

build-native:
	go build -o service

test:
	docker-compose build 
	docker-compose up -d 
	docker run --rm -it \
		-v `pwd`:/go/src/${REPO} \
		-w /go/src/${REPO} \
		--network host \
		golang:1.10.0-alpine3.7 \
		go test -v ./...
	docker-compose down

run:
	make build-linux
	docker-compose up

down:
	docker-compose down
