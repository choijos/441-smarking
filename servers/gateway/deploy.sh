# Calls the build script you created in the previous step to rebuild the API server Linux executable and API docker container image
# Pushes your API server Docker container image to Docker Hub
# Executes the appropriate docker rm -f, docker pull, and docker run commands on your API server
#   VM via ssh to stop the existing container instance, pull the updated container image from DockerHub,
#   and re-run a newly-updated instance.

# docker push choijos/choijos.me

# docker rm -f choijos.me

# docker pull choijos/choijos.me

# docker run \
#   -d \
#   -p 443:443 \
#   -v /etc/letsencrypt:/etc/letsencrypt:ro \
#   -e TLSCERT=/etc/letsencrypt/live/choijos.me/fullchain.pem \
#   -e TLSKEY=/etc/letsencrypt/live/choijos.me/privkey.pem \
#   --name choijos.me \
#   choijos/choijos.me

# ssh ec2-user@ec2-35-81-89-251.us-west-2.compute.amazonaws.com < build.sh
sh build.sh

ssh ec2-user@ec2-35-81-89-251.us-west-2.compute.amazonaws.com < clean.sh