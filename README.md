# Redis Go Sample


# Redis

## Run
docker run -p 6379:6379 --name my-redis -d redis

## Remove
docker rm -f my-redis

## Run Redis CLI
docker run --rm --name redis-cli -it redis:alpine redis-cli -h localhost
