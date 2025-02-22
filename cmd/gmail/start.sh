#!/bin/bash

GMAIL_DIR="/home/tmalik/git/archive_utility/cmd/gmail"
GMAIL_EXE="./gmail"

cd "$GMAIL_DIR"
echo "removing build.."
if [ -f "$GMAIL_EXE" ] && [ -x "$GMAIL_EXE" ]; then
    echo "located gmail_exe"
    rm "$GMAIL_EXE"
    echo "removed executable"
else
    echo "gmail_exe does not exist or is not executable"
fi

echo "building service.."
cd "$GMAIL_DIR"
go build .

echo "executing service.."
cd "$GMAIL_DIR"
./gmail
