package cmd

import (
	"bufio"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yousefvand/secret-service/pkg/crypto"
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

		if len(password) != 32 {
			panic("Wrong password length. Password should be exactly 32 characters.")
		}

		if exist, _ := fileOrFolderExists(input); !exist {
			panic("Input file doesn't exist")
		}

		file, err := os.Open(input)

		if err != nil {
			panic("failed to open input file")
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)
		var i int = 0
		var signature bool = false
		var encrypted []string

		for scanner.Scan() {
			i++
			line := scanner.Text()
			if line == " \"encrypted\": false," {
				signature = true
			}
			if i > 2 && !signature {
				panic("Wrong type of file. This file is not marked as non \"encrypted\"!")
			}
			if index := strings.Index(line, "secretText"); index > 0 {
				text := line[21 : len(line)-1]
				cipher, err := crypto.EncryptAESCBC256(password, text)
				if err != nil {
					panic("Encryption failed: " + err.Error())
				}
				newLine := line[:21] + cipher + "\""
				encrypted = append(encrypted, newLine)
			} else {
				encrypted = append(encrypted, line)
			}
		}

		fileContent := strings.Join(encrypted, "\n")
		err = ioutil.WriteFile(output, []byte(fileContent), 0755)
		if err != nil {
			panic("Writing to output file failed: " + output)
		}

	},
}
