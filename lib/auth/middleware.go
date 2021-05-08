package auth

import (
	"context"
	"errors"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

var ErrPeerRequired = errors.New("peer required")
var ErrTlsInformationRequired = errors.New("TLS required")
var ErrInvalidCertificate = errors.New("invalid certificate")

// TODO: deal with revoked certificates
func GRPCAuthenticationMiddleware(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	authenticatedCtx, err := mTLSAuthenticationFromContext(ctx)
	if err != nil {
		return nil, status.Errorf(
			codes.Unauthenticated, err.Error())
	}
	return handler(authenticatedCtx, req)
}

func GRPCAuthenticationStreamMiddleware(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	authenticatedCtx, err := mTLSAuthenticationFromContext(stream.Context())
	if err != nil {
		return status.Errorf(
			codes.Unauthenticated, err.Error())
	}
	wrapped := grpc_middleware.WrapServerStream(stream)
	wrapped.WrappedContext = authenticatedCtx
	return handler(srv, wrapped)
}

func mTLSAuthenticationFromContext(ctx context.Context) (context.Context, error) {
	peer, ok := peer.FromContext(ctx)
	if !ok {
		return nil, ErrPeerRequired
	}
	tlsInfo, ok := peer.AuthInfo.(credentials.TLSInfo)
	if !ok {
		return nil, ErrTlsInformationRequired
	}
	if len(tlsInfo.State.VerifiedChains) == 0 || len(tlsInfo.State.VerifiedChains[0]) == 0 {
		return nil, ErrInvalidCertificate
	}

	userID := tlsInfo.State.VerifiedChains[0][0].Subject.CommonName
	// TODO: look up user in system
	user := &User{
		ID: userID,
	}

	return ContextWithUser(ctx, user), nil
}
