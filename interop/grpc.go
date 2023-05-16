package interop

import (
	propertyProsAuthApi "github.com/vireocloud/property-pros-sdk/api/auth/v1"
	apiCommon "github.com/vireocloud/property-pros-sdk/api/common/v1"
	propertyProsFinanceApi "github.com/vireocloud/property-pros-sdk/api/finance/v1"
	propertyProsAgreementApi "github.com/vireocloud/property-pros-sdk/api/note_purchase_agreement/v1"
	propertyProsStatementApi "github.com/vireocloud/property-pros-sdk/api/statement/v1"
	"google.golang.org/grpc"
)

// note purchase agreements
type UnsafeNotePurchaseAgreementServiceServer = propertyProsAgreementApi.UnsafeNotePurchaseAgreementServiceServer
type UnimplementedNotePurchaseAgreementServiceServer = propertyProsAgreementApi.UnimplementedNotePurchaseAgreementServiceServer
type NotePurchaseAgreementServiceServer = propertyProsAgreementApi.NotePurchaseAgreementServiceServer

type RecordRequestPayload = apiCommon.RecordRequestPayload
type RecordResultPayload = apiCommon.RecordResultPayload
type RecordColection = apiCommon.RecordCollection

type NotePurchaseAgreement = propertyProsAgreementApi.NotePurchaseAgreementRecord

type SaveNotePurchaseAgreementRequest = propertyProsAgreementApi.SaveNotePurchaseAgreementRequest
type SaveNotePurchaseAgreementResponse = propertyProsAgreementApi.SaveNotePurchaseAgreementResponse

type GetNotePurchaseAgreementRequest = propertyProsAgreementApi.GetNotePurchaseAgreementRequest
type GetNotePurchaseAgreementResponse = propertyProsAgreementApi.GetNotePurchaseAgreementResponse

type GetNotePurchaseAgreementsRequest = propertyProsAgreementApi.GetNotePurchaseAgreementsRequest
type GetNotePurchaseAgreementsResponse = propertyProsAgreementApi.GetNotePurchaseAgreementsResponse

type GetNotePurchaseAgreementDocRequest = propertyProsAgreementApi.GetNotePurchaseAgreementDocRequest
type GetNotePurchaseAgreementDocResponse = propertyProsAgreementApi.GetNotePurchaseAgreementDocResponse

type NotePurchaseAgreementServiceClient = propertyProsAgreementApi.NotePurchaseAgreementServiceClient

var NewNotePurchaseAgreementServiceClient = propertyProsAgreementApi.NewNotePurchaseAgreementServiceClient

var RegisterNotePurchaseAgreementServiceHandlerFromEndpoint = propertyProsAgreementApi.RegisterNotePurchaseAgreementServiceHandlerFromEndpoint
var RegisterNotePurchaseAgreementServiceServer = propertyProsAgreementApi.RegisterNotePurchaseAgreementServiceServer
var RegisterNotePurchaseAgreementServiceHandler = propertyProsAgreementApi.RegisterNotePurchaseAgreementServiceHandler

// auth
type AuthenticationServiceServer = propertyProsAuthApi.AuthenticationServiceServer
type UnsafeAuthenticationServiceServer = propertyProsAuthApi.UnsafeAuthenticationServiceServer
type UnimplementedAuthenticationServiceServer = propertyProsAuthApi.UnimplementedAuthenticationServiceServer

type User = propertyProsAuthApi.User
type AuthenticateUserRequest = propertyProsAuthApi.AuthenticateUserRequest
type AuthenticateUserResponse = propertyProsAuthApi.AuthenticateUserResponse

type AuthenticationServiceClient = propertyProsAuthApi.AuthenticationServiceClient

var NewAuthenticationServiceClient = propertyProsAuthApi.NewAuthenticationServiceClient

var RegisterAuthenticationServiceHandlerFromEndpoint = propertyProsAuthApi.RegisterAuthenticationServiceHandlerFromEndpoint
var RegisterAuthenticationServiceServer = propertyProsAuthApi.RegisterAuthenticationServiceServer

//statements

type StatementServiceServer = propertyProsStatementApi.StatementServiceServer
type UnsafeStatementServiceServer = propertyProsStatementApi.UnsafeStatementServiceServer
type UnimplementedStatementServiceServer = propertyProsStatementApi.UnimplementedStatementServiceServer

type Statement = propertyProsStatementApi.Statement
type GetStatementsRequest = propertyProsStatementApi.GetStatementsRequest
type GetStatementsResponse = propertyProsStatementApi.GetStatementsResponse

type StatementServiceClient = propertyProsStatementApi.StatementServiceClient

var NewStatementServiceClient = propertyProsStatementApi.NewStatementServiceClient

var RegisterStatementServiceHandlerFromEndpoint = propertyProsStatementApi.RegisterStatementServiceHandlerFromEndpoint

type CallOption = grpc.CallOption

// finance
type FinanceServiceServer = propertyProsFinanceApi.FinanceServiceServer
type UnsafeFinanceServiceServer = propertyProsFinanceApi.UnsafeFinanceServiceServer
type UnimplementedFinanceServiceServer = propertyProsFinanceApi.UnimplementedFinanceServiceServer

type FinancialData = propertyProsFinanceApi.FinancialData
type Account = propertyProsFinanceApi.Account
type Balance = propertyProsFinanceApi.Balance
type Transaction = propertyProsFinanceApi.Transaction
type Location = propertyProsFinanceApi.Location
type PaymentMeta = propertyProsFinanceApi.PaymentMeta

type SaveFinancialItemRequest = propertyProsFinanceApi.SaveFinancialItemRequest
type SaveFinancialItemResponse = propertyProsFinanceApi.SaveFinancialItemResponse

// type GetFinancialItemRequest = propertyProsFinanceApi.GetFinancialItemRequest
// type GetFinancialItemResponse = propertyProsFinanceApi.GetFinancialItemResponse

type GetFinancialItemsRequest = propertyProsFinanceApi.GetFinancialItemsRequest
type GetFinancialItemsResponse = propertyProsFinanceApi.GetFinancialItemsResponse
