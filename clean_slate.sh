#!/bin/bash

set -e

echo Cleaning up test environment
docker rm --force \
       $(shell docker ps -qa --filter="name=ethereum") \
       $(shell docker ps -qa --filter="name=sommelier") \
       $(shell docker ps -qa --filter="name=orchestrator") \
       1>/dev/null \
       2>/dev/null \
       || true
docker wait \
       $(shell docker ps -qa --filter="name=ethereum") \
       $(shell docker ps -qa --filter="name=sommelier") \
       $(shell docker ps -qa --filter="name=orchestrator") \
       1>/dev/null \
       2>/dev/null \
       || true
docker network prune --force 1>/dev/null 2>/dev/null || true
cd integration_tests && go test -c
cd -
