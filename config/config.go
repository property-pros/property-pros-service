package config

import (
	"flag"
	"fmt"
	"os"

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
	DocumentContentProviderSource string `env:"DOCUMENT_CONTENT_PROVIDER" flag:"document-content-provider" desc:"document provider grpc server address host:port" validate:"required,hostname_port"`
	ListenAddress                 string `env:"LISTEN_ADDRESS" flag:"listen-address" desc:"grpc server listening address host" validate:"required,ip4_addr"`
	ListenPort                    string `env:"LISTEN_PORT" flag:"listen-port" desc:"grpc server listening address host:port" validate:"required,number"`
	AppEnv                        string `env:"APP_ENV" flag:"app-env" desc:"generic flag to describe the runtime environment dev/prod" validate:"required"`
	S3Endpoint                    string `env:"S3_ENDPOINT" flag:"s3-endpoint" desc:"s3 endpoint" validate:"required,hostname_port"`
}

func NewConfig() (*Config, error) {
	instance := &Config{}

	return instance, LoadConfig(instance, &os.Args)
}

func LoadConfig(cfg interface{}, osArgs *[]string) error {
	godotenv.Load(".env")

	// recursively iterates over each field of the nested struct
	fields, err := flat.View(cfg)
	if err != nil {
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
		if !ok {
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

	err = flagset.Parse(args[1:])
	if err != nil {
		return err
	}

	err = validator.New().Struct(cfg)
	if err != nil {
		return fmt.Errorf("config validation error: %w", err)
	}
	return nil
}
