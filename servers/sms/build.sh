docker build -t pbatalov/notifications .

docker push -t pbatalov/notifications

ssh ec2-user@api.vanessasgh.me < deploy.sh