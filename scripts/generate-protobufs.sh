#!/bin/bash

cd ../internal/protobufs

protoc *.proto --go_out=.. --go-grpc_out=..