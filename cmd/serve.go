package cmd

import (
	"crypto/tls"
	"net"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/cafebazaar/booker-resources/api"
)

// serveCmd represents the serve command
var (
	serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "Launches the grpc server",
		Run: func(cmd *cobra.Command, args []string) {
			serveService()
		},
	}
)

func init() {
	RootCmd.AddCommand(serveCmd)
}

func serveService() {

	printVersion()

	kp, err := keyPair()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to load the keyPair")
	}

	opts := []grpc.ServerOption{
		grpc.Creds(credentials.NewServerTLSFromCert(kp))}

	grpcServer := grpc.NewServer(opts...)

	api.RegisterServer(grpcServer)

	conn, err := net.Listen("tcp", addr)
	if err != nil {
		logrus.WithError(err).Fatalf("Failed to listen on %s", addr)
	}

	srv := &http.Server{
		Addr:    addr,
		Handler: healthzAddedHandlerFunc(grpcServer),
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{*kp},
		},
	}

	logrus.Infof("Starting at    %s", addr)
	err = srv.Serve(tls.NewListener(conn, srv.TLSConfig))

	if err != nil {
		logrus.WithError(err).Fatalf("Stopped serving at     %s", addr)
	}
}
