# 启动本地基础服务
local-infra-up:
	./deployments/local/infra/up.sh

# 关闭本地基础服务
local-infra-down:
	./deployments/local/infra/down.sh

# 运行单元测试
.PHONY: test
test: CONFIG ?= ""
test:
	./test/run.sh ${CONFIG}

# 执行 CI 构建
ci-build:
	./build/ci/build-code.sh
