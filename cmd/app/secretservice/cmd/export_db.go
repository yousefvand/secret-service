package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yousefvand/secret-service/pkg/client"
)

func init() {
	rootCmd.AddCommand(exportDbCmd)
}

var exportDbCmd = &cobra.Command{
	Use:   "export db",
	Short: "export db exports an unencrypted version of db",
	Long:  `export db exports an unencrypted version of credentials database`,
	Run: func(_ *cobra.Command, _ []string) {

		ssClient, _ := client.New()
		response, _ := ssClient.SecretServiceCommand("export database", "")

		if response == "ok" {
			fmt.Println("database export completed successfully")
		} else {
			fmt.Println("database export failed!")
		}
	},
}
