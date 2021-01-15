postgres:
	docker run --name sense-it -p 5433:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

createdb:
	docker exec -it sense-it createdb --username=root --owner=root sense_it

dropdb:
	docker exec -it sense-it dropdb sense_it

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5433/sense_it?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5433/sense_it?sslmode=disable" -verbose down