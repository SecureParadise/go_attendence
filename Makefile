DATABASE_URL=postgres://postgres:admin@localhost:5432/attendence?sslmode=disable
createdb:
	createdb -U postgres attendence
dropdb:
	dropdb -U postgres attendence
migrategen:
	migrate create -ext sql -dir internal/db/migrations -seq subject
migrateup:
	migrate -path internal/db/migrations -database "$(DATABASE_URL)" -verbose up
migrateup1:
	migrate -path internal/db/migrations -database "$(DATABASE_URL)" -verbose up 1
migratedown:
	migrate -path internal/db/migrations -database "$(DATABASE_URL)" -verbose down
migratedown1:
	migrate -path internal/db/migrations -database "$(DATABASE_URL)" -verbose down 1

migratestat:
	migrate -path internal/db/migrations -database "$(DATABASE_URL)" version
migrateclsforce:
	migrate -path internal/db/migrations -database "$(DATABASE_URL)" force 1
sqlc:
	sqlc generate
.PHONY: createdb dropdb migrategen migrateup migratedown sqlc