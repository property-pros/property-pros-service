package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/omeid/uconfig/flat"
)

const (
	TagEnv  = "env"
	TagFlag = "flag"
	TagDesc = "desc"
)

// Validation tags described here: https://github.com/go-playground/validator
type Config struct {
	DocumentContentProviderSource string `env:"DOCUMENT_CONTENT_PROVIDER" flag:"document-content-provider" desc:"document provider grpc server address url" validate:"required"`
	ListenAddress                 string `env:"LISTEN_ADDRESS" flag:"listen-address" desc:"grpc server listening address host" validate:"required,ip4_addr"`
	ListenPort                    string `env:"LISTEN_PORT" flag:"listen-port" desc:"grpc server listening port" validate:"required,number"`
	GatewayPort                   string `env:"GATEWAY_PORT" flag:"gateway-port" desc:"grpc gateway server listening port" validate:"required,number"`
	AppEnv                        string `env:"APP_ENV" flag:"app-env" desc:"generic flag to describe the runtime environment dev/prod" validate:"required"`
	S3Endpoint                    string `env:"S3_ENDPOINT" flag:"s3-endpoint" desc:"s3 endpoint"`
	S3Region           			  string `env:"S3_REGION" flag:"s3-region" desc:"s3 region"`
	S3AccessKey        			  string `env:"S3_ACCESS_KEY" flag:"s3-access-key" desc:"s3 access key" validate:"required"`
	S3PrivateKey       			  string `env:"S3_SECRET_KEY" flag:"s3-private-key" desc:"s3 private key" validate:"required"`
	S3ForcePathStyle 			  bool 	 `env:"S3_FORCE_PATH_STYLE" flag:"s3-force-path-style" desc:"s3 config to force the path style"`
	S3BucketName       			  string `env:"S3_BUCKET_NAME" flag:"s3-bucket-name" desc:"s3 bucket name"`
	DbConnectionString 			  string `env:"DB_CONNECTION_STRING" flag:"db-connection-string" desc:"db connection string" validate:"required"`
	EnableTLS                     bool   `env:"ENABLE_TLS" flag:"enable-tls" desc:"enable tls"`
}

func NewConfig() (*Config, error) {
	instance := &Config{}
	return instance, LoadConfig(instance, &os.Args)
}

func LoadConfig(cfg interface{}, osArgs *[]string) error {
	fmt.Printf("path: %s\n", filepath.Join(os.Getenv("PROJECT_PATH"), "./.env"))
	err := godotenv.Load(filepath.Join(os.Getenv("PROJECT_PATH"), ".env"))

fmt.Printf("env vars: %v\n", os.Environ())
fmt.Printf("args: %v\n", os.Args)
	if err != nil {
		fmt.Printf("godotenv.Load %v\n", err)
		return err
	}
	
	// recursively iterates over each field of the nested struct
	fields, err := flat.View(cfg)
	if err != nil {
		fmt.Printf("flat.ViewError %v\n", err)
		return err
	}

	flagset := flag.NewFlagSet("", flag.ContinueOnError)

	for _, field := range fields {

		envName, ok := field.Tag(TagEnv)
		if !ok {
			continue
		}

		envValue := os.Getenv(envName)

		field.Set(envValue)

		flagName, ok := field.Tag(TagFlag)
		
		if !ok || flagName == "test.paniconexit0" {
			continue
		}

		flagDesc, _ := field.Tag(TagDesc)

		// writes flag value to variable
		flagset.Var(field, flagName, flagDesc)
	}

	var args []string
	if osArgs != nil {
		args = *osArgs
	} else {
		args = os.Args
	}

	args = removeElements(args, []string{"-test.paniconexit0", "-test.timeout", "-test.coverprofile"})

	fmt.Printf("parsing args: %+v\n", args)
	err = flagset.Parse(args[1:])
	if err != nil {
		fmt.Printf("flatset.Parse %v\n", err)
		return err
	}

	err = validator.New().Struct(cfg)
	if err != nil {
		return fmt.Errorf("config validation error: %w", err)
	}
	return nil
}

func removeElements(args []string, elements []string) []string {
	var result []string
	for _, arg := range args {
		if !contains(elements, arg) {
			result = append(result, arg)
		}
	}
	return result
}

func contains(elements []string, arg string) bool {
	for _, e := range elements {
		if strings.Contains(arg, e) {
			return true
		}
	}
	return false
}
