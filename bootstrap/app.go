package bootstrap

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
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

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/improbable-eng/grpc-web/go/grpcweb"

	"github.com/vireocloud/property-pros-service/server/gateway"

	propertyProsApi "github.com/vireocloud/property-pros-docs/generated/notePurchaseAgreement"
	"github.com/vireocloud/property-pros-service/config"
	controllers "github.com/vireocloud/property-pros-service/server/controllers"
	"github.com/vireocloud/property-pros-service/server/third_party"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	enableTls = flag.Bool("enable_tls", false, "Use TLS - required for HTTP2.")
)

type App struct {
	config                          *config.Config
	notePurchaseAgreementController *controllers.NotePurchaseAgreementController
}

func NewApp(notePurchaseAgreementController *controllers.NotePurchaseAgreementController, configuration *config.Config) *App {
	return &App{
		notePurchaseAgreementController: notePurchaseAgreementController,
		config:                          configuration,
	}
}

func (a *App) Run() error {
	// controllers.chill()
	flag.Parse()
	// Adds gRPC internal logs. This is quite verbose, so adjust as desired!
	log := grpclog.NewLoggerV2(os.Stdout, ioutil.Discard, ioutil.Discard)
	grpclog.SetLoggerV2(log)

	if *enableTls {
		port := "9090"

		if *enableTls {
			port = "10000"
		}

		addr := fmt.Sprintf(":%v", port)

		grpcServer := grpc.NewServer()

		a.registerControllers(grpcServer)

		wrappedServer := grpcweb.WrapServer(grpcServer)

		handler := func(resp http.ResponseWriter, req *http.Request) {
			wrappedServer.ServeHTTP(resp, req)
		}

		httpServer := http.Server{
			Addr:    addr,
			Handler: http.HandlerFunc(handler),
		}

		b, err := json.MarshalIndent(grpcweb.ListGRPCResources(grpcServer), "", "  ")
		if err == nil {
			grpclog.Infof("Available Grpc Commands: %v", string(b))
		} else {
			return err
		}

		if *enableTls {
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
		} else {
			go func() {
				grpclog.Info("Serving gRPC on http://", addr)
				if err := httpServer.ListenAndServe(); err != nil {
					grpclog.Fatalf("failed starting http server: %v", err)
				}
			}()

			err = gateway.Run("dns:///"+addr, enableTls)

			if err != nil {
				log.Fatalln(err)
				return err
			}
		}
	} else {
		return a.StartInsecureServer()
	}

	return nil
}

func (a *App) registerControllers(grpcServer *grpc.Server) {

	propertyProsApi.RegisterNotePurchaseAgreementServiceServer(grpcServer, a.notePurchaseAgreementController)

}

func grpcHandlerFunc(grpcServer *grpc.Server, grpcWebServer *grpcweb.WrappedGrpcServer, restHandler http.Handler, oa http.Handler) http.Handler {

	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("handling request")
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

	grpcServer := grpc.NewServer()

	wrappedServer := grpcweb.WrapServer(grpcServer)

	a.registerControllers(grpcServer)

	gwmux := runtime.NewServeMux()

	ctx := context.Background()

	dopts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	port := a.config.ListenPort
	scheme := "dns:///"

	serverUrl := fmt.Sprintf("%v:%v", a.config.ListenAddress, port)
	dialUrl := fmt.Sprintf("%vlocalhost:%v", scheme, port)

	fmt.Println("server url: ", serverUrl)
	fmt.Println("dial url: ", dialUrl)

	wg.Add(1)

	go func() {
		fmt.Printf("Listening at %v", a.config.ListenAddress)
		if err := http.ListenAndServe(serverUrl, grpcHandlerFunc(grpcServer, wrappedServer, gwmux, getOpenAPIHandler())); err != nil {
			fmt.Println("Http listener failed: ", err)
			wg.Done()
		}
	}()

	err := propertyProsApi.RegisterNotePurchaseAgreementServiceHandlerFromEndpoint(ctx, gwmux, dialUrl, dopts)

	if err != nil {
		grpclog.Fatalf("failed starting http server: %v", err)
		return err
	}

	wg.Wait()

	return nil
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
	connection, err := grpc.Dial(config.DocumentContentProviderSource, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		panic(fmt.Errorf("NewGrpcConnection failed: %w", err))
	}

	return connection
}

func NewNotePurchaseAgreementClient(conn grpc.ClientConnInterface) propertyProsApi.NotePurchaseAgreementServiceClient {

	return propertyProsApi.NewNotePurchaseAgreementServiceClient(conn)

}
