package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yousefvand/secret-service/pkg/client"
)

func init() {
	rootCmd.AddCommand(decryptCmd)
	decryptCmd.Flags().StringP("password", "p", "", "database password")

	decryptCmd.Flags().StringP("input", "i", "", "input file")
	decryptCmd.Flags().StringP("output", "o", "", "output file")
}

var decryptCmd = &cobra.Command{
	Use:   "decrypt",
	Short: "decrypt a database",
	Long:  `decrypt a credential database using provided password`,
	Run: func(cmd *cobra.Command, _ []string) {

		password, _ := cmd.Flags().GetString("password")
		input, _ := cmd.Flags().GetString("input")
		output, _ := cmd.Flags().GetString("output")
		params := password + "\n" + input + "\n" + output

		ssClient, _ := client.New()
		response, _ := ssClient.SecretServiceCommand("decrypt database", params)

		if response == "ok" {
			fmt.Println("database decrypted successfully")
		} else {
			fmt.Println("database decryption failed!")
		}
	},
}
