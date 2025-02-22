#!/bin/bash

ARCHIVE_DIR="/home/tmalik/git/archive_utility/cmd/archive"
ARCHIVE_EXE="./archive"

echo "removing build.."
cd "$ARCHIVE_DIR"
if [ -f "$ARCHIVE_EXE" ] && [ -x "$ARCHIVE_EXE" ]; then
    echo "located archive_exe"
    rm "$ARCHIVE_EXE"
    echo "removed executable"
else
    echo "archive_exe does not exist or is not executable"
fi

echo "building service.."
cd "$ARCHIVE_DIR"
go build .

echo "executing service.."
cd "$ARCHIVE_DIR"
./archive
