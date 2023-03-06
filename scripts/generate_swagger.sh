#!/bin/bash
alias swagger='docker run --rm -it  --user $(id -u):$(id -g) -e GOCACHE=/tmp -e  GOPATH=$(go env GOPATH):/go -v $HOME:$HOME -w $(pwd) quay.io/goswagger/swagger'
swagger generate spec -o ./swagger.yaml --scan-models
