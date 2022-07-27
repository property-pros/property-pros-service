package main

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
	"github.com/jinzhu/gorm"
	"github.com/vireocloud/property-pros-service/server/gateway"

	propertyProsApi "github.com/vireocloud/property-pros-service/generated/notePurchaseAgreement"
	controllers "github.com/vireocloud/property-pros-service/server/controllers"
	"github.com/vireocloud/property-pros-service/server/third_party"
)

var (
	enableTls = flag.Bool("enable_tls", false, "Use TLS - required for HTTP2.")
)

func main() {
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

		propertyProsApi.RegisterNotePurchaseAgreementServiceServer(grpcServer, &controllers.NotePurchaseAgreementController{})

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
		}

		if *enableTls {
			go func() {
				grpclog.Info("Serving gRPC on https://", addr)
				if err := httpServer.ListenAndServeTLS("insecure/cert", "insecure/key"); err != nil {
					grpclog.Fatalf("failed starting http2 server: %v", err)
				}
			}()

			err = gateway.Run("dns:///"+addr, enableTls)

			log.Fatalln(err)
		} else {
			go func() {
				grpclog.Info("Serving gRPC on http://", addr)
				if err := httpServer.ListenAndServe(); err != nil {
					grpclog.Fatalf("failed starting http server: %v", err)
				}
			}()

			err = gateway.Run("dns:///"+addr, enableTls)

			log.Fatalln(err)
		}
	} else {
		StartInsecureServer()
	}
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

func StartInsecureServer() {
	wg := sync.WaitGroup{}

	grpcServer := grpc.NewServer()

	wrappedServer := grpcweb.WrapServer(grpcServer)

	propertyProsApi.RegisterNotePurchaseAgreementServiceServer(grpcServer, &controllers.NotePurchaseAgreementController{})

	gwmux := runtime.NewServeMux()

	ctx := context.Background()

	dopts := []grpc.DialOption{grpc.WithInsecure()}

	host := "0.0.0.0"
	port := "8020"
	scheme := "dns:///"

	serverUrl := fmt.Sprintf("%v:%v", host, port)
	dialUrl := fmt.Sprintf("%vlocalhost:%v", scheme, port)

	fmt.Println("server url: ", serverUrl)
	fmt.Println("dial url: ", dialUrl)

	wg.Add(1)

	go func() {
		fmt.Println("Listening on 8020")
		if err := http.ListenAndServe(serverUrl, grpcHandlerFunc(grpcServer, wrappedServer, gwmux, getOpenAPIHandler())); err != nil {
			fmt.Println("Http listener failed: ", err)
			wg.Done()
		}
	}()

	err := propertyProsApi.RegisterNotePurchaseAgreementServiceHandlerFromEndpoint(ctx, gwmux, dialUrl, dopts)

	if err != nil {
		grpclog.Fatalf("failed starting http server: %v", err)
	}

	wg.Wait()
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

func bootstrapControllers() (*controllers.GroupsController, *controllers.ActivitiesController, *controllers.MeetingAgendaItemsController, dependencyInjection.IDependencyInjectionContainer) {
	fmt.Println(os.Getenv("DB_CONN_STR"))
	db, err := gorm.Open("postgres", os.Getenv("DB_CONN_STR"))

	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect database")
	}

	maps.Init()

	dic := dependencyInjection.Initialize(db)

	var groupsController *controllers.GroupsController
	var activitiesController *controllers.ActivitiesController
	var meetingAgendasController *controllers.MeetingAgendaItemsController

	dic.GetInstance(func(instance *controllers.GroupsController) {
		groupsController = instance
	})

	dic.GetInstance(func(instance *controllers.ActivitiesController) {
		activitiesController = instance
	})

	dic.GetInstance(func(instance *controllers.MeetingAgendaItemsController) {
		meetingAgendasController = instance
	})

	return groupsController, activitiesController, meetingAgendasController, dic
}
