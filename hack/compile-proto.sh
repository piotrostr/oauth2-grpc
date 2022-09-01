#!/bin/bash

buf mod update
buf build
buf generate
go mod tidy
