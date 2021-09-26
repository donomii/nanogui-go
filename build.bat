rm dash
rm nfs.exe
#go build -v -o dash dashboard/dashboard.go
go build -v -o nfs.exe newfs/newfs.go
go build -v -o sample.exe sample1/sample.go
nfs.exe
