GOOS=linux go build
docker build -t pbatalov/finalgateway .
go clean

docker push pbatalov/finalgateway

exit