package main

import (
	"database/sql"
	"log"
	"net"
	"simplebank/api"
	db "simplebank/db/model"
	"simplebank/factory"
	"simplebank/gapi"
	"simplebank/pb"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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
