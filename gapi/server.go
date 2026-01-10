package gapi

import (
	"fmt"
	db "simplebank/db/model"
	"simplebank/factory"
	"simplebank/pb"
	token "simplebank/token"
)

// Server serves gRPC requests for our banking service.
type Server struct {
	pb.UnimplementedUserServiceServer
	config     factory.Config
	tokenMaker token.Maker
	store      db.Store
}

func NewServer(config factory.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %v", err)
	}
	server := &Server{store: store,
		tokenMaker: tokenMaker,
		config:     config}

	return server, nil
}
