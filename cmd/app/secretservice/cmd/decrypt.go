package cmd

import (
	"bufio"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yousefvand/secret-service/pkg/crypto"
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
		var decrypted []string

		for scanner.Scan() {
			i++
			line := scanner.Text()
			if line == " \"encrypted\": true," {
				signature = true
			}
			if i > 2 && !signature {
				panic("Wrong type of file. This file is not marked as \"encrypted\"!")
			}
			if index := strings.Index(line, "secretText"); index > 0 {
				cipher := line[21 : len(line)-1]
				decrypt, err := crypto.DecryptAESCBC256(password, cipher)
				if err != nil {
					panic("Decryption failed: " + err.Error())
				}
				newLine := line[:21] + decrypt + "\""
				decrypted = append(decrypted, newLine)
			} else {
				decrypted = append(decrypted, line)
			}
		}

		fileContent := strings.Join(decrypted, "\n")
		err = os.WriteFile(output, []byte(fileContent), 0755)
		if err != nil {
			panic("Writing to output file failed: " + output)
		}

	},
}

func fileOrFolderExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
