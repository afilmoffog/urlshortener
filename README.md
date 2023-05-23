# Url shortener service



### Contains:
* ##### Golang net/http based server
* ##### Tarantool database
* ##### Docker

Docker-compose includes 3 services: server, go test service, tarantool db.
For local development might take .env file with variables, otherwise takes environment variables.
It runs tests inside docker container.

##### How to run:
```shell script
# Build images
cd /build && docker-compose build
# Run images
docker-compose up
```
Go to POST http://127.0.0.1:8080/?source=<original_url> and take <response string>
After, GET http://127.0.0.1:8080/<response string> and you'll be redirected to original url
