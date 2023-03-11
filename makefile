GOPATH:=$(shell go env GOPATH)

.PHONY: consul
consul:
	@docker run -d --network=host --name=consul consul:latest agent -server -bootstrap -ui -node=1 -client='0.0.0.0' -bind='172.29.251.146'

.PHONY: getChatcha
chatcha:
	docker run --net=host -v ~/conf/:/conf --name getCaptcha registry.cn-heyuan.aliyuncs.com/whyy1/service:v1

.PHONY: user
user:
	@cd .\service\user\
	@go run main.go