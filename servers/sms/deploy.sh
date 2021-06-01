docker rm -f notifications
docker pull pbatalov/notifications

export NOTIFICATIONSADDR=:80
export TWILIO_ACCOUNT_SID=ACbf391b4d4234222e10ca877559fc0838
export TWILIO_AUTH_TOKEN=d07e3e6bfe47f2c0c579cfe2db2d13b8
export TWILIO_PHONE_NUMBER=+12512734782

docker run -d \
  -e NOTIFICATIONSADDR=$NOTIFICATIONSADDR \
  --name notifications \
  --network info441network \
  -e TWILIO_ACCOUNT_SID=$TWILIO_ACCOUNT_SID \
  -e TWILIO_AUTH_TOKEN=$TWILIO_AUTH_TOKEN \
  -e TWILIO_PHONE_NUMBER=$TWILIO_PHONE_NUMBER \
  pbatalov/notifications

exit