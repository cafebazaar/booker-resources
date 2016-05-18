package cmd

import (
	//	"fmt"

	"github.com/spf13/cobra"
)

var grpcServeAddress string
var grpcSecure bool

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "serves grpc",
	Long:  "Serves resources service over grpc on the specified address",
	Run: func(cmd *cobra.Command, args []string) {
		serveGrpc()
	},
}

func init() {
	RootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringVarP(&grpcServeAddress, "address", "a", "localhost:8089", "Address to listen on")
	serveCmd.Flags().BoolVarP(&grpcSecure, "secure", "s", false, "Whether or not to use tls")
}

func serveGrpc() {
	println(grpcServeAddress)
	println(grpcSecure)
}
