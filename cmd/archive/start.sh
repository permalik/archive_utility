#!/bin/bash

echo "archive_archive: removing build..."
ARCHIVE_DIR="/home/tmalik/git/archive_utility/cmd/archive"
ARCHIVE_EXE="./archive"

cd "$ARCHIVE_DIR"
if [ -f "$ARCHIVE_EXE" ] && [ -x "$ARCHIVE_EXE" ]; then
    echo "located archive_exe"
    sudo rm "$ARCHIVE_EXE"
    echo "removed executable"
else
    echo "archive_exe does not exist or is not executable"
fi

echo "archive_archive: building service..."
cd "$ARCHIVE_DIR"
go build .

echo "archive_archive: executing service..."
cd "$ARCHIVE_DIR"
./archive
