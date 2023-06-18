@echo off
set GOOS=js
set GOARCH=wasm
go build  -trimpath -ldflags="-s -w" -o wasm/bin/argon.wasm ./src