# Run your Go server with env vars
run:
	ENV=development go run main.go

# Run database migrations on app start
migrate:
	go run main.go

# Build binary
build:
	go build -o kingdup

# Clean up binary
clean:
	rm -f kingdup

# Print available commands
help:
	@echo "make run           - Start the Go server (development mode)"
	@echo "make migrate       - Run the app with auto migration (built-in)"
	@echo "make build         - Build the app"
	@echo "make clean         - Remove built binary"
