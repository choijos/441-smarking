

docker pull pbatalov/parking

export PARKINGOWNADDR=:80

docker rm -f parking

docker run -d \
    -e PARKINGADDR=$PARKINGOWNADDR \
    --name parking \
    --network info441network \
    pbatalov/parking
exit