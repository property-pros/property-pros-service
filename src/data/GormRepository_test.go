package data

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"regexp"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	gormMock "github.com/Selvatico/go-mocket"
	"github.com/jinzhu/gorm"
	i "github.com/vireocloud/property-pros-service/src/interfaces"
)

func SaveExpectations(model *FakeModel) {
	sqlDBMocker.ExpectQuery(`^SELECT \* FROM \"fake_models\" WHERE \(ID = \$1\)$`).WithArgs(model.ID).WillReturnRows(sqlmock.NewRows([]string{"ID", "Name", "Value"}))

	sqlDBMocker.ExpectBegin()

	sqlDBMocker.ExpectQuery(`^INSERT INTO \"fake_models\" \(\"name\"\,\"value\"\) VALUES \(\$1\,\$2\) RETURNING \"fake_models\"\.\"id\"$`).WithArgs(model.Name, model.Value).WillReturnRows(sqlmock.NewRows([]string{"ID"}).AddRow(expectedResult.ID))

	sqlDBMocker.ExpectCommit()
}

func TestShouldSaveWithInsert(t *testing.T) {

	model := &FakeModel{ID: 0, Name: "test name", Value: "test value"}

	SaveExpectations(model)

	sqlDBMocker.MatchExpectationsInOrder(true)

	payload, query := repo.GetQueryAndPayload(model)

	if query == nil {
		t.Error("expected query to not be nil;  recieved nil")
	}

	if payload == nil {
		t.Error("expected payload to not be nil;  recieved nil")
	}

	payloadFakeModel := payload.(*FakeModel)
	queryFakeModel := query.(*FakeModel)

	if payloadFakeModel.ID != uint32(0) {
		t.Errorf("expected payload.ID to be 0; recieved %v", payloadFakeModel.ID)
	}

	if queryFakeModel.ID != uint32(0) {
		t.Errorf("expected query.ID to be 0; recieved %v", queryFakeModel.ID)
	}

	result, saveErr := repo.Save(payload, query)

	if saveErr != nil {
		t.Error(saveErr)
	}

	err := sqlDBMocker.ExpectationsWereMet()

	if err != nil {
		t.Error(err)
	}

	concreteResult, ok := result.(*FakeModel)

	if !ok {
		t.Error("Something went wrong while running repo.Save.  Conversion to *FakeModel failed.")
	}

	if concreteResult == nil {
		t.Error("Expected 1 result from repo.Save.  Recieved nil")
	}
}

func TestShouldQuery(t *testing.T) {

	sqlDBMocker.ExpectQuery(`^SELECT \* FROM \"fake_models\" WHERE \(\"fake_models\"\.\"name\" = \$1\)$`).WithArgs("Test Name").WillReturnRows(sqlRows)

	results := repo.Query(&FakeModel{Name: "Test Name"})

	err := sqlDBMocker.ExpectationsWereMet()

	if err != nil {
		t.Error(err)
	}

	if len(results) != 1 {
		t.Errorf("Expected query result length to equal 1.  Recieved %v", len(results))
	}
}

func TestShouldFindOne(t *testing.T) {

	query := &FakeModel{}

	query.ID = expectedResult.ID

	sqlDBMocker.ExpectQuery(`^SELECT \* FROM \"fake_models\" WHERE \(ID \= \$1\)$`).WithArgs(query.ID).WillReturnRows(sqlRows)

	repo.FindOne(query)

	err := sqlDBMocker.ExpectationsWereMet()

	if err != nil {
		t.Error(err)
	}
}

func TestShouldDelete(t *testing.T) {

	log.Printf("TestShouldDelete testEntity: %v", expectedResult)

	sqlDBMocker.ExpectBegin()

	sqlDBMocker.ExpectExec(regexp.QuoteMeta("DELETE FROM \"fake_models\" WHERE \"fake_models\".\"id\" = $1")).
		WithArgs(uint64(expectedResult.ID)).
		WillReturnResult(sqlmock.NewResult(int64(expectedResult.ID), 1))

	sqlDBMocker.ExpectCommit()

	sqlDBMocker.MatchExpectationsInOrder(true)

	deleteResult, _ := repo.Delete(expectedResult)

	if deleteResult == nil {
		t.Error("Expected delete result to be not nil;  recieved nil")
	}

	err := sqlDBMocker.ExpectationsWereMet()

	if err != nil {
		t.Error(err)
	}
}

var repo i.IRepository

var sqlRows *sqlmock.Rows

var slqmockDbInstance *sql.DB

var sqlDBMocker sqlmock.Sqlmock

var gormDBMock *gormMock.MockCatcher

var commonMockReply *[]map[string]interface{}

var expectedResult *FakeModel

func TestMain(m *testing.M) {
	setup(m)
	code := m.Run()
	tearDown(m)
	os.Exit(code)
}

func setup(m *testing.M) {

	expectedResult = &FakeModel{ID: 1234, Name: "test name", Value: "test value"}

	slqmockDbInstance, sqlDBMocker, _ = sqlmock.New()

	sqlRows = sqlmock.NewRows([]string{"ID", "Name", "Value"}).AddRow(expectedResult.ID, expectedResult.Name, expectedResult.Value)

	// GORM
	db, err := gorm.Open("postgres", slqmockDbInstance) // Could be any connection string

	if err != nil {
		fmt.Println(err)
	}

	repo = &GormRepository{
		db: db,
	}

	gormDBMock = gormMock.Catcher

	commonMockReply = &[]map[string]interface{}{{"ID": expectedResult.ID, "Name": expectedResult.Name, "Value": expectedResult.Value}}

	if err != nil {
		panic(err)
	}
	// db.LogMode(true)
}

func tearDown(m *testing.M) {

}

type FakeModel struct {
	ID    uint32
	Name  string
	Value string
}
