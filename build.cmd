go build -o build\correlation-service-wrapper-0.1.3.exe -ldflags="-X main.VERSION=0.1.3 -X main.BUILD_TIME=20170504

SET GOOS=linux
SET GOARCH=amd64
go build -o build\correlation-service-wrapper-0.1.3-linux -ldflags="-X main.VERSION=0.1.3 -X main.BUILD_TIME=20170504

SET GOOS=darwin
SET GOARCH=amd64
go build -o build\correlation-service-wrapper-0.1.3-mac -ldflags="-X main.VERSION=0.1.3 -X main.BUILD_TIME=20170504
