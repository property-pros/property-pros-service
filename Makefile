run:
	wire && go run wire_gen.go
dev-init:
	go mod tidy
	docker-compose down --remove-orphans
	docker-compose up --build
dev:
	docker-compose up
lint:
	go vet ./...
	go fmt ./...