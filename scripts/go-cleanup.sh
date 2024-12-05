#!/bin/bash

if [[ -n "$1" ]]; then
  cd "$1" || exit 1
  gofmt -l -w . && go mod tidy
else
  echo "no argumenets were provided. arg 1 should be the path of the go code"
  exit 1
fi
