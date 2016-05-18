package cmd

import (
	//	"fmt"

	"github.com/spf13/cobra"
)

var jsonServeAddress string
var jsonSecure bool

// serve_jsonCmd represents the serve-json command
var serve_jsonCmd = &cobra.Command{
	Use:   "serve-json",
	Short: "serves json",
	Long:  `Serves resources service over http(s)/json on the specified address`,

	Run: func(cmd *cobra.Command, args []string) {
		serveJson()
	},
}

func init() {
	RootCmd.AddCommand(serve_jsonCmd)
	serve_jsonCmd.Flags().StringVarP(&jsonServeAddress, "address", "a", "localhost:8089", "Address to listen on")
	serve_jsonCmd.Flags().BoolVarP(&jsonSecure, "secure", "s", false, "Whether or not to use tls")
}

func serveJson() {
	println(jsonServeAddress)
	println(jsonSecure)
}
