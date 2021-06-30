#!/usr/bin/env bash

echo "$(tput setaf 3)""building secretserviced ...""$(tput sgr0)"
go build -race -o secretserviced cmd/app/secretserviced/main.go
echo "$(tput setaf 3)""building secretservice ...""$(tput sgr0)"

go build -race -o secretservice cmd/app/secretservice/main.go
du -bh secretservice*