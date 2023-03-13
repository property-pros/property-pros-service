package documents

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/vireocloud/property-pros-service/interfaces"
	"github.com/vireocloud/property-pros-service/interop"
)

var testDocumentContentManager interfaces.IDocumentContentService
var mockClient interop.NotePurchaseAgreementServiceClient
var mockDocUploader *DocUploaderMock

func TestBuildNotePurchaseAgreement(t *testing.T) {
	Setup()
	mockDocUploader.On("PutObject").Return("someId", nil)
	documentContent, err := testDocumentContentManager.BuildNotePurchaseAgreement(context.TODO(), &interop.NotePurchaseAgreement{})

	if err != nil {
		t.Error(err)
	}

	if documentContent == nil {
		t.Errorf("expected BuildNotePurchaseAgreement to not return nil")
	}

	if documentContent.GetDocContent() == nil {
		t.Errorf("expected BuildNotePurchaseAgreement return document with content")
	}

	if len(documentContent.GetDocContent()) != 3 {

		t.Errorf("expected BuildNotePurchaseAgreement return document with 3 bytes of content")
	}

	Teardown()
}

type DocUploaderMock struct {
	mock.Mock
}

func (m *DocUploaderMock) PutObject(ctx context.Context, content []byte) (string, error) {
	args := m.Called(ctx, content)

	return args.Get(0).(string), nil
}

func (m *DocUploaderMock) GetObject(ctx context.Context, url string) ([]byte, error) {
	args := m.Called(ctx, url)

	return args.Get(0).([]byte), nil
}
func TestBuildAccountStatement(t *testing.T) {}

func Setup() {

	mockClient = &ClientMock{}
	mockDocUploader = &DocUploaderMock{}

	testDocumentContentManager = NewDocumentContentService(mockClient, mockDocUploader)
}

func Teardown() {
	testDocumentContentManager = nil
	mockClient = nil
}

func setup() {}

func teardown() {}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}
