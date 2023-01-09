# go-redis-practice

Play with Redis in Golang

## Reference

[1] Redis in Action - 2013 - Manning

[2] [Redis in Action Github](https://github.com/josiahcarlson/redis-in-action)

[3] [Redis Pub/Sub underhood](https://making.pusher.com/redis-pubsub-under-the-hood/)

## Running

### Using Docker

- check the config messages in the config/config.go file first, you may need to set your config for redis
- run `docker-compose up -d` in the directory.
- use `docker exec -it go-redis-practice go test ./chap_0*/redis_test.go -v` to run the test, use number 1 through 8 to replace the `*` depending on the Chapter's examples you want to run.
