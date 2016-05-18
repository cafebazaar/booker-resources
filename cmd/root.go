package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var addr string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "resources",
	Short: "Server for resources microservice",
	Long: `To get started run the serve subcommand which will start a server
on localhost:10000:
    resources serve -a localhost:10000
Then test ir:
    curl -k https://localhost:10000/healthz
`,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&addr, "address", "a", "127.0.0.1:8000", "Address to listen on")
}
