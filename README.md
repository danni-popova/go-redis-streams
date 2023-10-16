# go-redis-streams
Go implementation of writing and reading from Redis Streams using multiple Consumer groups. 

**How to run**

You need to have Redis running locally and the easiest way to do this is to spin up a Docker container using this command:
```shell
docker run -d --name redis-stack -p 6379:6379 -p 8001:8001 redis/redis-stack:latest
```

