package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"

	"github.com/spf13/viper"

	"github.com/neibla/streaming-grpc-with-mtls/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	config, err := LoadConfig(".")
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}
	ctx := context.Background()
	client := newClient(config)
	stream, err := client.Ack(ctx)
	if err != nil {
		log.Fatalf("Error opening stream: %v", err)
	}

	done := make(chan bool)

	//sender
	go func() {
		for i := 0; i < 10; i++ {
			req := proto.AckRequest{Message: fmt.Sprintf("hello %d", i)}
			if err := stream.Send(&req); err != nil {
				log.Fatalf("Failed to send %v", err)
			}
		}
		if err := stream.CloseSend(); err != nil {
			log.Fatalf("Failed to close send %v", err)
		}
	}()

	//receiver
	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				close(done)
				return
			}
			if err != nil {
				log.Fatalf("Failed to receive %v", err)
			}
			log.Printf("Received: %s", resp.Message)
		}
	}()

	select {
	case <-done:
	case <-ctx.Done():
	}
}

type Config struct {
	AuthorizedCA string `mapstructure:"AUTHORIZED_CA"`
	ClientKey    string `mapstructure:"CLIENT_KEY"`
	ClientCert   string `mapstructure:"CLIENT_CERT"`
	ServerURL    string `mapstructure:"SERVER_URL"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

func newClient(config Config) proto.AcknowledgementServiceClient {
	tlsCredentials, err := loadTLSCredentials(config)
	if err != nil {
		log.Fatalf("Failed to load tlsCredentials: %v, config: %+v", err, config)
	}

	conn, err := grpc.Dial(config.ServerURL, grpc.WithTransportCredentials(tlsCredentials))
	if err != nil {
		log.Fatalf("Cannot connect with server. error: %v, config: %+v", err, config)
	}
	client := proto.NewAcknowledgementServiceClient(conn)
	return client
}

func loadTLSCredentials(config Config) (credentials.TransportCredentials, error) {

	cert, err := tls.LoadX509KeyPair(config.ClientCert, config.ClientKey)
	if err != nil {
		return nil, err
	}
	caCert, err := ioutil.ReadFile(config.AuthorizedCA)
	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		return nil, errors.New("Failed setup CA cert")
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
		MinVersion:   tls.VersionTLS13,
	}
	return credentials.NewTLS(tlsConfig), nil
}
