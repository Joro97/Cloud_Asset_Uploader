IMAGE_NAME := cloud-asset-uploader
IMAGE_TAG := latest

image:
	docker build . -t ${IMAGE_NAME}:${IMAGE_TAG}

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

astra:
	docker run -v ~/.aws/:/root/.aws:ro -v ACTUAL_ABSOLUTE_PATH_TO_CREDS_HERE:/root/secure-connect --env-file ACTUAL_PATH_TO_ENV_FILE_HERE -p 8090:8090 ${IMAGE_NAME}:${IMAGE_TAG}
