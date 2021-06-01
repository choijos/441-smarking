GOOS=linux go build
docker build -t pbatalov/finalgateway .
go clean

docker push pbatalov/finalgateway

exit
# docker build -t choijos/choijos.me .
# go clean

# docker push choijos/choijos.me

# exit

# unused - moved to another file
# ssh ec2-user@ec2-35-81-89-251.us-west-2.compute.amazonaws.com < deploy.sh