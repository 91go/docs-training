# 通用Makefile
.PHONY: all serve build

# xxx
CLI = xxx
BINARY = xxx
# 测试
GoTest = @go test
CoverageFile = coverage.out

all: serve

test:
	@go test -v ./...
	@go test -race ./...

# 通过benchstat对比测试数据
benchstat:
	@go test -run=NONE -benchmem -bench=Rss -count=20 | tee -a xxx.txt
	benchstat old.txt new.txt

# 查看测试覆盖率
cover:
	@go test `go list ./... | grep -v examples` -coverprofile=coverage.out -covermode=atomic
	@go tool cover -html=coverage.out

build:
	@go build -o $(CLI) main.go

# 交叉编译
cross-build:
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 @go build -a -o ./build/$(BINARY) -v ./

clean:
	@if [ -f $(BINARY) ] ; then rm $(BINARY) ; fi

mod/why:
	@go mod why $(package)

## Run data race detector
race:
	@go test -race -short

## Generate global code coverage report
coverage:
	sh ./coverage.sh;

## Generate global code coverage report in HTML
coverhtml:
	sh ./coverage.sh html;

help:
	@echo "make tidy 执行linter"
	@echo "make cover 查看测试覆盖率"
	@echo "make cross-build 交叉编译"
	@echo "make clean 移除二进制文件"
