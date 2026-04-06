package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"
	"simplebank/api"
	db "simplebank/db/model"
	"simplebank/factory"
	"simplebank/gapi"
	"simplebank/pb"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	config, err := factory.LoadConfig(".")
	if err != nil {
		log.Fatal("Error! No s'ha pogut carregar el .env", err)
	}
	testDB, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Error! No et pots connectar a la base de dades: ", err)
	}

	store := db.NewStore(testDB)
	go runGatewayServer(config, store)
	runGRPCServer(config, store)

}

func runGinServer(config factory.Config, store db.Store) {
	router, err := api.NewServer(config, store)

	if err != nil {
		log.Fatal("No es pot inicialitzar el server: ", err)
	}

	err = router.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("Error! No es pot inicialitzar el server: ", err)
	}
}

func runGRPCServer(config factory.Config, store db.Store) {
	router, err := gapi.NewServer(config, store)

	if err != nil {
		log.Fatal("No es pot inicialitzar el server: ", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServiceServer(grpcServer, router)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	log.Printf("Starting gRPC server at %s", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("Error! No es pot inicialitzar el listener: ", err)
	}
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("Error! No es pot inicialitzar el server: ", err)
	}
}

func runGatewayServer(config factory.Config, store db.Store) {
	router, err := gapi.NewServer(config, store)

	if err != nil {
		log.Fatal("No es pot inicialitzar el server:    ", err)
	}
	grpcMux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	}))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterSimpleBankServiceHandlerServer(ctx, grpcMux, router)
	if err != nil {
		log.Fatal("Error! No es pot inicialitzar el handler: ", err)
	}
	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	log.Printf("Starting gRPC gateway server at %s", config.HTTPServerAddress)
	if err != nil {
		log.Fatal("Error! No es pot inicialitzar el listener: ", err)
	}
	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal("Error! No es pot inicialitzar el server: ", err)
	}
}
