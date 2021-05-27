

docker pull pbatalov/messaging

export MESSAGESADDR=:80

docker rm -f messaging

docker run -d \
    -e MESSAGESADDR=$MESSAGESADDR \
    --name messaging \
    --network customNet \
    pbatalov/messaging
exit