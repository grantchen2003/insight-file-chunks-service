#!/bin/bash

cd ..

export ENV=dev

nodemon --exec "go run cmd/filechunks/main.go" --watch . -e go