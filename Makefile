run:
	wire && go run wire_gen.go
watch:
	reflex -c reflex.conf
dev-init:
	go mod tidy
	docker-compose down --remove-orphans
	docker-compose up --build db s3mock & make watch
dev:
	docker-compose up db s3mock & make watch
reset-db:
	docker-compose rm -f -s -v db
	docker-compose up -d db
lint:
	go vet ./...
	go fmt ./...
test-api:
	rm wire_gen.go || true
	go test ./... -tags=integration -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

