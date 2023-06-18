run:
	wire && go run wire_gen.go
dev-init:
	go mod tidy
	docker-compose down --remove-orphans
	docker-compose up --build
dev:
	docker-compose up
reset-db:
	docker-compose rm -f -s -v db
	docker-compose up -d db
lint:
	go vet ./...
	go fmt ./...
