package gapi

import (
	"context"
	"database/sql"
	db "simplebank/db/model"
	"simplebank/factory"
	"simplebank/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginOwner(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {

	owner, err := server.store.GetOwnerByEmail(ctx, req.GetEmail())
	if err != nil {

		if err == sql.ErrNoRows {

			return nil, status.Errorf(codes.NotFound, "owner not found: %s", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to get owner: %s", err)
	}

	err = factory.CheckPassword(req.GetPassword(), owner.HashedPassword)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "invalid password: %s", err)
	}

	accessToken, accesspayload, err := server.tokenMaker.CreateToken(owner.Email, server.config.TokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create access token: %s", err)

	}

	refreshToken, refreshpayload, err := server.tokenMaker.CreateToken(owner.Email, server.config.RefreshTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create refresh token: %s", err)
	}

	mtdt := server.extractMetadata(ctx)

	sessionParams := db.CreateSessionParams{
		ID:           refreshpayload.ID,
		OwnerID:      owner.ID,
		Email:        owner.Email,
		RefreshToken: refreshToken,
		ClientIp:     mtdt.ClientIP,
		UserAgent:    mtdt.UserAgent,
		IsBlocked:    false,
		ExpiresAt:    refreshpayload.ExpiredAt,
	}

	session, err := server.store.CreateSession(ctx, sessionParams)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create session: %s", err)
	}

	response := &pb.LoginUserResponse{
		SessionId:       string(session.ID.String()),
		AccessToken:     accessToken,
		AccessTokenExp:  timestamppb.New(accesspayload.ExpiredAt),
		RefreshToken:    refreshToken,
		RefreshTokenExp: timestamppb.New(refreshpayload.ExpiredAt),
		User:            convertOwnerResponse(owner),
	}
	return response, nil
}
