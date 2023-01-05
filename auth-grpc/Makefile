docker:
	docker run --name authdb \
	-e POSTGRES_HOST=localhost \
	-e POSTGRES_PASSWORD=postgres \
	-e POSTGRES_USER=postgres \
	-e POSTGRES_DB=authdb \
	-p 5432:5432 -d postgres

docker-start:
	docker start authdb

docker-exec: docker-start
	docker exec -it authdb psql -U postgres authdb

docker-stop:
	docker stop authdb

# migrations
migrate-create:
	migrate create -ext sql -dir ./migrations -seq authdb

migrate-up:
	migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/authdb?sslmode=disable" up

migrate-down:
	migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/authdb?sslmode=disable" down

run:
	go run main.go