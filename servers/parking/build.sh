docker build -t pbatalov/messaging .


docker push pbatalov/messaging

ssh ec2-user@api.pavelsrinidhi.me < deploy.sh

