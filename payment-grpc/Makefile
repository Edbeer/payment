build:
	@go build -o ./bin/api

run: build
	@./bin/api

docker:
	docker run --name paymentdb \
	-e POSTGRES_HOST=localhost \
	-e POSTGRES_PASSWORD=postgres \
	-e POSTGRES_USER=postgres \
	-e POSTGRES_DB=paymentdb \
	-p 5436:5432 -d postgres:15.1

docker-start:
	docker start paymentdb

docker-exec: docker-start
	docker exec -it paymentdb psql -U postgres paymentdb

docker-stop:
	docker stop paymentdb

# migrations
migrate-create:
	migrate create -ext sql -dir ./migrations -seq paymentdb

migrate-up:
	migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5436/paymentdb?sslmode=disable" up

migrate-down:
	migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5436/paymentdb?sslmode=disable" down