GOOS=linux go build
docker build -t choijos/sqla4 .
go clean

docker push choijos/sqla4

exit