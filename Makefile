IMAGE_NAME := cloud-asset-uploader

image:
	docker build . -t ${IMAGE_NAME}

format:
	gofmt -w .

linter:
	golint ./...

integration:
	go test -v -coverprofile cover.out ./... && \
    go tool cover -html=cover.out -o cover.html && \
    open cover.html

unit:
	go test -v ./server && \
	go test -v ./responses
