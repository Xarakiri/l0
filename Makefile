postgres:
	docker run --name psqld -p 6432:5432 -e POSTGRES_PASSWORD=admin -e POSTGRES_USER=root -e POSTGRES_DB=root -d postgres

nats:
	docker run --name natsq -e NATS_ADDRESS=nats:4222 -p 4222:4222 -d nats-streaming:0.9.2

createdb:
	docker exec -it psqld createdb --username=root --owner=root level0

dropdb:
	docker exec -it psqld dropdb level0

migrateup:
	migrate -path db/migration -database "postgresql://root:admin@localhost:6432/level0?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:admin@localhost:6432/level0?sslmode=disable" -verbose down

.PHONY: postgres createdb dropdb