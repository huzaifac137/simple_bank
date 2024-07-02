package gapi

import (
	"context"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	grpcGatewayUserAgent     = "grpcgateway-user-agent"
	grpcGatewayXForwardedFor = "x-forwarded-for"
	grpcUserAgent            = "user-agent"
)

type Metadata struct {
	UserAgent string
	ClientIp  string
}

func (s *Server) ExtractMetadata(ctx context.Context) *Metadata {
	mtdt := &Metadata{}

	if grpcMetadata, ok := metadata.FromIncomingContext(ctx); ok {

		if gUA := grpcMetadata.Get(grpcGatewayUserAgent); len(gUA) > 0 {
			mtdt.UserAgent = gUA[0]
		}
		if GClientIp := grpcMetadata.Get(grpcGatewayXForwardedFor); len(GClientIp) > 0 {
			mtdt.ClientIp = GClientIp[0]
		}

		if ua := grpcMetadata.Get(grpcUserAgent); len(ua) > 0 {
			mtdt.UserAgent = ua[0]
		}
	}

	if p, ok := peer.FromContext(ctx); ok {
		if ip := p.Addr.String(); len(ip) > 0 {
			mtdt.ClientIp = ip
		}
	}
	return mtdt
}
