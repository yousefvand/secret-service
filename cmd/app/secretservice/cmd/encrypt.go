package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yousefvand/secret-service/pkg/client"
)

func init() {
	rootCmd.AddCommand(encryptCmd)
	encryptCmd.Flags().StringP("password", "p", "", "database password")

	encryptCmd.Flags().StringP("input", "i", "", "input file")
	encryptCmd.Flags().StringP("output", "o", "", "output file")
}

var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "encrypt a database",
	Long:  `encrypt a credential database using provided password`,
	Run: func(cmd *cobra.Command, _ []string) {

		password, _ := cmd.Flags().GetString("password")
		input, _ := cmd.Flags().GetString("input")
		output, _ := cmd.Flags().GetString("output")
		params := password + "\n" + input + "\n" + output

		ssClient, _ := client.New()
		response, _ := ssClient.SecretServiceCommand("encrypt database", params)

		if response == "ok" {
			fmt.Println("database encrypted successfully")
		} else {
			fmt.Println("database encryption failed!")
		}
	},
}
