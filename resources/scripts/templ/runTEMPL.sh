#!/bin/bash

go mod tidy
go get github.com/a-h/templ/cmd/templ@latest

go run github.com/a-h/templ/cmd/templ@latest generate ./app/views
