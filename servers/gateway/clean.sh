# DB & network
docker rm -f redisServer
docker run -d --name redisServer --network info441network redis

export MYSQL_ROOT_PASSWORD=thisbetterwork
export DB_NAME=userinfo

docker pull choijos/sqla4
docker rm -f sqla4

docker run -d \
  -p 3306:3306 \
  --name sqla4 \
  --network info441network \
  -e MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD \
  -e MYSQL_DATABASE=$DB_NAME \
  choijos/sqla4


# Gateway
docker rm -f choijos.me
docker pull choijos/choijos.me

export TLSCERT=/etc/letsencrypt/live/api.choijos.me/fullchain.pem
export TLSKEY=/etc/letsencrypt/live/api.choijos.me/privkey.pem
export SESSIONKEY=$(openssl rand -base64 18)
export REDISADDR=redisServer:6379
export DSN=root:$MYSQL_ROOT_PASSWORD@tcp\(sqla4:3306\)/$DB_NAME
export MESSAGESADDR=http://messaging:80
export SUMMARYADDR=http://summary:80

docker run -d \
  --name choijos.me \
  --network info441network \
  --restart unless-stopped \
  -p 443:443 \
  -v /etc/letsencrypt:/etc/letsencrypt:ro \
  -e TLSCERT=$TLSCERT \
  -e TLSKEY=$TLSKEY \
  -e SESSIONKEY=$SESSIONKEY \
  -e REDISADDR=$REDISADDR \
  -e DSN=$DSN \
  -e MESSAGESADDR=$MESSAGESADDR \
  -e SUMMARYADDR=$SUMMARYADDR \
  choijos/choijos.me

exit
