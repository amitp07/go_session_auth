### Create Migration file 
migrate create --ext sql -dir internal/database/migrations -seq <file_name> 
### Apply Migration 
migrate -database $GO_SESSION_AUTH_DSN -source file://<path_to_file> up