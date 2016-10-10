#!/bin/bash

if [ -e  ./.env ];then
	echo "sourcing wnvironment vars from file"
	source .env
	#go test -v
#else
	#go test -v
fi
case "$1" in
	"all")
		go test -v
		;;
	"run")
		go test -v -run $2
		;;
	*)
		go test -v
esac
