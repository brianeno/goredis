# Redis Go Sample


# Redis
docker run -p 6379:6379 --name my-redis -d redis

docker rm -f my-redis

docker run --rm -it redis:alpine redis-cli -h localhost
