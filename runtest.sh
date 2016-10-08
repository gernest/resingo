#!/bin/bash

if [ -e  ./.env ];then
	echo "sourcing wnvironment vars from file"
	source .env
	go test -v
else
	go test -v
fi
