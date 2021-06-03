docker build -t pbatalov/finalclient .

docker push pbatalov/finalclient

ssh ec2-user@vanessasgh.me < deploy.sh