package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/rs/cors"

	controllers "github.com/vireocloud/property-pros-service/src/controllers"
	dependencyInjection "github.com/vireocloud/property-pros-service/src/dependencyInjectionManagement"
	lib "github.com/vireocloud/property-pros-service/src/grpc"
	middleware "github.com/vireocloud/property-pros-service/src/middleware"
	maps "github.com/vireocloud/property-pros-service/src/objectMaps"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

var (
	enableTls = flag.Bool("enable_tls", false, "Use TLS - required for HTTP2.")
	// tlsCertFilePath = flag.String("tls_cert_file", "../misc/localhost.crt", "Path to the CRT/PEM file.")
	// tlsKeyFilePath  = flag.String("tls_key_file", "../misc/localhost.key", "Path to the private key file.")
	GO_ENV                       = os.Getenv("GO_ENV")
	groupsController             *controllers.GroupsController
	activitiesController         *controllers.ActivitiesController
	meetingAgendaItemsController *controllers.MeetingAgendaItemsController
	dic                          dependencyInjection.IDependencyInjectionContainer
)

// main start a gRPC server and waits for connection
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	flag.Parse()

	port, err := strconv.ParseInt(os.Getenv("HTTP_PLATFORM_PORT"), 10, 32)

	if err != nil {
		port = 8900
	}

	if *enableTls {
		port = 9091
	}

	http2Port := port + 1

	groupsController, activitiesController, meetingAgendaItemsController, dic = getControllers()

	defer dic.Teardown()

	var serverWaitgroup sync.WaitGroup
	serverWaitgroup.Add(2)

	go func() {
		defer serverWaitgroup.Done()
		err := startGRPCHttp1Server(port)
		if err != nil {
			log.Fatalf("failed to start gRPC server: %s", err)
		}
	}()

	go func() {
		defer serverWaitgroup.Done()
		err := startGRPCHttp2Server(http2Port)
		if err != nil {
			log.Fatalf("failed to start gRPC server: %s", err)
		}

	}()

	serverWaitgroup.Wait()
	// infinite loop
	log.Printf("Entering infinite loop")
	select {}
}

func startGRPCHttp2Server(http2Port int64) error {

	grpclog.Printf("Starting server. tcp port: %d, with TLS: %v", http2Port, *enableTls)

	// create a listener on TCP port 7777
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", http2Port))

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcPingServer := grpc.NewServer()
	// attach the Ping service to the server

	lib.RegisterGroupServiceServer(grpcPingServer, groupsController)
	lib.RegisterActivityServiceServer(grpcPingServer, activitiesController)
	lib.RegisterMeetingAgendaItemServiceServer(grpcPingServer, meetingAgendaItemsController)

	// start the server
	return grpcPingServer.Serve(lis)
}

func startGRPCHttp1Server(port int64) error {

	grpcServer := grpc.NewServer()

	lib.RegisterGroupServiceServer(grpcServer, groupsController)
	lib.RegisterActivityServiceServer(grpcServer, activitiesController)
	lib.RegisterMeetingAgendaItemServiceServer(grpcServer, meetingAgendaItemsController)

	grpclog.SetLogger(log.New(os.Stdout, "localize: ", log.LstdFlags))

	wrappedServer := grpcweb.WrapServer(grpcServer)

	router := chi.NewRouter()

	var corsOptions cors.Options

	corsOptions = cors.Options{
		AllowedOrigins:   []string{"*" /*"http://localhost:3002"*/}, //don't be fooled by the slice type.  this only works with one value.  it will break if you add more than one domain
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"x-grpc-web", "Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Access-Control-Request-Headers", "Access-Control-Request-Method", "DNT", "Origin", "Referer", "User-Agent"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
		Debug:            true,
	}

	router.Use(chiMiddleware.Logger)
	router.Use(chiMiddleware.Recoverer) // Must come before general CORS handling
	router.Use(middleware.NewGrpcWebMiddleware(wrappedServer, corsOptions).Handler)
	router.Use(cors.New(corsOptions).Handler)

	// the server won't work without this
	// apparently go-chi requires at least one route to work
	router.Get("/test", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "test")
	}))

	grpclog.Printf("Starting server. http port: %d, with TLS: %v", port, *enableTls)

	endpoints := grpcweb.ListGRPCResources(grpcServer)
	grpclog.Printf("With registered endpoints: ")
	for _, endpoint := range endpoints {
		grpclog.Printf("%v\n", endpoint)
	}
	// grpclog.Printf("With registered endpoints: \r\n%v\n", grpcweb.ListGRPCResources(grpcServer))

	return http.ListenAndServe(":"+strconv.Itoa(int(port)), router)
}

func getControllers() (*controllers.GroupsController, *controllers.ActivitiesController, *controllers.MeetingAgendaItemsController, dependencyInjection.IDependencyInjectionContainer) {
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
