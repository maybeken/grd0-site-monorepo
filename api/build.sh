#!/bin/bash
rm -rf build/
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -tags lambda.norpc -o build/bootstrap server.go
cp -r data build/data
cd build
zip -x data/data.go -r bootstrap.zip bootstrap data
