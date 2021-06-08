docker pull pbatalov/parking

export PARKINGOWNADDR=:80

docker rm -f parking

docker run -d \
    -e PARKINGADDR=$PARKINGOWNADDR \
    -e TWILIO_ACCOUNT_SID=$TWILIO_ACCOUNT_SID \
    -e TWILIO_AUTH_TOKEN=$TWILIO_AUTH_TOKEN \
    -e TWILIO_PHONE_NUMBER=$TWILIO_PHONE_NUMBER \
    --name parking \
    --network info441network \
    pbatalov/parking
exit