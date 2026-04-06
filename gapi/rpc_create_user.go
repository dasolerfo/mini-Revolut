package gapi

import (
	"context"
	db "simplebank/db/model"
	"simplebank/factory"
	"simplebank/pb"
	"time"

	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateOwner(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	hash_pass, err := factory.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash the password: %s", err)

	}

	date, err := time.Parse("2006-01-02", req.GetBornAt())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to parse birth date: %s", err)

	}

	arg := db.
		CreateOwnerParams{
		FirstName:      req.GetFirstName(),
		FirstSurname:   req.GetFirstSurname(),
		SecondSurname:  req.GetSecondSurname(),
		HashedPassword: hash_pass,
		Nationality:    req.GetNationality(),
		BornAt:         date,
		Email:          req.GetEmail(),
	}

	owner, err := server.store.CreateOwner(ctx, arg)

	if err != nil {
		if pqError, ok := err.(*pq.Error); ok {
			switch pqError.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "failed to create owner, already exists: %s", err)
			}
		}
		return nil, status.Errorf(codes.Internal, "failed to create owner: %s", err)
	}

	ownerResponse := &pb.CreateUserResponse{User: convertOwnerResponse(owner)}

	return ownerResponse, nil
}
