run:
	wire && go run wire_gen.go
dev-init:
	docker-compose down && docker-compose up --build
dev:
	docker-compose up