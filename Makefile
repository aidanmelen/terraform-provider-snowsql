TEST?=$$(go list ./... | grep -v 'vendor')
REGISTRY=registry.terraform.io
HOSTNAME=aidanmelen
NAME=snowsql
BINARY=terraform-provider-${NAME}
VERSION=1.2.0
OS_ARCH=darwin_amd64

help: ## This help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-35s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help

build: ## build the provider
	go build -o ${BINARY}

clean: ## clean the repo
	rm terraform-provider-snowflake 2>/dev/null || true
	go clean
	rm -rf dist

doc: ## generate the provider documentation
	tfplugindocs generate

install: build ## build and install the provider
	mkdir -p ~/.terraform.d/plugins/${REGISTRY}/${HOSTNAME}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${REGISTRY}/${HOSTNAME}/${NAME}/${VERSION}/${OS_ARCH}

pre-commit: ## run pre-commit checks
	git init
	git add -A
	pre-commit run -a
	git add -A

test: ## run the provider unittests
	go test -i $(TEST) || exit 1
	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

testacc: ## run the provider acceptance tests
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

tools: ## Get the development tools
	go get -u golang.org/x/lint/golint
	go get golang.org/x/tools/cmd/goimports
	go get github.com/bflad/tfproviderdocs
	go get github.com/bflad/tfproviderlint/cmd/tfproviderlint
	go get github.com/katbyte/terrafmt
	go get github.com/hashicorp/terraform-plugin-docs

local-release: ## Release provider locally
	GOOS=darwin GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_darwin_amd64
	GOOS=freebsd GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_freebsd_386
	GOOS=freebsd GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_freebsd_amd64
	GOOS=freebsd GOARCH=arm go build -o ./bin/${BINARY}_${VERSION}_freebsd_arm
	GOOS=linux GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_linux_386
	GOOS=linux GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_linux_amd64
	GOOS=linux GOARCH=arm go build -o ./bin/${BINARY}_${VERSION}_linux_arm
	GOOS=openbsd GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_openbsd_386
	GOOS=openbsd GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_openbsd_amd64
	GOOS=solaris GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_solaris_amd64
	GOOS=windows GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_windows_386
	GOOS=windows GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_windows_amd64

release: ## Trigger a release
	git checkout main
	git tag v$(VERSION)
	git push --tag