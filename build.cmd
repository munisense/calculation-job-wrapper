go build -o build\calculation-job-wrapper-0.1.5.exe -ldflags="-X main.VERSION=0.1.5 -X main.BUILD_TIME=20170518

SET GOOS=linux
SET GOARCH=amd64
go build -o build\calculation-job-wrapper-0.1.5-linux -ldflags="-X main.VERSION=0.1.5 -X main.BUILD_TIME=20170518

SET GOOS=darwin
SET GOARCH=amd64
go build -o build\calculation-job-wrapper-0.1.5-mac -ldflags="-X main.VERSION=0.1.5 -X main.BUILD_TIME=20170518
