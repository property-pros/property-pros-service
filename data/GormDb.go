package data

import (
	"fmt"
	"os"
	"time"

	"github.com/vireocloud/property-pros-service/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewGormDatabase(config *config.Config ) (*gorm.DB, error) {
	// for docker compose
	// postgresConfig := postgres.Open("host=db port=5432 user=postgres dbname=PropertyPros password=postgres")

	// for local
	postgresConfig := postgres.Open(config.DbConnectionString)
	db, err := gorm.Open(postgresConfig, &gorm.Config{})

	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&User{})
	db.AutoMigrate(&NotePurchaseAgreement{})
	db.AutoMigrate(&Statement{})
	
	userFixture(db)  // Create the user fixture
	statementFixture(db)  // Create the statement fixture with the created user
	notePurchaseAgreementFixture(db)
	return db, nil
}

func notePurchaseAgreementFixture(db *gorm.DB) {

  // Get the test user
  var user User
  db.Where("email_address = ?", "srt0422@yahoo.com").First(&user)

  // Create a test note purchase agreement
  npa := NotePurchaseAgreement{
	Id: "test-npa-1", 
	UserId: user.Id,
	CreatedOn: time.Now(),
	FundsCommitted: 0,
  }

  // Save the test NPA
  db.FirstOrCreate(&npa)
}

func userFixture(db *gorm.DB) {
	user := User{
		Id:           os.Getenv("TEST_ACCOUNT_ID"),
		EmailAddress: "srt0422@yahoo.com",
		Password:     "password",  // Consider using a secure way to store passwords
	}

	db.FirstOrCreate(&user, User{Id: user.Id})
}

func statementFixture(db *gorm.DB) {

	email := "srt0422@yahoo.com"
	userID := os.Getenv("TEST_ACCOUNT_ID")

	principle := 100000.0
	rate := 0.04 / 12 // monthly interest rate
	balance := principle
	totalIncome := 0.0

	for i := 0; i < 3; i++ {
		startDate := time.Now().AddDate(0, -i-1, 0) // start of previous month
		endDate := time.Now().AddDate(0, -i, 0)     // end of previous month

		startDate = time.Date(startDate.Year(), startDate.Month(), 1, 0, 0, 0, 0, time.UTC) // first day of previous month
		endDate = time.Date(endDate.Year(), endDate.Month()+1, 0, 0, 0, 0, 0, time.UTC)     // last day of previous month

		income := balance * rate // monthly interest income
		balance += income        // new balance after interest
		totalIncome += income    // cumulative income

		statement := Statement{
			Id:              fmt.Sprintf("%s-%d", userID, i+1),
			UserId:          userID,
			EmailAddress:    email,
			StartPeriodDate: startDate,
			EndPeriodDate:   endDate,
			Balance:         balance,
			TotalIncome:     totalIncome,
			Principle:       principle,
		}

		db.FirstOrCreate(&statement, Statement{Id: statement.Id})
	}
}