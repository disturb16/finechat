# Finechat - The financial chat
This is a simple web chat application in go.

# Setup

## Docker deployment
the application is dockerized to make everything more simple. In order to run in docker, first you need to initialize in swarm mode.

`docker swarm init`

Then you need to build the docker images from the root path.

The app:

`docker build -t finechatapp ./finechatapp`

The bot:

`docker build -t finechatbot ./finechatbot`

You can also download the depenencies before deploying the stack:

- `docker pull mariadb:10.6.4`
- `docker pull rabbitmq:3.9.12-management`

After this you can deploy the stack with the following.

`docker stack deploy -c docker-compose.yml finechat`

Wait for it to initialize completely, and visit http://localhost:8081.

## Removing the stack

To stop everything you can run the following command.

`docker stack rm finechat`

And to leave the swarm mode just type:

`docker swarm leave`

# Finechat bot
Process the stock commands (e.g. "/stock=googl.us"). The finechatbot is a microservice that listenes to the `stock.command` topic and sends back the response or error.

# Database
This project uses mariadb as the database engine. The `docker-compose.yml` file loads the `finechatapp/finechat_schema.sql` script to initialize the database schema and creates the appuser which is used by the api to connect.

# Browser Support

- Chrome
- Firefox
- Edge

# Project structure

This repo contains two applications, the app and the bot. Each one is located on its own folder `finechatapp` and `finechatbot`. The reason for this is that if one of the containers crashes the other wont be affected.

- The **finechatapp** is responsable of the web application, users, messages and chatrooms.
- The **finechatbot** is responsable to process any stock command message.

The frontend app is written in Vuejs and is located in `finechatapp/internal/web` which combined with golang embed files, makes working with SPAs really easy.
