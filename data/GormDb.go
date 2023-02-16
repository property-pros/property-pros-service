package data

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewGormDatabase() (*gorm.DB, error) {
	postgresConfig := postgres.Open("host=db port=5432 user=postgres dbname=PropertyPros password=postgres")

	db, err := gorm.Open(postgresConfig, &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return db, nil
}

func CreateFixtures(string) {

}
