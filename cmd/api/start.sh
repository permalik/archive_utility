#!/bin/bash

echo "archive_api: removing build..."
rm api

echo "archive_api: building service..."
go build .

echo "archive_api: starting service..."
./api
