run:
	go run *.go --config=config.ini

install:
	go get gopkg.in/mgo.v2;
	go get gopkg.in/mgo.v2/bson;
	go get gopkg.in/ini.v1
	go get github.com/gorilla/websocket
	go get github.com/stretchr/testify/assert

build:
	go build -o $(GOPATH)/bin/osx-omws *.go;
	GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -o $(GOPATH)/bin/omws *.go;

release:
	make build;
	echo "Release"
