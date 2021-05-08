package main

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"log"
	"net"

	"github.com/spf13/viper"

	"github.com/neibla/streaming-grpc-with-mtls/lib/auth"
	"github.com/neibla/streaming-grpc-with-mtls/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {

	config, err := LoadConfig(".")
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}
	tlsCredentials, err := loadTLSCredentials(config)
	if err != nil {
		log.Fatalf("Failed load tls credentials from config %+v, err: %v", config, err)
	}
	grpcServer := grpc.NewServer(grpc.Creds(tlsCredentials),
		grpc.StreamInterceptor(auth.GRPCAuthenticationStreamMiddleware),
		grpc.UnaryInterceptor(auth.GRPCAuthenticationMiddleware))

	proto.RegisterAcknowledgementServiceServer(grpcServer, server{})

	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Unable to listen to port: %v", err)
	}

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}

type Config struct {
	AuthorizedCA string `mapstructure:"AUTHORIZED_CA"`
	ServerKey    string `mapstructure:"SERVER_KEY"`
	ServerCert   string `mapstructure:"SERVER_CERT"`
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

func loadTLSCredentials(config Config) (credentials.TransportCredentials, error) {
	cert, err := tls.LoadX509KeyPair(config.ServerCert, config.ServerKey)
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
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    caCertPool,
		MinVersion:   tls.VersionTLS13,
	}
	return credentials.NewTLS(tlsConfig), nil
}
