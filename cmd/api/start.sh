#!/bin/bash

API_DIR="/home/tmalik/git/archive_utility/cmd/api"
API_EXE="./api"

cd "$API_DIR"
if [ -f "$API_EXE" ] && [ -x "$API_EXE" ]; then
    echo "located api_exe"
    sudo rm "$API_EXE"
    echo "removed executable"
else
    echo "api_exe does not exist or is not executable"
fi

echo "archive_api: removing build..."
cd "$API_DIR"

echo "archive_api: building service..."
cd "$API_DIR"
go build .

echo "archive_api: starting service..."
cd "$API_DIR"
./api
