package main

import (
	"fmt"
	"io"
	"log"

	"github.com/neibla/streaming-grpc-with-mtls/lib/auth"
	"github.com/neibla/streaming-grpc-with-mtls/proto"
)

type server struct {
	proto.UnimplementedAcknowledgementServiceServer
}

func (s server) Ack(stream proto.AcknowledgementService_AckServer) error {
	ctx := stream.Context()
	user := auth.GetUserFromContext(ctx)

	for {
		//exit if stream context complete
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Printf("Error on receive: %v", err)
			continue
		}

		resp := proto.AckResponse{Message: fmt.Sprintf("%s - acknowledged, user %s", req.Message, user.ID)}
		if err := stream.Send(&resp); err != nil {
			log.Printf("Failed to send: %v", err)
		}
	}
}
