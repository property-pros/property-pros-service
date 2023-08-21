run:
	wire
	go run wire_gen.go
build: 
	wire && go build wire_gen.go
watch:
	reflex -c reflex.conf -v
dev-init:
	go mod tidy
	docker-compose down --remove-orphans
	docker-compose up --build db s3mock & make dev
dev:
	docker-compose up db s3mock & make run
reset-db:
	rm -rf ./db && \
	docker-compose up -d db
lint:
	go vet ./...
	go fmt ./...
pre-test-api:
	make kill-app && \
	while ! curl --output /dev/null --silent --head --fail http://localhost:8000; do sleep .2; done && \
	echo "listening on port 8000"
test-api:
	make pre-test-api && \
	echo "starting tests..." && \
	export PROJECT_PATH=$(shell pwd) TEST_SERVICE="localhost:8000" && \
	(go test $(shell find . -name '*api_test.go') -coverprofile=coverage.out | tee test-output.txt) && \
	go tool cover -html=coverage.out & \
	make post-test-api
post-test-api:
	make reset-db && \
	docker-compose up -d s3mock & \
	while ! (export PGPASSWORD=postgres; echo '\q' | psql -h localhost -p 5432 -U postgres -d PropertyPros > /dev/null 2>&1); do sleep .2; done && \
	make run
kill-app:
	kill $(shell lsof -t -i:8000) || true
install-pg_isready:
	wget https://github.com/pganalyze/pg_isready/releases/download/v2.4.0/pg_isready-2.4.0.tar.gz
	tar -xzf pg_isready-2.4.0.tar.gz 
	cd pg_isready-2.4.0 && ./configure && make && sudo make install
	cd .. && \
	rm -rf pg_isready-2.4.0 pg_isready-2.4.0.tar.gz
	echo "pg_isready installed!"