# Fiber Gorm Postgresql Boilerplate
ON PROGRESS

## Requirements
- Go ver 1.19
- Postgres ver 14

## Instalation (for development)
1. install swaggo `go install github.com/swaggo/swag/cmd/swag@latest`
1. install golang migrate `go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest`
1. copy .env.example to .env and fill based on your postgres configuration and set env as "dev"
1. generate/refresh swagger.json `swag init`
1. migrate database `go run main.go migrate-db up`
1. init superuser `go run main.go init-superuser --username {username} --email {email} --password {password}`
1. run server `go run main.go runserver`
1. open swagger "http://{SERVER_HOST}:{SERVER_PORT}/docs/index.html"
1. login using username and password

## Instalation (for Production)
TODO

## Golang migrate
### Install golang migrate
refrences https://github.com/golang-migrate/migrate/tree/master/cmd/migrate
- install golang migrate cli `go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest`
### create new migration
1. `migrate create -ext sql --dir ./migrations/migrations_files/ {migration name}`
2. Edit newly created up.sql and down.sql
### Migrate database
#### Using golang migrate
- change every {} based on your postgresql configuration 
`migrate -source file://migrations/migrations_files/ -database postgres://{username}:{password}@{host}:{port}/{database}?sslmode={require/verify-full/verify-ca/disable} up`
#### Using cli
- see `go run main.go migrate-db --help`

## Testing

- run all testing `go test ./...`
- run all test in folder `go test ./{folder name}/... ./{another folder name}/...`
- run all test in file `go test ./{folder name}/{file name}`
- run specific test function `go test --run '^{function name}$' ./{folder name}/{file name}` (Note: --run input is a regex)
- run test verbosely (show log) `go test ./... -v`
- remove all test cache `go clean -testcache`
- disable parallel test `go test ./... -p 1`
