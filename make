#! /bin/bash

echo "Building ======================================================================="
go generate ./... \
    && go test -coverprofile=coverage.out ./... \
    && go build -ldflags "-X github.com/jacobsimpson/msh/interpreter/builtin.Version=`date -u +%Y%m%d.%H%M%S`" . \
    && echo "Success."
