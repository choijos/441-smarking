# DB & network

docker network create info441network


docker rm -f redisServer
docker run -d --name redisServer --network info441network redis

export MYSQL_ROOT_PASSWORD=thisbetterwork
export DB_NAME=userinfo

docker pull pbatalov/finaldb
docker rm -f finaldb

docker run -d \
  -p 3306:3306 \
  --name finaldb \
  --network info441network \
  -e MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD \
  -e MYSQL_DATABASE=$DB_NAME \
  pbatalov/finaldb


# Gateway
docker rm -f finalgateway
docker pull pbatalov/finalgateway

export TLSCERT=/etc/letsencrypt/live/api.vanessasgh.me/fullchain.pem
export TLSKEY=/etc/letsencrypt/live/api.vanessasgh.me/privkey.pem
export SESSIONKEY=$(openssl rand -base64 18)
export REDISADDR=redisServer:6379
export DSN=root:$MYSQL_ROOT_PASSWORD@tcp\(sqla4:3306\)/$DB_NAME
export PARKINGADDR=http://parking:80

docker run -d \
  --name finalgateway \
  --network info441network \
  --restart unless-stopped \
  -p 443:443 \
  -v /etc/letsencrypt:/etc/letsencrypt:ro \
  -e TLSCERT=$TLSCERT \
  -e TLSKEY=$TLSKEY \
  -e SESSIONKEY=$SESSIONKEY \
  -e REDISADDR=$REDISADDR \
  -e DSN=$DSN \
  -e PARKINGADDR=$PARKINGADDR \
  pbatalov/finalgateway

exit
