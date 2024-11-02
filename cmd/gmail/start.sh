#!/bin/bash

echo "archive_gmail: removing build..."
rm gmail

echo "archive_gmail: building service..."
go build .

echo "archive_gmail: executing service..."
./gmail
