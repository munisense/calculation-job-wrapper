go build -o build\calculation-job-wrapper-0.1.6.exe -ldflags="-X main.VERSION=0.1.6 -X main.BUILD_TIME=20170521

SET GOOS=linux
SET GOARCH=amd64
go build -o build\calculation-job-wrapper-0.1.6-linux -ldflags="-X main.VERSION=0.1.6 -X main.BUILD_TIME=20170521

SET GOOS=darwin
SET GOARCH=amd64
go build -o build\calculation-job-wrapper-0.1.6-mac -ldflags="-X main.VERSION=0.1.6 -X main.BUILD_TIME=20170521
