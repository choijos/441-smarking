docker rm -f customMongoContainer


docker run -d \
  --name customMongoContainer \
  --network info441network \
  mongo