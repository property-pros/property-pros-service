code:
	wire && go run wire_gen.go
container:
	go mod tidy
	docker-compose down --remove-orphans
	docker-compose up --build
db:
	docker run -e POSTGRES_PASSWORD=postgres -e POSTGRES_USER=postgres -e POSTGRES_DB=PropertyPros -p 5432:5432 postgres:14.1-alpine
dev:
	reflex -c ./reflex.conf
lint:
	go vet ./...
	go fmt ./...