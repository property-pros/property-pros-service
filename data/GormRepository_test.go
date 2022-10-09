package data

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"log"
	"os"
	"regexp"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	gormMock "github.com/Selvatico/go-mocket"
	i "github.com/vireocloud/property-pros-service/interfaces"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InsertExpectations(model *FakeModel) {
	sqlDBMocker.ExpectBegin()
	// sqlDBMocker.ExpectQuery(`^SELECT \* FROM \"fake_models\" WHERE \(ID = \$1\)$`).WithArgs(model.ID).WillReturnRows(sqlmock.NewRows([]string{"ID", "Name", "Value"}))

	queryExpectation := sqlDBMocker.ExpectExec(`^INSERT INTO \"fake_models\" \(\"id\",\"name\"\,\"value\"\) VALUES \(\$1\,\$2\,\$3\)$`)
	// RETURNING \"fake_models\"\.\"id\"
	queryExpectation.WithArgs(model.GetId(), model.Name, model.Value).WillReturnResult(driver.ResultNoRows) //.NewRows([]string{"id"}).AddRow(expectedResult.ID)

	sqlDBMocker.ExpectCommit()
}

func TestShouldSaveWithInsert(t *testing.T) {

	payload := &FakeModel{Name: "test name", Value: "test value"}

	InsertExpectations(payload)

	sqlDBMocker.MatchExpectationsInOrder(true)

	concreteResult, saveErr := repo.Save(payload)

	if saveErr != nil {
		t.Error(saveErr)
	}

	err := sqlDBMocker.ExpectationsWereMet()

	if err != nil {
		t.Error(err)
	}

	if concreteResult == nil {
		t.Error("Expected 1 result from repo.Save.  Recieved nil")
	}
}

func UpdateExpectations(model *FakeModel) {
	sqlDBMocker.ExpectBegin()

	queryExpectation := sqlDBMocker.ExpectExec(regexp.QuoteMeta(`UPDATE "fake_models" SET "name"=$1,"value"=$2 WHERE "id" = $3`))

	queryExpectation.WithArgs(model.Name, model.Value, model.GetId()).WillReturnResult(driver.RowsAffected(1))

	sqlDBMocker.ExpectCommit()
}

func TestShouldSaveWithUpdate(t *testing.T) {

	payload := &FakeModel{ID: "test_id", Name: "test name", Value: "test value"}
	t.Logf("payload: %+#v\n\n", payload)
	t.Logf("payload.GetId(): %v\n\n", payload.GetId())
	UpdateExpectations(payload)

	sqlDBMocker.MatchExpectationsInOrder(true)

	concreteResult, saveErr := repo.Save(payload)

	if saveErr != nil {
		t.Error(saveErr)
	}

	err := sqlDBMocker.ExpectationsWereMet()

	if err != nil {
		t.Error(err)
	}

	if concreteResult == nil {
		t.Error("Expected 1 result from repo.Save.  Recieved nil")
	}
}

func TestShouldQuery(t *testing.T) {

	sqlDBMocker.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "fake_models" WHERE "fake_models"."name" = $1`)).WithArgs("Test Name").WillReturnRows(sqlRows)

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

	sqlDBMocker.ExpectQuery(`^SELECT \* FROM \"fake_models\" WHERE \"fake_models\".\"id\" \= \$1$`).WithArgs(query.ID).WillReturnRows(sqlRows)

	repo.FindOne(query)

	err := sqlDBMocker.ExpectationsWereMet()

	if err != nil {
		t.Error(err)
	}
}

func TestShouldDelete(t *testing.T) {

	log.Printf("TestShouldDelete testEntity: %v", expectedResult)

	sqlDBMocker.ExpectBegin()

	sqlDBMocker.ExpectExec(regexp.QuoteMeta("DELETE FROM \"fake_models\" WHERE \"fake_models\".\"id\" = $1")).WithArgs(string(expectedResult.ID)).WillReturnResult(driver.ResultNoRows)

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

var repo i.IRepository[FakeModel]

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

	expectedResult = &FakeModel{ID: "1234", Name: "test name", Value: "test value"}

	slqmockDbInstance, sqlDBMocker, _ = sqlmock.New()

	sqlRows = sqlmock.NewRows([]string{"ID", "Name", "Value"}).AddRow(expectedResult.ID, expectedResult.Name, expectedResult.Value)

	// GORM
	// slqmockDbInstance.
	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: slqmockDbInstance,
	})) // Could be any connection string

	if err != nil {
		fmt.Println(err)
	}

	repo = &GormRepository[FakeModel, *FakeModel]{
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
	ID    string
	Name  string
	Value string
}

func (model *FakeModel) Getcontext() {}
func (model *FakeModel) GetId() string {
	return model.ID
}
func (model *FakeModel) GetPayload() {}
