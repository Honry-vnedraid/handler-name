# handler-name

## Prerequisites

```shell
brew install go
brew install docker
brew install golang-migrate
```

## How to deploy

```shell
docker-compose up -d
```

```shell
go build
```

```shell
migrate -path db/migrations \
  -database "postgres://postgres:supersecret@localhost:5433/newsdb?sslmode=disable" up
```

```shell
./handler-service
```
