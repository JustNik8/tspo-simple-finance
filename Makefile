migrateup:
	 migrate -path migrations -database "postgresql://user:user@localhost:5440/finance_db?sslmode=disable" --verbose up

migrateup1:
	 migrate -path migrations -database "postgresql://user:user@localhost:5440/finance_db?sslmode=disable" --verbose up 1

migratedown:
	 migrate -path migrations -database "postgresql://user:user@localhost:5440/finance_db?sslmode=disable" --verbose down

migratedown1:
	 migrate -path migrations -database "postgresql://user:user@localhost:5440/finance_db?sslmode=disable" --verbose down 1