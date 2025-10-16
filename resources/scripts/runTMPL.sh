#!/bin/bash

go mod tidy
go get github.com/a-h/templ
go run github.com/a-h/templ/cmd/templ@latest generate
