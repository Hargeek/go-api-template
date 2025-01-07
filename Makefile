ALL_SRC := $(shell find . -name "*_test.go" | grep -v -e vendor \
	-e ".*/\..*" \
	-e ".*/_.*" \
	-e ".*/mocks.*")
TEST_DIRS := $(sort $(dir $(filter %_test.go,$(ALL_SRC))))
COVERAGE_FILE := coverage.out

SERVICE_NAME ?= go-api-template
BRANCH 		 ?= $(shell git name-rev --name-only HEAD|cut -d '/' -f 3-)
REVISION     ?= $(shell git rev-parse HEAD)
BUILD_DATE   ?= $(shell date -I'seconds')
BUILD_USER   ?= $(shell whoami)@$(shell hostname)
TAG_VERSION  ?= $(shell git describe --tags --abbrev=0)

VERSION_LDFLAGS := \
	-X ${SERVICE_NAME}/common/types.Branch=$(BRANCH) \
	-X ${SERVICE_NAME}/common/types.Revision=$(REVISION) \
	-X ${SERVICE_NAME}/common/types.BuildDate=$(BUILD_DATE) \
	-X ${SERVICE_NAME}/common/types.BuildUser=$(BUILD_USER)

default: help
help:
	@# print help first, so it's visible
	@printf "\033[36m%-20s\033[0m %s\n" 'help' '打印帮助信息'
	@# then everything matching "target: ## magic comments"
	@# @awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
	@cat $(MAKEFILE_LIST) | grep -e "^[a-zA-Z_\-]*:.* ## .*" | awk 'BEGIN {FS = ":.*? ## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

install-hook: ## 安装一些git hook
	cp -r .githooks/* .git/hooks/
	chmod +x .git/hooks/*

swagger: ## 生成swagger接口文档
	sh -c 'swag fmt && swag init  --parseDependency --parseInternal --parseDepth 1'

fieldalignment: ## 内存对齐检查
	fieldalignment -fix  ./...

vet: ## go vet检查
	go vet ./...

generate-error: ## 生成错误码
	sh -c 'cd common/error && go generate .'

run: ## 运行默认环境
	go run main.go

run-air: ## 以air运行
	air -c .air.toml -d

build: clean ## 编译二进制文件
	go build -v -o bin/app main.go

buildx: ## 用于cicd pipeline中docker编译二进制文件
	go build -ldflags="-s -w $(VERSION_LDFLAGS)" -o /bin/server main.go

test: ## 运行测试用例及测试编译
	echo $(TEST_DIRS)
	@rm -f $(COVERAGE_FILE)
	set -o pipefail;\
	for dir in $(TEST_DIRS); do \
		go test -v -timeout 20m -coverprofile="test.temp" "$$dir" | tee -a $(COVERAGE_FILE) || exit 1; \
	done;
	go build -o /dev/null

local: ## 运行本地环境
	rm config/conf.yaml || true
	cp config/conf-local.yaml config/conf.yaml
	go run main.go

.PHONY: clean
clean:
	go clean
