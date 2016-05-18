package cmd

import (
	"crypto/tls"
	"net"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/gengo/grpc-gateway/runtime"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/cafebazaar/booker-resources/common"
	"github.com/cafebazaar/booker-resources/proto"
)

var grpcAddr string
var throughTLS bool

// serveCmd represents the serve command
var (
	serveGWCmd = &cobra.Command{
		Use:   "serve-json",
		Short: "Launches the grpc-gateway on the given address",
		Run: func(cmd *cobra.Command, args []string) {
			serveGW()
		},
	}
)

func init() {
	serveGWCmd.Flags().StringVarP(&grpcAddr, "service-address", "s", "127.0.0.1:8000", "Address of the service, to be proxified")
	serveGWCmd.Flags().BoolVarP(&throughTLS, "https", "", false, "Serve on https protocol")
	RootCmd.AddCommand(serveGWCmd)
}

func serveGW() {

	printVersion()

	var err error
	ctx := context.Background()

	cp, err := certPool()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to load the certPool")
	}

	dcreds := credentials.NewTLS(&tls.Config{
		ServerName: common.ConfigString("SERVER_NAME"),
		RootCAs:    cp,
	})
	dopts := []grpc.DialOption{grpc.WithTransportCredentials(dcreds)}

	var mux http.Handler
	gwmux := runtime.NewServeMux()
	mux = healthzAddedHandlerFunc(gwmux)

	err = proto.RegisterResourcesHandlerFromEndpoint(ctx, gwmux, grpcAddr, dopts)
	if err != nil {
		logrus.WithError(err).Fatal("Failed while RegisterResourcesHandlerFromEndpoint")
		return
	}

	if throughTLS {
		kp, err := keyPair()
		if err != nil {
			logrus.WithError(err).Fatal("Failed to load the keyPair")
		}

		conn, err := net.Listen("tcp", addr)
		if err != nil {
			logrus.WithError(err).Fatalf("Failed to listen on %s", addr)
		}

		srv := &http.Server{
			Addr:    addr,
			Handler: mux,
			TLSConfig: &tls.Config{
				Certificates: []tls.Certificate{*kp},
			},
		}

		logrus.Infof("Starting at %s (HTTPS)", addr)
		err = srv.Serve(tls.NewListener(conn, srv.TLSConfig))
	} else {
		logrus.Infof("Starting at %s (HTTP)", addr)
		err = http.ListenAndServe(addr, mux)
	}

	if err != nil {
		logrus.WithError(err).Fatalf("Stopped serving at %s", addr)
	}
}
