

docker pull pbatalov/parking

export MESSAGESADDR=:80

docker rm -f parking

docker run -d \
    -e PARKINGADDR=$PARKINGADDR \
    --name parking \
    --network info441network \
    pbatalov/parking
exit