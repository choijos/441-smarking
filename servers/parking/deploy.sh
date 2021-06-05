docker pull pbatalov/parking

export PARKINGOWNADDR=:80
export TWILIO_ACCOUNT_SID=ACbf391b4d4234222e10ca877559fc0838
export TWILIO_PHONE_NUMBER=+12512734782
export TWILIO_AUTH_TOKEN=cf1117117861ffb99e70e80c8457a1e7

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