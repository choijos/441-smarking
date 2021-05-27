docker rm -f customMongoContainer


docker run -d \
  --name customMongoContainer \
  --network customNet \
  mongo