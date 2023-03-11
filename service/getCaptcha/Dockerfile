FROM golang:1.18.1-alpine AS builder
#FROM golang:alpine AS builder

# Set Go env
ENV CGO_ENABLED=0 GOOS=linux GOPROXY=https://goproxy.cn,direct GO111MODULE=on
WORKDIR /go/src/getCaptcha

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

# Install dependencies
RUN apk --update --no-cache add ca-certificates gcc libtool make musl-dev protoc git

# Build Go binary
COPY Makefile go.mod go.sum ./
RUN make init && go mod download 
COPY . .
RUN make proto tidy build

# Deployment container
FROM scratch

COPY --from=builder /etc/ssl/certs /etc/ssl/certs
COPY --from=builder /go/src/getCaptcha/getCaptcha /getCaptcha
ENTRYPOINT ["/getCaptcha"]
CMD []
