package gapi

import (
	"context"
	"log"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	grpcMetadataUserAgentKey = "rpcgateway-user-agent"
	userAgentHeaderKey       = "user-agent"
	grpcMetadataClientIPKey  = "x-forwarded-for"
)

type Metadata struct {
	UserAgent string
	ClientIP  string
}

func (server *Server) extractMetadata(ctx context.Context) *Metadata {
	mtdt := &Metadata{}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		log.Printf("Metadata received: %v\n", md)

		userAgents := md.Get(grpcMetadataUserAgentKey)
		if len(userAgents) > 0 {
			mtdt.UserAgent = userAgents[0]
		}
		userAgents = md.Get(userAgentHeaderKey)
		if len(userAgents) > 0 {
			mtdt.UserAgent = userAgents[0]
		}

		clientIPs := md.Get(grpcMetadataClientIPKey)
		if len(clientIPs) > 0 {
			mtdt.ClientIP = clientIPs[0]
		}
	}

	if p, ok := peer.FromContext(ctx); ok {
		log.Printf("Peer info received: %v\n", p)
		if p.Addr != nil {
			mtdt.ClientIP = p.Addr.String()
		}
	}

	return mtdt

}
