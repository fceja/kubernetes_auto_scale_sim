# Creating Docker Image

## Docker Compose

\*\* Note - to use external secrets with Docker Compose, you must deploy the stack in Docker Swarm mode.

\*\* Make sure Docker swarm is initialized and then use `docker stack deploy` instead of `docker compose up -d` (full command in All-in-one)

`docker compose up -d`

## Build Image

`docker build -t kafka_producer_image .`

## Deploy Stack

\*\* Note - Need to build image before deploying to stack.

`docker stack deploy -c docker-compose.yaml my_stack`

## All-in-one

Delete previous stack, build docker image, deploy stack from rebuilt image
`docker stack rm my_stack && docker build -t kafka_producer . && docker stack deploy -c docker-compose.yaml my_stack`
