main_package_path = ./cmd/github-activity/
binary_name = github-activity
build_dir = ./bin

.PHONY: help
help:
	@echo "Usage:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ] 

## build: build the application
.PHONY: build
build:
	@mkdir -p $(build_dir)
	@if [[ "$$DEBUG" = "1" ]]; then \
		echo "Building with debug flags..."; \
		go build --gcflags="all=-N -l" -o $(build_dir)/$(binary_name) $(main_package_path); \
	else \
		echo "Building production binary..."; \
		go build -o $(build_dir)/$(binary_name) $(main_package_path); \
	fi

## run: run the application
.PHONY: run
run: build
# 	${build_dir}/${binary_name} ivansevryukov1995
	${build_dir}/${binary_name} snow-ghost
# 	${build_dir}/${binary_name} valyala
# 	${build_dir}/${binary_name} yangwenmai

## clear-cache: clear cache the application
.PHONY: clear-cache
clear-cache: build
	${build_dir}/${binary_name} --clear-cache

## test: testing the application
.PHONY: test
test:
	go test -v -count=1 ./...

## clean: clean dir bin
.PHONY: clean
clean: confirm
	@echo "Clean up..."
	@rm -rf ${build_dir}

## debug: debugging binary file with debug flags
.PHONY: debug
debug:
	DEBUG=1 $(MAKE) build
	dlv exec $(build_dir)/$(binary_name) -- ivansevryukov1995