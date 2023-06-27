package bootstrap

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"mime"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"gorm.io/gorm"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/improbable-eng/grpc-web/go/grpcweb"

	"github.com/vireocloud/property-pros-service/config"
	"github.com/vireocloud/property-pros-service/interop"
	controllers "github.com/vireocloud/property-pros-service/server/controllers"
	"github.com/vireocloud/property-pros-service/server/gateway"
	"github.com/vireocloud/property-pros-service/server/interceptors"
	"github.com/vireocloud/property-pros-service/server/third_party"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	enableTls = flag.Bool("enable_tls", false, "Use TLS - required for HTTP2.")
)

type App struct {
	Db                              *gorm.DB
	Config                          *config.Config
	AuthController                  *controllers.AuthController
	NotePurchaseAgreementController *controllers.NotePurchaseAgreementController
	StatementController             *controllers.StatementController

	apiInterceptor  *interceptors.AuthValidationInterceptor
	grpcInterceptor *interceptors.GrpcInterceptor
}

func NewApp(notePurchaseAgreementController *controllers.NotePurchaseAgreementController, authController *controllers.AuthController, statementController *controllers.StatementController, configuration *config.Config, grpcInterceptor *interceptors.GrpcInterceptor, authInterceptor *interceptors.AuthValidationInterceptor) *App {

	return &App{
		AuthController:                  authController,
		NotePurchaseAgreementController: notePurchaseAgreementController,
		StatementController:             statementController,
		Config:                          configuration,
		grpcInterceptor:                 grpcInterceptor,
		apiInterceptor:                  authInterceptor,
	}
}

func (a *App) Run() error {
	// controllers.chill()
	flag.Parse()
	// Adds gRPC internal logs. This is quite verbose, so adjust as desired!
	log := grpclog.NewLoggerV2(os.Stdout, io.Discard, io.Discard)
	grpclog.SetLoggerV2(log)

	if *enableTls {
		port := "9090"

		if *enableTls {
			port = "10000"
		}

		addr := fmt.Sprintf(":%v", port)

		grpcServer := grpc.NewServer(
			grpc.UnaryInterceptor(a.grpcInterceptor.HandleRequest),
			grpc.UnaryInterceptor(a.apiInterceptor.Validate),
		)

		wrappedServer := grpcweb.WrapServer(grpcServer)

		handler := func(resp http.ResponseWriter, req *http.Request) {
			wrappedServer.ServeHTTP(resp, req)
		}

		httpServer := http.Server{
			Addr:    addr,
			Handler: http.HandlerFunc(handler),
		}

		err := a.LogAvailableGrpcMethods(grpcServer)

		if err != nil {
			return err
		}

		go func() {
			grpclog.Info("Serving gRPC on https://", addr)
			if err := httpServer.ListenAndServeTLS("insecure/cert", "insecure/key"); err != nil {
				grpclog.Fatalf("failed starting http2 server: %v", err)
			}
		}()

		err = gateway.Run("dns:///"+addr, enableTls)

		if err != nil {
			log.Fatalln(err)
			return err
		}

		// a.registerControllers(grpcServer)
	} else {
		return a.StartInsecureServer()
	}

	return nil
}

func (*App) LogAvailableGrpcMethods(grpcServer *grpc.Server) error {
	b, err := json.MarshalIndent(grpcweb.ListGRPCResources(grpcServer), "", "  ")

	if err == nil {
		grpclog.Infof("Available Grpc Commands: %v", string(b))
	} else {
		return err
	}
	return err
}

func (a *App) registerControllers(grpcServer *grpc.Server, ctx context.Context, gwmux *runtime.ServeMux, dialUrl string, dopts []grpc.DialOption) error {
	interop.RegisterNotePurchaseAgreementServiceServer(grpcServer, a.NotePurchaseAgreementController)
	interop.RegisterAuthenticationServiceServer(grpcServer, a.AuthController)
	interop.RegisterStatementServiceServer(grpcServer, a.StatementController)

	err := interop.RegisterNotePurchaseAgreementServiceHandlerFromEndpoint(ctx, gwmux, dialUrl, dopts)

	if err != nil {
		grpclog.Fatalf("failed starting http server: %v", err)
		return err
	}

	err = interop.RegisterAuthenticationServiceHandlerFromEndpoint(ctx, gwmux, dialUrl, dopts)

	if err != nil {
		grpclog.Fatalf("failed starting http server: %v", err)
		return err
	}

	err = interop.RegisterStatementServiceHandlerFromEndpoint(ctx, gwmux, dialUrl, dopts)
	if err != nil {
		grpclog.Fatalf("failed starting http server: %v", err)
		return err
	}

	return nil
}

