#!/usr/bin/make -f

PACKAGES_SIMTEST=$(shell go list ./... | grep '/simulation')
VERSION := $(shell git describe --abbrev=6 --dirty --always --tags)
COMMIT := $(shell git log -1 --format='%H')
LEDGER_ENABLED ?= true
SDK_PACK := $(shell go list -m github.com/cosmos/cosmos-sdk | sed  's/ /\@/g')
BUILDDIR ?= $(CURDIR)/build
TEST_DOCKER_REPO=jackzampolin/sommtest
HTTPS_GIT := https://github.com/peggyjv/sommelier.git
DOCKER := $(shell which docker)
DOCKER_BUF := $(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace bufbuild/buf

export GO111MODULE = on

# process build tags

build_tags = netgo
ifeq ($(LEDGER_ENABLED),true)
  ifeq ($(OS),Windows_NT)
    GCCEXE = $(shell where gcc.exe 2> NUL)
    ifeq ($(GCCEXE),)
      $(error gcc.exe not installed for ledger support, please install or set LEDGER_ENABLED=false)
    else
      build_tags += ledger
    endif
  else
    UNAME_S = $(shell uname -s)
    ifeq ($(UNAME_S),OpenBSD)
      $(warning OpenBSD detected, disabling ledger support (https://github.com/cosmos/cosmos-sdk/issues/1988))
    else
      GCC = $(shell command -v gcc 2> /dev/null)
      ifeq ($(GCC),)
        $(error gcc not installed for ledger support, please install or set LEDGER_ENABLED=false)
      else
        build_tags += ledger
      endif
    endif
  endif
endif

ifeq (cleveldb,$(findstring cleveldb,$(SOM_BUILD_OPTIONS)))
  build_tags += gcc
endif
build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))

whitespace :=
whitespace += $(whitespace)
comma := ,
build_tags_comma_sep := $(subst $(whitespace),$(comma),$(build_tags))

# process linker flags

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=sommelier \
		  -X github.com/cosmos/cosmos-sdk/version.AppName=sommelier \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
		  -X "github.com/cosmos/cosmos-sdk/version.BuildTags=$(build_tags_comma_sep)"

ifeq (cleveldb,$(findstring cleveldb,$(SOM_BUILD_OPTIONS)))
  ldflags += -X github.com/cosmos/cosmos-sdk/types.DBBackend=cleveldb
endif
ifeq (,$(findstring nostrip,$(SOM_BUILD_OPTIONS)))
  ldflags += -w -s
endif
ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

BUILD_FLAGS := -tags "$(build_tags)" -ldflags '$(ldflags)'
# check for nostrip option
ifeq (,$(findstring nostrip,$(SOM_BUILD_OPTIONS)))
  BUILD_FLAGS += -trimpath
endif

# The below include contains the tools target.
include contrib/devtools/Makefile

###############################################################################
###                              Documentation                              ###
###############################################################################

all: install lint test

build: go.sum
	mkdir -p $(BUILDDIR)
	go build -mod=readonly $(BUILD_FLAGS) -o $(BUILDDIR)/ ./cmd/sommelier

build-linux: go.sum
	LEDGER_ENABLED=false GOOS=linux GOARCH=amd64 $(MAKE) build

build-contract-tests-hooks:
	mkdir -p $(BUILDDIR)
	go build -mod=readonly $(BUILD_FLAGS) -o $(BUILDDIR)/ ./cmd/contract_tests

install: go.sum
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/sommelier

go-mod-cache: go.sum
	@echo "--> Download go modules to local cache"
	@go mod download

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	@go mod verify

draw-deps:
	@# requires brew install graphviz or apt-get install graphviz
	go get github.com/RobotsAndPencils/goviz
	@goviz -i ./cmd/sommelier -d 2 | dot -Tpng -o dependency-graph.png

clean:
	rm -rf snapcraft-local.yaml build/

distclean: clean
	rm -rf vendor/

###############################################################################
###                                 Devdoc                                  ###
###############################################################################

build-docs:
	@cd docs && \
	while read p; do \
		(git checkout $${p} && npm install && VUEPRESS_BASE="/$${p}/" npm run build) ; \
		mkdir -p ~/output/$${p} ; \
		cp -r .vuepress/dist/* ~/output/$${p}/ ; \
		cp ~/output/$${p}/index.html ~/output ; \
	done < versions ;

sync-docs:
	cd ~/output && \
	echo "role_arn = ${DEPLOYMENT_ROLE_ARN}" >> /root/.aws/config ; \
	echo "CI job = ${CIRCLE_BUILD_URL}" >> version.html ; \
	aws s3 sync . s3://${WEBSITE_BUCKET} --profile terraform --delete ; \
	aws cloudfront create-invalidation --distribution-id ${CF_DISTRIBUTION_ID} --profile terraform --path "/*" ;
.PHONY: sync-docs


###############################################################################
###                           Tests & Simulation                            ###
###############################################################################

test: test-unit test-build

test-all: check test-race test-cover

test-unit:
	@VERSION=$(VERSION) go test -mod=readonly -tags='ledger test_ledger_mock' ./...

test-race:
	@VERSION=$(VERSION) go test -mod=readonly -race -tags='ledger test_ledger_mock' ./...

test-cover:
	@go test -mod=readonly -timeout 30m -race -coverprofile=coverage.txt -covermode=atomic -tags='ledger test_ledger_mock' ./...

test-build: build
	@go test -mod=readonly -p 4 `go list ./cli_test/...` -tags=cli_test -v

benchmark:
	@go test -mod=readonly -bench=. ./...


###############################################################################
###                                Linting                                  ###
###############################################################################

lint:
	golangci-lint run
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs gofmt -d -s

format:
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./client/lcd/statik/statik.go" | xargs gofmt -w -s
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./client/lcd/statik/statik.go" | xargs misspell -w
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./client/lcd/statik/statik.go" | xargs goimports -w -local github.com/cosmos/cosmos-sdk

###############################################################################
###                                Localnet                                 ###
###############################################################################

build-docker-sommnode:
	$(MAKE) -C networks/local

# Run a 4-node testnet locally
localnet-start: build-linux localnet-stop
	@if ! [ -f build/node0/somm/config/genesis.json ]; then docker run --rm -v $(CURDIR)/build:/somm:Z tendermint/sommnode testnet --v 4 -o . --starting-ip-address 192.168.10.2 --keyring-backend=test ; fi
	docker-compose up -d

# Stop testnet
localnet-stop:
	docker-compose down

test-docker:
	@docker build -f contrib/Dockerfile.test -t ${TEST_DOCKER_REPO}:$(shell git rev-parse --short HEAD) .
	@docker tag ${TEST_DOCKER_REPO}:$(shell git rev-parse --short HEAD) ${TEST_DOCKER_REPO}:$(shell git rev-parse --abbrev-ref HEAD | sed 's#/#_#g')
	@docker tag ${TEST_DOCKER_REPO}:$(shell git rev-parse --short HEAD) ${TEST_DOCKER_REPO}:latest

test-docker-push: test-docker
	@docker push ${TEST_DOCKER_REPO}:$(shell git rev-parse --short HEAD)
	@docker push ${TEST_DOCKER_REPO}:$(shell git rev-parse --abbrev-ref HEAD | sed 's#/#_#g')
	@docker push ${TEST_DOCKER_REPO}:latest

.PHONY: all build-linux install format lint \
	go-mod-cache draw-deps clean build \
	setup-transactions setup-contract-tests-data start-somm run-lcd-contract-tests contract-tests \
	test test-all test-build test-cover test-unit test-race \
	benchmark \
	build-docker-sommnode localnet-start localnet-stop \
	docker-single-node

###############################################################################
###                           Protobuf                                    ###
###############################################################################

proto-all: proto-format proto-lint proto-gen

proto-gen:
	@echo "Generating Protobuf files"
	$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace tendermintdev/sdk-proto-gen sh ./scripts/protocgen.sh

proto-format:
	@echo "Formatting Protobuf files"
	$(DOCKER) run --rm -v $(CURDIR):/workspace \
	--workdir /workspace tendermintdev/docker-build-proto \
	find ./ -not -path "./third_party/*" -name *.proto -exec .clang-format -i {} \;

proto-swagger-gen:
	@./scripts/protoc-swagger-gen.sh

proto-lint:
	@$(DOCKER_BUF) check lint --error-format=json

proto-check-breaking:
	@$(DOCKER_BUF) check breaking --against-input $(HTTPS_GIT)#branch=main

TM_URL           = https://raw.githubusercontent.com/tendermint/tendermint/v0.34.0/proto/tendermint
GOGO_PROTO_URL   = https://raw.githubusercontent.com/regen-network/protobuf/cosmos
COSMOS_PROTO_URL = https://raw.githubusercontent.com/regen-network/cosmos-proto/master
COSMOS_SDK_PROTO_URL = https://raw.githubusercontent.com/cosmos/cosmos-sdk/master/proto/cosmos/base
GOOGLE_PROTO_URL   = https://raw.githubusercontent.com/googleapis/googleapis/master/google/api
PROTOBUF_GOOGLE_URL = https://raw.githubusercontent.com/protocolbuffers/protobuf/master/src/google/protobuf

TM_CRYPTO_TYPES     = third_party/proto/tendermint/crypto
TM_ABCI_TYPES       = third_party/proto/tendermint/abci
TM_TYPES     	    = third_party/proto/tendermint/types
TM_VERSION 			= third_party/proto/tendermint/version
TM_LIBS				= third_party/proto/tendermint/libs/bits

GOGO_PROTO_TYPES    = third_party/proto/gogoproto
COSMOS_PROTO_TYPES  = third_party/proto/cosmos_proto
GOOGLE_PROTO_TYPES  = third_party/proto/google/api
PROTOBUF_GOOGLE_TYPES = third_party/proto/google/protobuf

SDK_ABCI_TYPES  	= third_party/proto/cosmos/base/abci/v1beta1
SDK_QUERY_TYPES  	= third_party/proto/cosmos/base/query/v1beta1
SDK_COIN_TYPES  	= third_party/proto/cosmos/base/v1beta1

proto-update-deps:
	mkdir -p $(GOGO_PROTO_TYPES)
	curl -sSL $(GOGO_PROTO_URL)/gogoproto/gogo.proto > $(GOGO_PROTO_TYPES)/gogo.proto

	mkdir -p $(COSMOS_PROTO_TYPES)
	curl -sSL $(COSMOS_PROTO_URL)/cosmos.proto > $(COSMOS_PROTO_TYPES)/cosmos.proto

	mkdir -p $(TM_ABCI_TYPES)
	curl -sSL $(TM_URL)/abci/types.proto > $(TM_ABCI_TYPES)/types.proto

	mkdir -p $(TM_VERSION)
	curl -sSL $(TM_URL)/version/types.proto > $(TM_VERSION)/types.proto

	mkdir -p $(TM_TYPES)
	curl -sSL $(TM_URL)/types/types.proto > $(TM_TYPES)/types.proto
	curl -sSL $(TM_URL)/types/evidence.proto > $(TM_TYPES)/evidence.proto
	curl -sSL $(TM_URL)/types/params.proto > $(TM_TYPES)/params.proto
	curl -sSL $(TM_URL)/types/validator.proto > $(TM_TYPES)/validator.proto

	mkdir -p $(TM_CRYPTO_TYPES)
	curl -sSL $(TM_URL)/crypto/proof.proto > $(TM_CRYPTO_TYPES)/proof.proto
	curl -sSL $(TM_URL)/crypto/keys.proto > $(TM_CRYPTO_TYPES)/keys.proto

	mkdir -p $(TM_LIBS)
	curl -sSL $(TM_URL)/libs/bits/types.proto > $(TM_LIBS)/types.proto

	mkdir -p $(SDK_ABCI_TYPES)
	curl -sSL $(COSMOS_SDK_PROTO_URL)/abci/v1beta1/abci.proto > $(SDK_ABCI_TYPES)/abci.proto

	mkdir -p $(SDK_QUERY_TYPES)
	curl -sSL $(COSMOS_SDK_PROTO_URL)/query/v1beta1/pagination.proto > $(SDK_QUERY_TYPES)/pagination.proto

	mkdir -p $(SDK_COIN_TYPES)
	curl -sSL $(COSMOS_SDK_PROTO_URL)/v1beta1/coin.proto > $(SDK_COIN_TYPES)/coin.proto

	mkdir -p $(GOOGLE_PROTO_TYPES)
	curl -sSL $(GOOGLE_PROTO_URL)/annotations.proto > $(GOOGLE_PROTO_TYPES)/annotations.proto
	curl -sSL $(GOOGLE_PROTO_URL)/http.proto > $(GOOGLE_PROTO_TYPES)/http.proto
	curl -sSL $(GOOGLE_PROTO_URL)/httpbody.proto > $(GOOGLE_PROTO_TYPES)/httpbody.proto

	mkdir -p $(PROTOBUF_GOOGLE_TYPES)
	curl -sSL $(PROTOBUF_GOOGLE_URL)/any.proto > $(PROTOBUF_GOOGLE_TYPES)/any.proto

PREFIX ?= /usr/local
BIN ?= $(PREFIX)/bin
UNAME_S ?= $(shell uname -s)
UNAME_M ?= $(shell uname -m)

BUF_VERSION ?= 0.11.0

PROTOC_VERSION ?= 3.11.2
ifeq ($(UNAME_S),Linux)
  PROTOC_ZIP ?= protoc-${PROTOC_VERSION}-linux-x86_64.zip
endif
ifeq ($(UNAME_S),Darwin)
  PROTOC_ZIP ?= protoc-${PROTOC_VERSION}-osx-x86_64.zip
endif

proto-tools: proto-tools-stamp buf

proto-tools-stamp:
	echo "Installing protoc compiler..."
	(cd /tmp; \
	curl -OL "https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/${PROTOC_ZIP}"; \
	unzip -o ${PROTOC_ZIP} -d $(PREFIX) bin/protoc; \
	unzip -o ${PROTOC_ZIP} -d $(PREFIX) 'include/*'; \
	rm -f ${PROTOC_ZIP})

	echo "Installing protoc-gen-gocosmos..."
	go install github.com/regen-network/cosmos-proto/protoc-gen-gocosmos

	# Create dummy file to satisfy dependency and avoid
	# rebuilding when this Makefile target is hit twice
	# in a row
	touch $@

buf: buf-stamp

buf-stamp:
	echo "Installing buf..."
	curl -sSL \
    "https://github.com/bufbuild/buf/releases/download/v${BUF_VERSION}/buf-${UNAME_S}-${UNAME_M}" \
    -o "${BIN}/buf" && \
	chmod +x "${BIN}/buf"

	touch $@

tools-clean:
	rm -f proto-tools-stamp buf-stampmodule

.PHONY: proto-all proto-gen proto-swagger-gen proto-format proto-lint proto-check-breaking proto-update-deps

#####################
# Integration tests #
#####################

#PACKAGES_E2E=$(shell go list ./... | grep '/integration_tests')
ORCHESTRATOR_IMAGE := "ghcr.io/peggyjv/gravity-bridge-orchestrator:v0.1.25"

e2e_build_images:
	@docker pull $(ORCHESTRATOR_IMAGE)
	@docker build -t sommelier:prebuilt -f Dockerfile .
	@docker build -t ethereum:prebuilt -f integration_tests/ethereum/Dockerfile integration_tests/ethereum/

e2e_clean_slate:
	@docker rm --force \
		$(shell docker ps -qa --filter="name=ethereum") \
		$(shell docker ps -qa --filter="name=sommelier") \
		$(shell docker ps -qa --filter="name=orchestrator") \
		1>/dev/null \
		2>/dev/null \
		|| true
	@docker wait \
		$(shell docker ps -qa --filter="name=ethereum") \
		$(shell docker ps -qa --filter="name=sommelier") \
		$(shell docker ps -qa --filter="name=orchestrator") \
		1>/dev/null \
		2>/dev/null \
		|| true
	@docker network rm testnet 1>/dev/null 2>/dev/null || true
	@cd integration_tests && go test -c

e2e_basic: e2e_clean_slate
	@integration_tests/integration_tests.test -test.run TestBasicChain -test.failfast -test.v || make -s fail

e2e_rebalance: e2e_clean_slate
	@integration_tests/integration_tests.test -test.run TestRebalance -test.failfast -test.v || make -s fail
