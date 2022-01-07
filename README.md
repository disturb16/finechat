# Finechat - the financial chat
This is a simple web chat application in go.

# Setup

## Docker deployment
the application is dockerized to make everything more simple. In order to run in docker, first you need to initialize in swarm mode.

`docker swarm init`

Then you need to build the docker image of the application from the root path.

`docker build -t finechat .`

After this you can deploy the stack with the following.

`docker stack deploy -c docker-compose.yml finechat`

The application will now be running on port 8081 after it initilizes.

## Removing the stack

To stop everything you can run the follwing command.

`docker stack rm finechat`

And to leave the swarm mode just type:

`docker swarm leave`
