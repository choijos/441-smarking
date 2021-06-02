

docker pull pbatalov/parking

export MESSAGESADDR=:80
export TWILIO_ACCOUNT_SID=ACbf391b4d4234222e10ca877559fc0838
export TWILIO_AUTH_TOKEN=6c500f3f30f434a30543cc50d2599406
export TWILIO_PHONE_NUMBER=+12512734782

docker rm -f parking

docker run -d \
    -e PARKINGADDR=$PARKINGADDR \
    -e TWILIO_ACCOUNT_SID=$TWILIO_ACCOUNT_SID \
    -e TWILIO_AUTH_TOKEN=$TWILIO_AUTH_TOKEN \
    -e TWILIO_PHONE_NUMBER=$TWILIO_PHONE_NUMBER \
    --name parking \
    --network info441network \
    pbatalov/parking
exit