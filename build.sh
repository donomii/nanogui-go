#!/bin/sh
rm dash
go build -v -o dash dashboard/dashboard.go
go build -v -o newfs newfs/newfs.go
./newfs
