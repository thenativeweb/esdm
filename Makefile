OS_NAME := $(shell uname)
ifeq ($(OS_NAME), Darwin)
OPEN := open
else
OPEN := xdg-open
endif

WITH_RACE_DETECTION ?= false

qa: analyze lint-documentation-examples test


analyze:
	@go vet ./...
	@go tool staticcheck --checks=all ./...


lint-documentation-examples:
	@if [ -d documentation/guides ] && find documentation/guides -name '*.esdm.yaml' -print -quit | grep -q .; then \
		go run ./cmd/esdm lint --directory documentation/guides; \
	fi


dev-documentation: lint-documentation-examples generate-version-json
	@docker run --rm -it -p 3000:3000 -v ./documentation:/docs \
		squidfunk/mkdocs-material:9.7.6 serve -a 0.0.0.0:3000


build-documentation: lint-documentation-examples generate-version-json
	@docker build -t esdm-documentation:local ./documentation


generate-version-json:
	$(eval VERSION=$(shell git describe --tags --abbrev=0 2>/dev/null | sed 's/^v//'))
	$(eval VERSION=$(or $(VERSION),0.0.0))
	@printf '{"version": "%s"}\n' "$(VERSION)" > documentation/docs/version.json


build-schema-host:
	@docker build -t esdm-schema:local ./schema


generate-reference-snippets:
	@go run ./cmd/refgen


test:
ifeq ($(WITH_RACE_DETECTION), false)
	@go test -failfast -cover ./...
else
	@go test -race -count=1 -failfast -cover ./...
endif


coverage: test
	@mkdir -p ./coverage
	@go test -failfast -coverprofile=./coverage/cover.out ./...
	@go tool cover -html=./coverage/cover.out -o ./coverage/cover.html
	@$(OPEN) ./coverage/cover.html


clean:
	@echo "Cleaning the build directory..."
	@rm -rf build/


build: clean
	$(eval VERSION=$(shell git tag --points-at HEAD))
	$(eval VERSION=$(or $(VERSION), (version unavailable)))
	$(eval GIT_VERSION=$(shell git rev-parse HEAD))

	@echo "Building esdm-darwin-arm64..."
	@GOOS=darwin GOARCH=arm64 go build \
		-trimpath \
		-ldflags="-buildid= -s -w -X 'github.com/thenativeweb/esdm/cmd/cmdutils.Version=$(VERSION)' -X 'github.com/thenativeweb/esdm/cmd/cmdutils.GitVersion=$(GIT_VERSION)'" \
		-o ./build/esdm-darwin-arm64 \
		./cmd/esdm

	@echo "Building esdm-darwin-amd64..."
	@GOOS=darwin GOARCH=amd64 go build \
		-trimpath \
		-ldflags="-buildid= -s -w -X 'github.com/thenativeweb/esdm/cmd/cmdutils.Version=$(VERSION)' -X 'github.com/thenativeweb/esdm/cmd/cmdutils.GitVersion=$(GIT_VERSION)'" \
		-o ./build/esdm-darwin-amd64 \
		./cmd/esdm

	@echo "Building esdm-linux-arm64..."
	@GOOS=linux GOARCH=arm64 go build \
		-trimpath \
		-ldflags="-buildid= -s -w -X 'github.com/thenativeweb/esdm/cmd/cmdutils.Version=$(VERSION)' -X 'github.com/thenativeweb/esdm/cmd/cmdutils.GitVersion=$(GIT_VERSION)'" \
		-o ./build/esdm-linux-arm64 \
		./cmd/esdm

	@echo "Building esdm-linux-amd64..."
	@GOOS=linux GOARCH=amd64 go build \
		-trimpath \
		-ldflags="-buildid= -s -w -X 'github.com/thenativeweb/esdm/cmd/cmdutils.Version=$(VERSION)' -X 'github.com/thenativeweb/esdm/cmd/cmdutils.GitVersion=$(GIT_VERSION)'" \
		-o ./build/esdm-linux-amd64 \
		./cmd/esdm

	@echo "Building esdm-windows-arm64.exe..."
	@GOOS=windows GOARCH=arm64 go build \
		-trimpath \
		-ldflags="-buildid= -s -w -X 'github.com/thenativeweb/esdm/cmd/cmdutils.Version=$(VERSION)' -X 'github.com/thenativeweb/esdm/cmd/cmdutils.GitVersion=$(GIT_VERSION)'" \
		-o ./build/esdm-windows-arm64.exe \
		./cmd/esdm

	@echo "Building esdm-windows-amd64.exe..."
	@GOOS=windows GOARCH=amd64 go build \
		-trimpath \
		-ldflags="-buildid= -s -w -X 'github.com/thenativeweb/esdm/cmd/cmdutils.Version=$(VERSION)' -X 'github.com/thenativeweb/esdm/cmd/cmdutils.GitVersion=$(GIT_VERSION)'" \
		-o ./build/esdm-windows-amd64.exe \
		./cmd/esdm


.PHONY: analyze \
				build \
				build-documentation \
				build-schema-host \
				clean \
				coverage \
				dev-documentation \
				generate-reference-snippets \
				generate-version-json \
				lint-documentation-examples \
				qa \
				test
