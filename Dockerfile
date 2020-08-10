FROM golang:1.14 AS builder

LABEL maintainer="Georgi Karov"

WORKDIR /app

COPY . .

RUN go mod download

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

######## Start a new stage from scratch #######
FROM alpine:3.9

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy pre-build image from previous stage
COPY --from=builder /app/main .

CMD ["./main"]
