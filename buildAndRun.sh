#!/bin/bash

go build -o ./main ./main.go

cd ./frontend

npm run build

cd ..

./main