docker rm -f finalclient

docker pull pbatalov/finalclient

docker run \
    -d \
    --name finalclient \
    -p 80:80 -p 443:443 \
    -v /etc/letsencrypt:/etc/letsencrypt:ro \
    pbatalov/finalclient