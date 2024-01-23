#!/bin/sh
set -e
cd "$( dirname -- "$( readlink -f $0 )"; )"
rm -rf covdatafiles
mkdir covdatafiles
GOCOVERDIR=covdatafiles
go build -cover ../cmd/struc-gen
GOCOVERDIR=covdatafiles ./struc-gen -file example.go
GOCOVERDIR=covdatafiles ./struc-gen -file bitfield_example.go
GOCOVERDIR=covdatafiles ./struc-gen -file default_byteoder_example.go -little
go test -coverprofile=covdatafiles/coverage.out ../cmd/...
go tool covdata textfmt -i=covdatafiles -o covdatafiles/profile.cov
