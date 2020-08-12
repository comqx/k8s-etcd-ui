#!/usr/bin/env sh

npm run build 
go-bindata -o tpls.go ../dist/...

sed -i "" "s/package main/package tpls/g" tpls.go
