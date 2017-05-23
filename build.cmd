go build -o build\calculation-job-wrapper-0.1.7.exe -ldflags="-X main.VERSION=0.1.7 -X main.BUILD_TIME=20170523

SET GOOS=linux
SET GOARCH=amd64
go build -o build\calculation-job-wrapper-0.1.7-linux -ldflags="-X main.VERSION=0.1.7 -X main.BUILD_TIME=20170523

SET GOOS=darwin
SET GOARCH=amd64
go build -o build\calculation-job-wrapper-0.1.7-mac -ldflags="-X main.VERSION=0.1.7 -X main.BUILD_TIME=20170523
