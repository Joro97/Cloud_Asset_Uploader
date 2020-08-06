IMAGE_NAME := cloud-asset-uploader

format:
	gofmt -w .

linter:
	golint ./...

image:
	docker build . -t ${IMAGE_NAME}