func grpcHandlerFunc(grpcServer *grpc.Server, grpcWebServer *grpcweb.WrappedGrpcServer, restHandler http.Handler, oa http.Handler) http.Handler {

	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		grpclog.Infof("url: %v\r\n\r\ncontent type: %v\r\n\r\n", r.URL.Path, r.Header.Get("Content-Type"))

		if strings.Contains(r.Header.Get("Content-Type"), "application/grpc-web+proto") {
			grpclog.Infoln("grpc-web request")
			grpcWebServer.ServeHTTP(w, r)
			return
		}

		if strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpclog.Infoln("grpc request")

			grpcServer.ServeHTTP(w, r)
			return
		}

		if matched, err := regexp.MatchString("v\\d", r.URL.Path); err == nil && matched {
			grpclog.Infoln("rest api request")
			restHandler.ServeHTTP(w, r)
			return
		}

		oa.ServeHTTP(w, r)

	}), &http2.Server{})
}

func (a *App) StartInsecureServer() error {
	wg := sync.WaitGroup{}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(a.grpcInterceptor.HandleRequest))

	wrappedServer := grpcweb.WrapServer(grpcServer)

	gwmux := runtime.NewServeMux()

	parentContext := context.Background()

	ctx, cancel := context.WithCancel(parentContext)

	dopts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	port := a.Config.ListenPort
	scheme := "dns:///"

	serverUrl := fmt.Sprintf("%v:%v", a.Config.ListenAddress, port)
	dialUrl := fmt.Sprintf("%vlocalhost:%v", scheme, port)

	fmt.Println("server url: ", serverUrl)
	fmt.Println("dial url: ", dialUrl)
	var err error
	go (func(grpcServer *grpc.Server, ctx context.Context, gwmux *runtime.ServeMux, dialUrl string, dopts []grpc.DialOption, returnErr *error) {

		err = a.registerControllers(grpcServer, ctx, gwmux, dialUrl, dopts)

		if err != nil {
			fmt.Println("registerControllers failed: ", err)

			if returnErr != nil {
				returnErr = &err
			}

			cancel()
		}

		err = a.LogAvailableGrpcMethods(grpcServer)

		if err != nil {
			returnErr = &err
			cancel()
		}
	})(grpcServer, ctx, gwmux, dialUrl, dopts, &err)

	wg.Add(1)

	go func() {
		fmt.Printf("Listening at %v", a.Config.ListenAddress)
		if err := http.ListenAndServe(serverUrl, grpcHandlerFunc(grpcServer, wrappedServer, gwmux, getOpenAPIHandler())); err != nil {
			fmt.Println("Http listener failed: ", err)
			wg.Done()
		}
	}()

	wg.Wait()

	return err
}

func getOpenAPIHandler() http.Handler {
	mime.AddExtensionType(".svg", "image/svg+xml")
	// Use subdirectory in embedded files
	subFS, err := fs.Sub(third_party.OpenAPI, "OpenAPI")
	if err != nil {
		panic("couldn't create sub filesystem: " + err.Error())
	}
	return http.FileServer(http.FS(subFS))
}

func NewGrpcConnection(config *config.Config) grpc.ClientConnInterface {
	// connection, err := grpc.Dial(config.DocumentContentProviderSource, grpc.WithTransportCredentials(insecure.NewCredentials()))
	connection, err := grpc.Dial("localhost:8020", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(fmt.Errorf("NewGrpcConnection failed: %w", err))
	}

	return connection
}

func NewNotePurchaseAgreementClient(conn grpc.ClientConnInterface) interop.NotePurchaseAgreementServiceClient {

	return interop.NewNotePurchaseAgreementServiceClient(conn)
}
