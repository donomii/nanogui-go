#!/bin/sh
rm dash
go build -v -o dash dashboard/dashboard.go
./dash
