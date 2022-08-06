package documents

import (
	"context"
	"os"
	"testing"

	"github.com/vireocloud/property-pros-service/interfaces"
	"github.com/vireocloud/property-pros-service/interop"
)

var testDocumentContentManager interfaces.IDocumentContentService
var mockClient interop.NotePurchaseAgreementServiceClient

func TestBuildNotePurchaseAgreement(t *testing.T) {
	Setup()

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

type ClientMock struct{}

func (*ClientMock) GetNotePurchaseAgreementDoc(ctx context.Context, in *interop.GetNotePurchaseAgreementDocRequest, opts ...interop.CallOption) (*interop.GetNotePurchaseAgreementDocResponse, error) {
	return &interop.GetNotePurchaseAgreementDocResponse{
		FileContent: []byte{'1', '2', '3'},
	}, nil
}

func TestBuildAccountStatement(t *testing.T) {}

func Setup() {

	mockClient = &ClientMock{}

	testDocumentContentManager = NewDocumentContentManager(mockClient)
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
