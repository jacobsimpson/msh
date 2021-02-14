#! /bin/bash

go generate . \
    && go test -coverprofile=coverage.out ./... \
    && go build . \
    && echo "Success."
