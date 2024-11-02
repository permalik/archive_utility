#!/bin/bash

echo "archive_archive: removing build..."
rm archive

echo "archive_archive: building service..."
go build .

echo "archive_archive: executing service..."
./archive
