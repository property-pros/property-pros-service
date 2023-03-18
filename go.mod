module github.com/vireocloud/property-pros-service

go 1.18

require (
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/Selvatico/go-mocket v1.0.7
	github.com/go-playground/validator/v10 v10.2.0
	github.com/google/uuid v1.1.2
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.15.0
	github.com/improbable-eng/grpc-web v0.15.0
	github.com/joho/godotenv v1.4.0
	github.com/omeid/uconfig v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.8.0
	github.com/vireocloud/property-pros-sdk v0.0.17
	golang.org/x/net v0.8.0
	google.golang.org/grpc v1.51.0
	gorm.io/driver/postgres v1.3.8
	gorm.io/gorm v1.23.8
)

require (
	github.com/aws/aws-sdk-go v1.44.214 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-playground/locales v0.13.0 // indirect
	github.com/go-playground/universal-translator v0.17.0 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.12.1 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20200714003250-2b9c44734f2b // indirect
	github.com/jackc/pgtype v1.11.0 // indirect
	github.com/jackc/pgx/v4 v4.16.1 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/kr/pretty v0.3.0 // indirect
	github.com/leodido/go-urn v1.2.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.8.0 // indirect
	github.com/stretchr/objx v0.4.0 // indirect
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

require (
	github.com/cenkalti/backoff/v4 v4.1.3 // indirect
	github.com/desertbit/timer v0.0.0-20180107155436-c41aec40b27f // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/wire v0.5.0
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/klauspost/compress v1.15.11 // indirect
	github.com/rs/cors v1.8.2 // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/text v0.8.0 // indirect
	google.golang.org/genproto v0.0.0-20221207170731-23e4bf6bdc37 // indirect
	google.golang.org/protobuf v1.28.1
	nhooyr.io/websocket v1.8.7 // indirect
)

replace github.com/omeid/uconfig => github.com/vireocloud/uconfig v0.4.0

// replace github.com/vireocloud/property-pros-sdk => ../property-pros-sdk
