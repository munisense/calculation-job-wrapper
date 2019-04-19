go build -o build\calculation-job-wrapper-0.1.11.exe -ldflags="-X main.VERSION=0.1.11 -X main.BUILD_TIME=20170110

SET GOOS=linux
SET GOARCH=amd64
go build -o build\calculation-job-wrapper-0.1.11-linux -ldflags="-X main.VERSION=0.1.11 -X main.BUILD_TIME=20170110

SET GOOS=darwin
SET GOARCH=amd64
go build -o build\calculation-job-wrapper-0.1.11-mac -ldflags="-X main.VERSION=0.1.11 -X main.BUILD_TIME=20170110

docker build -t munisense/tools-calculation-job-wrapper:latest .