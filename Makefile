USER ?= dhf0820
HUB ?= dhf0820
TAG ?= latest
TEST = dhf0820
PROD = vertisoft
VERSION ?= $(TAG)
IMAGE_NAME ?= email_connector
BINARY_NAME=$(IMAGE_NAME)
BINARY_UNIX=$(BINARY_NAME)_linux
LINUX_IMAGE_NAME ?= email_linux
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_MAC=$(IMAGE_NAME)_mac
MAC_IMAGE_NAME= $(BINARY_MAC)

api:
	@protoc -I ./protobufs/ \
		--proto_path=./ \
		--go_out=plugins=grpc:./ \
		./protobufs/*.proto


dep: ## Get the dependencies
	@go get -v -d ./...

tidy: # add all new includes
	@go mod tidy

release:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v
	docker build -t $(HUB)/$(IMAGE_NAME):$(VERSION) -f Dockerfile .
	docker push $(NS)/$(IMAGE_NAME):$(VERSION)


build:
	$(GOBUILD) -o $(BINARY_NAME) -v

mac:
	CGO_ENABLED=0 $(GOBUILD) -o $(BINARY_MAC) -v

test:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v
	docker build -t $(TEST)/$(IMAGE_NAME):$(VERSION) -f Dockerfile .
	docker push $(TEST)/$(IMAGE_NAME):$(VERSION)
prod:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v
	docker build -t $(PROD)/$(IMAGE_NAME):$(VERSION) -f Dockerfile .
	docker push $(PROD)/$(IMAGE_NAME):$(VERSION)

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v

docker-build:
	docker build -t $(HUB)/$(IMAGE_NAME):$(VERSION) -f Dockerfile .

docker-push: # push to docker
	docker push $(NS)/$(IMAGE_NAME):$(VERSION)



build_client:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o ./src/client/client ./src/client/

build_server:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o ./src/server/client ./src/server/

run_server:
	go run main.go

run-client:
	go run cmd/client/client.go  -address localhost:9010 -o 4