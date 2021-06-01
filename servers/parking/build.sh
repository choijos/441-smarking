docker build -t pbatalov/parking .


docker push pbatalov/parking

ssh ec2-user@api.vanessasgh.me < deploy.sh

