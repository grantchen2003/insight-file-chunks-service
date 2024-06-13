#!/bin/bash

pids=$(netstat -aon | grep ':50051' | awk '{print $5}' | cut -d ':' -f 1)

# Loop through each PID and kill the corresponding process
for pid in $pids; do
    taskkill /PID $pid /F
done

cd ..

export ENV=dev

nodemon --exec "go run cmd/filechunks/main.go" --watch . -e go