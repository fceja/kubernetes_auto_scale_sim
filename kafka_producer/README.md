# Creating Docker Image

## Build Image

`docker build -t kafka_producer_image .`

## Deploy Stack

`docker stack deploy -c docker-compose.yml my_stack`
