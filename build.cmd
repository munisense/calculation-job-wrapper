go build -o build\calculation-job-wrapper-0.1.9.exe -ldflags="-X main.VERSION=0.1.9 -X main.BUILD_TIME=20170530

SET GOOS=linux
SET GOARCH=amd64
go build -o build\calculation-job-wrapper-0.1.9-linux -ldflags="-X main.VERSION=0.1.9 -X main.BUILD_TIME=20170530

SET GOOS=darwin
SET GOARCH=amd64
go build -o build\calculation-job-wrapper-0.1.9-mac -ldflags="-X main.VERSION=0.1.9 -X main.BUILD_TIME=20170530
