package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yousefvand/secret-service/pkg/client"
)

func init() {
	rootCmd.AddCommand(pingCmd)
}

var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Ping secretserviced daemon",
	Long:  `Send ping signal to secretserviced daemon and wait for pong response`,
	Run: func(cmd *cobra.Command, args []string) {

		ssClient, _ := client.New()
		response, _ := ssClient.SecretServiceCommand("ping", "")

		if response == "pong" {
			fmt.Println("secretserviced is up and responsive")
		} else {
			fmt.Println("Something is wrong with secretserviced")
		}
	},
}
