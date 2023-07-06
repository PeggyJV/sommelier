#!/usr/bin/make -f

PACKAGES_SIMTEST=$(shell go list ./... | grep '/simulation')
VERSION := $(shell git describe --abbrev=6 --dirty --always --tags)
COMMIT := $(shell git log -1 --format='%H')
LEDGER_ENABLED ?= true
SDK_PACK := $(shell go list -m github.com/cosmos/cosmos-sdk | sed  's/ /\@/g')
BUILDDIR ?= $(CURDIR)/build
BINDIR ?= $(GOPATH)/bin
SIMAPP = ./app

# for dockerized protobuf tools
DOCKER := $(shell which docker)
#BUF_IMAGE=bufbuild/buf@sha256:3cb1f8a4b48bd5ad8f09168f10f607ddc318af202f5c057d52a45216793d85e5 #v1.4.0
BUF_IMAGE=bufbuild/buf:0.48.2
DOCKER_BUF := $(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace $(BUF_IMAGE)
HTTPS_GIT := https://github.com/peggyjv/sommelier.git

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



test: test-unit test-build

test-all: check test-race test-cover

test-unit:
	@VERSION=$(VERSION) go test -mod=readonly -tags='ledger test_ledger_mock' $(shell go list ./... | grep -v '/integration_tests')

test-race:
	@VERSION=$(VERSION) go test -mod=readonly -race -tags='ledger test_ledger_mock' ./...

test-cover:
	@go test -mod=readonly -timeout 30m -race -coverprofile=coverage.txt -covermode=atomic -tags='ledger test_ledger_mock' ./...

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


###############################################################################
###                           Protobuf                                    ###
###############################################################################
protoVer=0.11.6
protoImageName=ghcr.io/cosmos/proto-builder:$(protoVer)
protoImage=$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace $(protoImageName)

proto-all: proto-format proto-lint proto-gen format

proto-gen:
	@echo "Generating Protobuf files"
	@$(protoImage) sh ./scripts/protocgen.sh

proto-format:
	@echo "Formatting Protobuf files"
	@$(protoImage) find ./ -name "*.proto" -exec clang-format -i {} \;

proto-swagger-gen:
	@./scripts/protoc-swagger-gen.sh

proto-lint:
	@$(DOCKER_BUF) lint --error-format=json

proto-check-breaking:
	@$(DOCKER_BUF) breaking --against $(HTTPS_GIT)#branch=main

proto-version:
	@$(protoImage) buf --version

#####################
# Integration tests #
#####################

ORCHESTRATOR_IMAGE := "ghcr.io/peggyjv/gravity-bridge-orchestrator:latest"

e2e_build_images: e2e_clean_slate
	@docker pull $(ORCHESTRATOR_IMAGE)
	@docker tag $(ORCHESTRATOR_IMAGE) orchestrator:prebuilt
	@docker build -t sommelier:prebuilt -f Dockerfile .
	@docker build -t ethereum:prebuilt -f integration_tests/ethereum/Dockerfile integration_tests/ethereum/

e2e_clean_slate:
	@./clean_slate.sh

e2e_basic: e2e_clean_slate
	@integration_tests/integration_tests.test -test.run TestBasicChain -test.failfast -test.v || make -s fail

e2e_cork_test: e2e_clean_slate
	@E2E_SKIP_CLEANUP=true integration_tests/integration_tests.test -test.failfast -test.v -test.run IntegrationTestSuite -testify.m TestCork || make -s fail

e2e_scheduled_cork_test: e2e_clean_slate
	@E2E_SKIP_CLEANUP=true integration_tests/integration_tests.test -test.failfast -test.v -test.run IntegrationTestSuite -testify.m TestScheduledCork || make -s fail

e2e_incentives_test: e2e_clean_slate
	@E2E_SKIP_CLEANUP=true integration_tests/integration_tests.test -test.failfast -test.v -test.run IntegrationTestSuite -testify.m TestIncentives || make -s fail

fail:
	@echo 'test failed; dumping container logs into ./testlogs for review'
	@docker logs ethereum > testlogs/ethereum.log 2>&1 || true
	@docker logs sommelier0 > testlogs/sommelier0.log 2>&1 || true
	@docker logs sommelier1 > testlogs/sommelier1.log 2>&1 || true
	@docker logs sommelier2 > testlogs/sommelier2.log 2>&1 || true
	@docker logs sommelier3 > testlogs/sommelier3.log 2>&1 || true
	@docker logs orchestrator0 > testlogs/orchestrator0.log 2>&1 || true
	@docker logs orchestrator1 > testlogs/orchestrator1.log 2>&1 || true
	@docker logs orchestrator2 > testlogs/orchestrator2.log 2>&1 || true
	@docker logs orchestrator3 > testlogs/orchestrator3.log 2>&1 || true
	@false
