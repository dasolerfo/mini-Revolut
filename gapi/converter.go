package gapi

import (
	db "simplebank/db/model"
	"simplebank/pb"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertOwnerResponse(owner db.Owner) *pb.User {
	return &pb.User{
		FirstName:     owner.FirstName,
		FirstSurname:  owner.FirstSurname,
		SecondSurname: owner.SecondSurname,
		Nationality:   owner.Nationality,
		CreatedAt:     timestamppb.New(owner.CreatedAt),
		Email:         owner.Email,
	}

}
