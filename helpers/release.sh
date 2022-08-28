#!/bin/bash

mkdir release-$1

echo "Creating release $1"

echo "Building backend"
GOOS=windows GOARCH=amd64 go build -o ./release-$1/main-amd64.exe ./main.go
echo "Built Windows x64"
GOOS=darwin GOARCH=amd64 go build -o ./release-$1/main-amd64-darwin ./main.go
echo "Built Mac x64"
GOOS=linux GOARCH=amd64 go build -o ./release-$1/main-amd64-linux ./main.go
echo "Built Linux x64"

cp config.toml ./release-$1/config.toml

cd frontend

echo "Building frontend"
npm run build
echo "Built frontend"

cd ..
mkdir release-$1/frontend

cp -r frontend/dist release-$1/frontend

echo "Creating zips"
zip -r release-$1/release-$1-windows-x64.zip release-$1/frontend release-$1/main-amd64.exe release-$1/config.toml
zip -r release-$1/release-$1-mac-x64.zip release-$1/frontend release-$1/main-amd64-darwin release-$1/config.toml
zip -r release-$1/release-$1-linux-x64.zip release-$1/frontend release-$1/main-amd64-linux release-$1/config.toml
